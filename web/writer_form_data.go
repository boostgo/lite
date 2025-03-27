package web

import (
	"bytes"
	"io"
	"mime/multipart"

	"github.com/boostgo/convert"
	"github.com/boostgo/lite/errs"
)

// FormDataWriter uses for sending form-data request body
type FormDataWriter interface {
	Add(key string, value any) error
	AddFile(name, fileName string, file []byte) error
	Set(data map[string]any) error
	Boundary() string
	ContentType() string
	Buffer() *bytes.Buffer
	Close() error
}

type formData struct {
	body    bytes.Buffer
	writer  *multipart.Writer
	errType string
}

// NewFormData creates FormDataWriter
func NewFormData(initial ...map[string]any) FormDataWriter {
	const errType = "Form Data Writer"

	fd := &formData{
		body:    bytes.Buffer{},
		errType: errType,
	}

	fd.writer = multipart.NewWriter(&fd.body)

	if len(initial) > 0 {
		fd.Set(initial[0])
	}

	return fd
}

func (fd *formData) Add(key string, value any) (err error) {
	defer errs.Wrap(fd.errType, &err, "Add")
	return fd.writer.WriteField(key, convert.String(value))
}

func (fd *formData) AddFile(name, fileName string, file []byte) (err error) {
	defer errs.Wrap(fd.errType, &err, "AddFile")

	fileWriter, err := fd.writer.CreateFormFile(name, fileName)
	if err != nil {
		return err
	}

	_, err = io.Copy(fileWriter, bytes.NewReader(file))
	if err != nil {
		return err
	}

	return nil
}

func (fd *formData) Set(data map[string]any) (err error) {
	defer errs.Wrap(fd.errType, &err, "Set")

	if data == nil || len(data) == 0 {
		return nil
	}

	for key, value := range data {
		if err = fd.Add(key, value); err != nil {
			return err
		}
	}

	return nil
}

func (fd *formData) Boundary() string {
	return fd.writer.Boundary()
}

func (fd *formData) ContentType() string {
	return fd.writer.FormDataContentType()
}

func (fd *formData) Buffer() *bytes.Buffer {
	return &fd.body
}

func (fd *formData) Close() (err error) {
	defer errs.Wrap(fd.errType, &err, "Close")
	return fd.writer.Close()
}
