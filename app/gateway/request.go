package gateway

type Request interface {
	RequestBody() []byte
	Headers() map[string]any
	Cookies() map[string]any
}

type gwRequest struct {
	requestBody []byte
	headers     map[string]any
	cookies     map[string]any
}

func NewRequest(requestBody []byte, headers, cookies map[string]any) Request {
	return &gwRequest{
		requestBody: requestBody,
		headers:     headers,
		cookies:     cookies,
	}
}

func (s *gwRequest) RequestBody() []byte {
	return s.requestBody
}

func (s *gwRequest) Headers() map[string]any {
	return s.headers
}

func (s *gwRequest) Cookies() map[string]any {
	return s.cookies
}
