package sql

import "fmt"

type Arguments struct {
	args    []any
	counter int
}

// NewArguments created instance of Arguments object.
// Helps manage query arguments count & their values
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

func (a *Arguments) AddMany(args ...any) string {
	values := "("
	for idx, arg := range args {
		values += a.Add(arg).Number()
		if idx < len(args)-1 {
			values += ", "
		}
	}
	values += ")"
	return values
}

func (a *Arguments) Number() string {
	return fmt.Sprintf("$%d", a.counter)
}

func (a *Arguments) Args() []any {
	return a.args
}
