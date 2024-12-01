package gateway

type Route interface {
	Method() string
	CatchPath() string
	RedirectPath() string
}

type route struct {
	method       string
	catchPath    string
	redirectPath string
}

func NewRoute(method, catchPath, redirectPath string) Route {
	return &route{
		method:       method,
		catchPath:    catchPath,
		redirectPath: redirectPath,
	}
}

func (r *route) Method() string {
	return r.method
}

func (r *route) CatchPath() string {
	return r.catchPath
}

func (r *route) RedirectPath() string {
	return r.redirectPath
}
