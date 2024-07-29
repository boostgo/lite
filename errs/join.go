package errs

import "strings"

type joinErrors struct {
	errors []error
}

func Join(errors ...error) error {
	return &joinErrors{
		errors: errors,
	}
}

func (je joinErrors) Error() string {
	message := strings.Builder{}
	for i := 0; i < len(je.errors); i++ {
		message.WriteString(je.errors[i].Error())
		if i < len(je.errors)-1 {
			message.WriteString("\n")
		}
	}
	return message.String()
}
