package errs

import "strings"

type joinErrors struct {
	errors []error
}

// Join all provided errors into one error
func Join(errors ...error) error {
	errsList := make([]error, 0, len(errors))
	for _, err := range errors {
		errsList = append(errsList, err)
	}
	return &joinErrors{
		errors: errsList,
	}
}

// Unwrap return all errors slice
func (je joinErrors) Unwrap() []error {
	if je.errors == nil {
		je.errors = make([]error, 0)
	}

	return je.errors
}

// Error join all errors into one string
func (je joinErrors) Error() string {
	if je.errors == nil || len(je.errors) == 0 {
		return ""
	}

	message := strings.Builder{}
	for i := 0; i < len(je.errors); i++ {
		message.WriteString(je.errors[i].Error())
		if i < len(je.errors)-1 {
			message.WriteString(" - ")
		}
	}
	return message.String()
}
