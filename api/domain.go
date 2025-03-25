package api

import (
	"encoding/json"
	"github.com/boostgo/lite/errs"
)

type createdID struct {
	ID any `json:"id"`
}

func newCreatedID(id any) createdID {
	return createdID{
		ID: id,
	}
}

const (
	statusSuccess = "Success"
	statusFailure = "Failure"
)

type errorOutput struct {
	Status  string         `json:"status"`
	Type    string         `json:"type,omitempty"`
	Message string         `json:"message"`
	Inner   string         `json:"inner,omitempty"`
	Context map[string]any `json:"context,omitempty"`
}

func WrapErrorBlob(err error) []byte {
	const defaultErrorType = "ERROR"

	var output errorOutput
	output.Status = statusFailure

	// build/collect error output
	custom, ok := errs.TryGet(err)
	if ok {
		output.Message = custom.Message()
		output.Type = custom.Type()
		output.Context = custom.Context()
		if custom.InnerError() != nil {
			output.Inner = custom.InnerError().Error()
		}
	} else {
		output.Message = err.Error()
		output.Type = defaultErrorType
	}

	// clear from trace
	if output.Context != nil {
		if _, traceExist := output.Context["trace"]; traceExist {
			delete(output.Context, "trace")
		}
	}

	outputBlob, _ := json.Marshal(output)
	return outputBlob
}

type successOutput struct {
	Status string `json:"status"`
	Body   any    `json:"body"`
}

func newSuccess(body any) successOutput {
	return successOutput{
		Status: statusSuccess,
		Body:   body,
	}
}
