package sql

import "fmt"

// Arguments Helps manage query arguments count & their values
type Arguments struct {
	args    []any
	counter int
}

// NewArguments created instance of Arguments object
func NewArguments(args ...any) *Arguments {
	return &Arguments{
		args:    args,
		counter: len(args),
	}
}

// Add new argument and increment count.
// After adding argument use Number method to get "$number" string
func (a *Arguments) Add(arg any) *Arguments {
	a.args = append(a.args, arg)
	a.counter++
	return a
}

// AddMany adds new many arguments and return "($1, $2, $3...)" string
func (a *Arguments) AddMany(args ...any) string {
	if len(args) == 0 {
		return ""
	}

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

// Number returns current $number value as a string
func (a *Arguments) Number() string {
	return fmt.Sprintf("$%d", a.counter)
}

// Args return all provided arguments for executing query
func (a *Arguments) Args() []any {
	return a.args
}
