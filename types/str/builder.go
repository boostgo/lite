package str

import (
	"fmt"
	"strings"

	"github.com/boostgo/collection/slicex"
)

type Builder struct {
	builder *strings.Builder
}

func NewBuilder(grow ...int) *Builder {
	builder := &Builder{
		builder: &strings.Builder{},
	}

	if len(grow) > 0 {
		builder.Grow(grow[0])
	}

	return builder
}

func (builder *Builder) Grow(n int) *Builder {
	if n <= 0 {
		return builder
	}

	builder.builder.Grow(n)
	return builder
}

func (builder *Builder) Reset() *Builder {
	builder.builder.Reset()
	return builder
}

func (builder *Builder) Cap() int {
	return builder.builder.Cap()
}

func (builder *Builder) Len() int {
	return builder.builder.Len()
}

func (builder *Builder) String() string {
	return builder.builder.String()
}

func (builder *Builder) Write(p []byte) (int, error) {
	return builder.builder.Write(p)
}

func (builder *Builder) WriteString(s ...string) *Builder {
	if len(s) == 0 {
		return builder
	}

	s = slicex.Filter(s, func(s string) bool {
		return s != ""
	})

	if len(s) == 0 {
		return builder
	}

	for _, p := range s {
		builder.builder.WriteString(p)
	}

	return builder
}

func (builder *Builder) WriteRune(r rune) *Builder {
	builder.builder.WriteRune(r)
	return builder
}

func (builder *Builder) WriteByte(b byte) *Builder {
	builder.builder.WriteByte(b)
	return builder
}

func (builder *Builder) WriteFormat(format string, args ...any) *Builder {
	_, _ = fmt.Fprintf(builder.builder, format, args...)
	return builder
}

func String(s ...string) string {
	var grow int
	for _, v := range s {
		grow += len(v)
	}

	return NewBuilder(grow).
		WriteString(s...).
		String()
}

func Format(format string, args ...any) string {
	return NewBuilder().
		WriteFormat(format, args...).
		String()
}
