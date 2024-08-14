package web

import "context"

var (
	Default = New()
)

func Get(ctx context.Context, url string, params ...any) (*Response, error) {
	return Default.
		R(ctx).
		GET(url, params...)
}

func Post(ctx context.Context, body any, url string) (*Response, error) {
	return Default.
		R(ctx).
		POST(url, body)
}

func Put(ctx context.Context, body any, url string) (*Response, error) {
	return Default.
		R(ctx).
		PUT(url, body)
}

func Delete(ctx context.Context, url string, params ...any) (*Response, error) {
	return Default.
		R(ctx).
		DELETE(url, params...)
}
