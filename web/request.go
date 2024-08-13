package web

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"github.com/boostgo/lite/log"
	"github.com/boostgo/lite/system/trace"
	"github.com/boostgo/lite/types/flex"
	"github.com/boostgo/lite/types/to"
	"io"
	"net/http"
	"time"
)

type RequestOption func(request *http.Request)

type Request struct {
	ctx     context.Context
	baseURL string
	client  *http.Client
	export  any

	queryVariables map[string]any
	headers        map[string]any
	cookies        map[string]any

	retryCount int
	retryWait  time.Duration

	timeout time.Duration

	logging   bool
	traceMode bool

	basic       basicAuth
	bearerToken string

	options []RequestOption

	req      *http.Request
	resp     *http.Response
	response *Response
}

func R(ctx context.Context) *Request {
	if ctx == nil {
		ctx = context.Background()
	}

	return &Request{
		ctx: ctx,

		queryVariables: make(map[string]any),
		headers:        make(map[string]any),
		cookies:        make(map[string]any),

		retryCount: 1,
		retryWait:  time.Millisecond * 100,

		options: make([]RequestOption, 0),

		logging:   true,
		traceMode: trace.AmIMaster(),
	}
}

func (request *Request) setBaseURL(baseURL string) *Request {
	request.baseURL = baseURL
	return request
}

func (request *Request) Options(opts ...RequestOption) *Request {
	if len(opts) == 0 {
		return request
	}

	request.options = opts
	return request
}

func (request *Request) Logging(logging bool) *Request {
	request.logging = logging
	return request
}

func (request *Request) Client(client *http.Client) *Request {
	if client == nil {
		return request
	}

	request.client = client
	return request
}

func (request *Request) RetryCount(count int) *Request {
	if count <= 1 {
		return request
	}

	request.retryCount = count
	return request
}

func (request *Request) RetryWait(wait time.Duration) *Request {
	if wait <= 0 {
		return request
	}

	request.retryWait = wait
	return request
}

func (request *Request) Timeout(timeout time.Duration) *Request {
	if timeout <= 0 {
		return request
	}

	request.timeout = timeout
	return request
}

func (request *Request) BasicAuth(username, password string) *Request {
	if username == "" {
		return request
	}

	request.basic = basicAuth{username, password}
	return request
}

func (request *Request) BearerToken(token string) *Request {
	request.bearerToken = token
	return request
}

func (request *Request) Query(key string, value any) *Request {
	request.queryVariables[key] = value
	return request
}

func (request *Request) Queries(queries map[string]any) *Request {
	for key, value := range queries {
		request.queryVariables[key] = value
	}
	return request
}

func (request *Request) Result(export any) *Request {
	if !flex.Type(export).IsPtr() {
		return request
	}

	request.export = export
	return request
}

func (request *Request) Header(key string, value any) *Request {
	request.headers[key] = value
	return request
}

func (request *Request) Headers(headers map[string]any) *Request {
	for key, value := range headers {
		request.headers[key] = value
	}
	return request
}

func (request *Request) Cookie(key string, value any) *Request {
	request.cookies[key] = value
	return request
}

func (request *Request) Cookies(cookies map[string]any) *Request {
	for key, value := range cookies {
		request.cookies[key] = value
	}
	return request
}

func (request *Request) GET(url string, body ...any) (*Response, error) {
	return request.retryDo(http.MethodGet, url, body...)
}

func (request *Request) POST(url string, body ...any) (*Response, error) {
	return request.retryDo(http.MethodPost, url, body...)
}

func (request *Request) PUT(url string, body ...any) (*Response, error) {
	return request.retryDo(http.MethodPut, url, body...)
}

func (request *Request) PATCH(url string, body ...any) (*Response, error) {
	return request.retryDo(http.MethodPatch, url, body...)
}

func (request *Request) DELETE(url string, body ...any) (*Response, error) {
	return request.retryDo(http.MethodDelete, url, body...)
}

func (request *Request) OPTIONS(url string, body ...any) (*Response, error) {
	return request.retryDo(http.MethodOptions, url, body...)
}

func (request *Request) HEAD(url string, body ...any) (*Response, error) {
	return request.retryDo(http.MethodHead, url, body...)
}

func (request *Request) initRequest(method, url string, body ...any) error {
	if request.req != nil {
		return nil
	}

	// building request path (URL)
	var err error
	var fullURL string
	if request.baseURL != "" {
		fullURL = request.baseURL + url
	} else {
		fullURL = url
	}

	// creating request
	if len(body) > 0 && body[0] != nil {
		var bodyBlob []byte
		bodyBlob, err = json.Marshal(body[0])
		if err != nil {
			return err
		}

		request.req, err = http.NewRequest(method, fullURL, bytes.NewReader(bodyBlob))
	} else {
		request.req, err = http.NewRequest(method, fullURL, nil)
	}
	if err != nil {
		return err
	}

	// query variables
	query := request.req.URL.Query()
	for key, value := range request.queryVariables {
		query.Set(key, to.String(value))
	}
	request.req.URL.RawQuery = query.Encode()

	// auth
	request.initAuth()

	// headers
	for key, value := range request.headers {
		request.req.Header.Set(key, to.String(value))
	}

	// cookies
	for key, value := range request.cookies {
		request.req.AddCookie(&http.Cookie{Name: key, Value: to.String(value)})
	}

	// options
	for _, opt := range request.options {
		opt(request.req)
	}

	return nil
}

func (request *Request) retryDo(method, url string, body ...any) (*Response, error) {
	if request.timeout > 0 {
		var cancel context.CancelFunc
		request.ctx, cancel = context.WithTimeout(context.Background(), request.timeout)
		defer cancel()
	}

	var err error
	for i := 0; i < request.retryCount; i++ {
		isLast := i == request.retryCount-1

		request.response, err = request.do(method, url, body...)
		if err != nil && isLast {
			return nil, err
		}

		time.Sleep(request.retryWait)
	}

	return request.response, nil
}

func (request *Request) do(method, url string, body ...any) (*Response, error) {
	logger := log.Context(request.ctx, "web.request")

	var err error

	if err = request.initRequest(method, url, body...); err != nil {
		return nil, err
	}

	request.resp, err = request.getClient().Do(request.req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = request.resp.Body.Close(); err != nil {
			request.logError(logger, err, "Close response body")
		}
	}()

	request.response = newResponse(request.resp)

	var respBlob []byte
	respBlob, err = io.ReadAll(request.resp.Body)
	if err != nil {
		return nil, err
	}
	request.response.bodyBlob = respBlob

	if request.export != nil {
		if err = request.response.Parse(request.export); err != nil {
			request.logError(logger, err, "Parse response body")
		}
	}

	return request.response, nil
}

func (request *Request) logError(logger log.Logger, err error, msg string) {
	if !request.logging {
		return
	}

	logger.Error().Err(err).Msg(msg)
}

func (request *Request) getClient() *http.Client {
	if request.client != nil {
		return request.client
	}

	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}

func (request *Request) initAuth() {
	if request.basic != (basicAuth{}) && request.basic.username != "" {
		request.req.SetBasicAuth(request.basic.username, request.basic.password)
		return
	}

	if request.bearerToken != "" {
		request.req.Header.Set("Authorization", "Bearer "+request.bearerToken)
	}
}
