package web

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/boostgo/lite/system/trace"
	"github.com/boostgo/lite/types/flex"
	"net/http"
)

type Response struct {
	raw      *http.Response
	bodyBlob []byte
}

func newResponse(resp *http.Response) *Response {
	return &Response{
		raw: resp,
	}
}

func (response Response) Raw() *http.Response {
	return response.raw
}

func (response Response) Status() string {
	return response.raw.Status
}

func (response Response) StatusCode() int {
	return response.raw.StatusCode
}

func (response Response) BodyRaw() []byte {
	return response.bodyBlob
}

func (response Response) Parse(export any) error {
	if response.bodyBlob == nil {
		return nil
	}

	if !flex.Type(export).IsPtr() {
		return errors.New("provided export is not a pointer")
	}

	return json.Unmarshal(response.bodyBlob, export)
}

func (response Response) Context(ctx context.Context) context.Context {
	return trace.Set(ctx, trace.FromResponse(response.raw))
}

func (response Response) TraceID() string {
	return trace.FromResponse(response.raw)
}
