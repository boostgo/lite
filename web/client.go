package web

import (
	"context"
	"github.com/boostgo/lite/types/to"
	"net/http"
	"time"
)

// Client web client which allow to send HTTP requests.
// It can simplify sending requests by containing base url, headers and cookies for created requests from this client.
// There is retry mechanism for many request sending till it proceed success.
type Client struct {
	baseURL string
	logging bool
	client  *http.Client

	retryCount int
	retryWait  time.Duration

	timeout time.Duration

	basic       basicAuth
	bearerToken string

	options []RequestOption

	headers        map[string]any
	cookies        map[string]any
	queryVariables map[string]any
}

func New() *Client {
	return &Client{
		logging: true,

		headers:        make(map[string]any),
		cookies:        make(map[string]any),
		queryVariables: make(map[string]any),

		options: make([]RequestOption, 0),
	}
}

// SetBaseURL sets base url for every nested requests
func (client *Client) SetBaseURL(baseURL string) *Client {
	client.baseURL = baseURL
	return client
}

// Logging setting logging mode.
// Logging mode turns on inner logs (mostly errors).
func (client *Client) Logging(logging bool) *Client {
	client.logging = logging
	return client
}

// Client set default http client for every nested request
func (client *Client) Client(httpClient http.Client) *Client {
	client.client = &httpClient
	return client
}

// RetryCount sets count of retries need.
// By default, retry count is 1
func (client *Client) RetryCount(count int) *Client {
	if count <= 1 {
		return client
	}

	client.retryCount = count
	return client
}

// RetryWait sets wait time between retry requests.
// Default wait time is 100ms.
func (client *Client) RetryWait(wait time.Duration) *Client {
	if wait <= 0 {
		return client
	}

	client.retryWait = wait
	return client
}

// Timeout sets timeout for waiting for request.
// By default, there is no timeout
func (client *Client) Timeout(timeout time.Duration) *Client {
	if timeout <= 0 {
		return client
	}

	client.timeout = timeout
	return client
}

// BasicAuth sets username & password for basic auth mechanism
func (client *Client) BasicAuth(username, password string) *Client {
	if username == "" {
		return client
	}

	client.basic = basicAuth{username, password}
	return client
}

// BearerToken sets token for "Authorization" header.
// Prefix "Bearer " sets automatically.
func (client *Client) BearerToken(token string) *Client {
	client.bearerToken = token
	return client
}

// Options sets option functions which can modify created request.
func (client *Client) Options(opts ...RequestOption) *Client {
	if len(opts) == 0 {
		return client
	}

	client.options = opts
	return client
}

// Header set one more key-value pair to headers.
// If key already exist it rewrites existing key value.
func (client *Client) Header(key string, value any) *Client {
	client.headers[key] = to.String(value)
	return client
}

// Headers sets map of key-value pairs.
// Existing keys will be rewritten.
func (client *Client) Headers(headers map[string]any) *Client {
	for key, value := range headers {
		client.headers[key] = value
	}
	return client
}

// Authorization sets authorization token for "Authorization" header.
// Prefix "Bearer " will set automatically
func (client *Client) Authorization(token string) *Client {
	client.headers["Authorization"] = "Bearer " + token
	return client
}

// ContentType sets header value to "Content-Type" header.
func (client *Client) ContentType(contentType string) *Client {
	client.headers["Content-Type"] = contentType
	return client
}

// Cookie sets new cookie to request.
// Existing key will be rewritten.
func (client *Client) Cookie(key string, value any) *Client {
	client.cookies[key] = to.String(value)
	return client
}

// Cookies sets new cookies map to request.
// Existing keys will be rewritten.
func (client *Client) Cookies(cookies map[string]any) *Client {
	for key, value := range cookies {
		client.cookies[key] = value
	}

	return client
}

// Query add new query param to nested requests.
func (client *Client) Query(key string, value any) *Client {
	client.queryVariables[key] = value
	return client
}

// Queries sets query params to request.
// Existing keys will be rewritten
func (client *Client) Queries(queries map[string]any) *Client {
	for key, value := range queries {
		client.queryVariables[key] = value
	}
	return client
}

// R creates Request object with nested settings from current Client
func (client *Client) R(ctx context.Context) *Request {
	return R(ctx).
		setBaseURL(client.baseURL).
		Logging(client.logging).
		Client(client.client).
		Headers(client.headers).
		Cookies(client.cookies).
		Queries(client.queryVariables).
		RetryCount(client.retryCount).
		RetryWait(client.retryWait).
		Timeout(client.timeout).
		BasicAuth(client.basic.username, client.basic.password).
		BearerToken(client.bearerToken).
		Options(client.options...)
}
