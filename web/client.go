package web

import (
	"context"
	"github.com/boostgo/lite/types/to"
	"net/http"
	"time"
)

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

func (client *Client) SetBaseURL(baseURL string) *Client {
	client.baseURL = baseURL
	return client
}

func (client *Client) Logging(logging bool) *Client {
	client.logging = logging
	return client
}

func (client *Client) Client(httpClient http.Client) *Client {
	client.client = &httpClient
	return client
}

func (client *Client) RetryCount(count int) *Client {
	if count <= 1 {
		return client
	}

	client.retryCount = count
	return client
}

func (client *Client) RetryWait(wait time.Duration) *Client {
	if wait <= 0 {
		return client
	}

	client.retryWait = wait
	return client
}

func (client *Client) Timeout(timeout time.Duration) *Client {
	if timeout <= 0 {
		return client
	}

	client.timeout = timeout
	return client
}

func (client *Client) BasicAuth(username, password string) *Client {
	if username == "" {
		return client
	}

	client.basic = basicAuth{username, password}
	return client
}

func (client *Client) BearerToken(token string) *Client {
	client.bearerToken = token
	return client
}

func (client *Client) Options(opts ...RequestOption) *Client {
	if len(opts) == 0 {
		return client
	}

	client.options = opts
	return client
}

func (client *Client) Header(key string, value any) *Client {
	client.headers[key] = to.String(value)
	return client
}

func (client *Client) Headers(headers map[string]any) *Client {
	for key, value := range headers {
		client.headers[key] = value
	}
	return client
}

func (client *Client) Cookie(key string, value any) *Client {
	client.cookies[key] = to.String(value)
	return client
}

func (client *Client) Cookies(cookies map[string]any) *Client {
	for key, value := range cookies {
		client.cookies[key] = value
	}

	return client
}

func (client *Client) Queries(queries map[string]any) *Client {
	for key, value := range queries {
		client.queryVariables[key] = value
	}
	return client
}

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
