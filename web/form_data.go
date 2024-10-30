package web

import (
	"bytes"
	"github.com/boostgo/lite/types/to"
	"io"
	"mime/multipart"
)

type FormDataWriter interface {
	Add(key string, value any) FormDataWriter
	AddFile(name, fileName string, file []byte) FormDataWriter
	Boundary() string
	ContentType() string
	Buffer() *bytes.Buffer
	Close() error
}

type formData struct {
	body   bytes.Buffer
	writer *multipart.Writer
}

func NewFormData() FormDataWriter {
	fd := &formData{
		body: bytes.Buffer{},
	}

	fd.writer = multipart.NewWriter(&fd.body)
	return fd
}

func (fd *formData) Add(key string, value any) FormDataWriter {
	_ = fd.writer.WriteField(key, to.String(value))
	return fd
}

func (fd *formData) AddFile(name, fileName string, file []byte) FormDataWriter {
	fileWriter, err := fd.writer.CreateFormFile(name, fileName)
	if err != nil {
		return fd
	}

	_, _ = io.Copy(fileWriter, bytes.NewReader(file))
	return fd
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

func (fd *formData) Close() error {
	return fd.writer.Close()
}
