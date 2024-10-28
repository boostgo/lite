package sql

import "fmt"

type Arguments struct {
	args    []any
	counter int
}

func NewArguments(args ...any) *Arguments {
	return &Arguments{
		args:    args,
		counter: len(args),
	}
}

func (a *Arguments) Add(arg any) *Arguments {
	a.args = append(a.args, arg)
	a.counter++
	return a
}

func (a *Arguments) Number() string {
	return fmt.Sprintf("$%d", a.counter)
}

func (a *Arguments) Args() []any {
	return a.args
}
