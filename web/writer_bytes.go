package web

import (
	"bytes"
	"github.com/boostgo/lite/errs"
	"io"
)

type BytesWriter interface {
	io.Writer
	ContentType() string
	Bytes() []byte
	SetContentType(contentType string) BytesWriter
	Reader() io.Reader
}

type bytesBuffer struct {
	buffer      *bytes.Buffer
	contentType string
	errType     string
}

func NewBytesWriter() BytesWriter {
	const (
		defaultContentType = "application/octet-stream"
		errType            = "Bytes Writer"
	)

	return &bytesBuffer{
		buffer:      bytes.NewBuffer(make([]byte, 0)),
		contentType: defaultContentType,
		errType:     errType,
	}
}

func (writer *bytesBuffer) Write(bytes []byte) (n int, err error) {
	defer errs.Wrap(writer.errType, &err, "Write")
	return writer.buffer.Write(bytes)
}

func (writer *bytesBuffer) ContentType() string {
	return writer.contentType
}

func (writer *bytesBuffer) SetContentType(contentType string) BytesWriter {
	writer.contentType = contentType
	return writer
}

func (writer *bytesBuffer) Bytes() []byte {
	return writer.buffer.Bytes()
}

func (writer *bytesBuffer) Reader() io.Reader {
	return writer.buffer
}
