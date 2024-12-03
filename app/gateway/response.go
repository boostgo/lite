package gateway

import "github.com/boostgo/lite/web"

// Response is gateway response
type Response interface {
	Body() []byte
	StatusCode() int
	ContentType() string
}

type gwResponse struct {
	response *web.Response
}

func newResponse(response *web.Response) Response {
	return &gwResponse{
		response: response,
	}
}

func (r *gwResponse) Body() []byte {
	return r.response.BodyRaw()
}

func (r *gwResponse) StatusCode() int {
	return r.response.StatusCode()
}

func (r *gwResponse) ContentType() string {
	return r.response.ContentType()
}
