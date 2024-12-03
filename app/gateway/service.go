package gateway

import (
	"context"
	"github.com/boostgo/lite/errs"
	"github.com/boostgo/lite/web"
	"strings"
)

// Service search for need route & redirect request to match path.
type Service interface {
	// RegisterRoute append new routes
	RegisterRoute(routes ...Route) Service
	// Routes return all service routes
	Routes() []Route
	// Match find match [Route] searching by method & path
	Match(method, path string) (Route, bool)
	// Proxy make request to redirect path and proxy all headers, cookies and request body
	Proxy(ctx context.Context, r Route, request Request) (Response, error)
}

type service struct {
	client  *web.Client
	routes  []Route
	errType string
}

func NewService(baseURL string) Service {
	const errType = "Gateway Service"
	return &service{
		client: web.
			New().
			SetBaseURL(baseURL),
		routes:  make([]Route, 0),
		errType: errType,
	}
}

func (s *service) RegisterRoute(routes ...Route) Service {
	if len(routes) == 0 {
		return s
	}

	s.routes = append(s.routes, routes...)
	return s
}

func (s *service) Routes() []Route {
	return s.routes
}

func (s *service) Match(method string, path string) (Route, bool) {
	method = strings.ToLower(method)
	for _, r := range s.routes {
		if method != "any" && r.Method() != method {
			continue
		}

		if r.CatchPath() == path {
			return r, true
		}
	}

	return nil, false
}

func (s *service) Proxy(ctx context.Context, r Route, request Request) (_ Response, err error) {
	defer errs.Wrap(s.errType, &err, "Proxy")

	writer := web.NewBytesWriter()
	if request.RequestBody() != nil {
		_, err = writer.Write(request.RequestBody())
		if err != nil {
			return nil, err
		}
	}

	var response *web.Response
	response, err = s.client.
		R(ctx).
		Headers(request.Headers()).
		Cookies(request.Cookies()).
		Do(r.Method(), r.RedirectPath(), writer)
	if err != nil {
		return nil, err
	}

	return newResponse(response), nil
}
