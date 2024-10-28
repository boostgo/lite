package api

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
