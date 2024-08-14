package web

import "context"

func Get(ctx context.Context, url string, params ...any) (*Response, error) {
	return R(ctx).
		GET(url, params...)
}

func Post(ctx context.Context, body any, url string) (*Response, error) {
	return R(ctx).
		POST(url, body)
}

func Put(ctx context.Context, body any, url string) (*Response, error) {
	return R(ctx).
		PUT(url, body)
}

func Delete(ctx context.Context, url string, params ...any) (*Response, error) {
	return R(ctx).
		DELETE(url, params...)
}
