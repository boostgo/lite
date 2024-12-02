package format

import (
	"testing"
)

func TestAlpha(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "With digits",
			args: args{input: "Hello world 777"},
			want: "Hello world",
		},
		{
			name: "With symbols",
			args: args{input: "Hello world !!!"},
			want: "Hello world",
		},
		{
			name: "Spaces before",
			args: args{input: "    Hello world !!!"},
			want: "Hello world",
		},
		{
			name: "Spaces after",
			args: args{input: "Hello world    "},
			want: "Hello world",
		},
		{
			name: "Spaces middle",
			args: args{input: "Hello     world"},
			want: "Hello world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Alpha(tt.args.input); got != tt.want {
				t.Errorf("Alpha() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCode(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Simple",
			args: args{input: "Hello World"},
			want: "hello_world",
		},
		{
			name: "With digits",
			args: args{input: "Hello World 123"},
			want: "hello_world_123",
		},
		{
			name: "Spaces before",
			args: args{input: "   Hello World"},
			want: "hello_world",
		},
		{
			name: "Spaces after",
			args: args{input: "Hello World   "},
			want: "hello_world",
		},
		{
			name: "Spaces middle",
			args: args{input: "Hello    World"},
			want: "hello_world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Code(tt.args.input); got != tt.want {
				t.Errorf("Code() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCyrillic(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Simple",
			args: args{input: "привет"},
			want: "privet",
		},
		{
			name: "Uppercase",
			args: args{input: "ПРИВЕТ"},
			want: "privet",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Cyrillic(tt.args.input); got != tt.want {
				t.Errorf("Cyrillic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEveryTitle(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Lowercase",
			args: args{input: "hello world"},
			want: "Hello World",
		},
		{
			name: "Uppercase",
			args: args{input: "HELLO WORLD"},
			want: "Hello World",
		},
		{
			name: "Spaces before",
			args: args{input: "   hello world"},
			want: "Hello World",
		},
		{
			name: "Spaces after",
			args: args{input: "hello world   "},
			want: "Hello World",
		},
		{
			name: "Spaces middle",
			args: args{input: "hello    world"},
			want: "Hello World",
		},
		{
			name: "Save digits",
			args: args{input: "hello world 123"},
			want: "Hello World 123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EveryTitle(tt.args.input); got != tt.want {
				t.Errorf("EveryTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestName(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Spaces and digits",
			args: args{input: "  john 123  "},
			want: "John",
		},
		{
			name: "Spaces and digits (full name)",
			args: args{input: "  john 123  Smith  _123* "},
			want: "John Smith",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Name(tt.args.input); got != tt.want {
				t.Errorf("Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumeric(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "No Digits",
			args: args{input: "Hello World"},
			want: "",
		},
		{
			name: "Only digits",
			args: args{input: "123"},
			want: "123",
		},
		{
			name: "Sentence with digits",
			args: args{input: "Hello World123"},
			want: "123",
		},
		{
			name: "Sentence with digits (digits around all words)",
			args: args{input: "1Hello2World3"},
			want: "123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Numeric(tt.args.input); got != tt.want {
				t.Errorf("Numeric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTitle(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "All lowercase",
			args: args{input: "hello world"},
			want: "Hello world",
		},
		{
			name: "All uppercase",
			args: args{input: "HELLO WORLD"},
			want: "Hello world",
		},
		{
			name: "Spaces before",
			args: args{input: "    hello world"},
			want: "Hello world",
		},
		{
			name: "Spaces after",
			args: args{input: "hello world    "},
			want: "Hello world",
		},
		{
			name: "Spaces middle",
			args: args{input: "hello    world"},
			want: "Hello world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Title(tt.args.input); got != tt.want {
				t.Errorf("Title() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_clearInput(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Before spaces",
			args: args{input: "     hello world"},
			want: "hello world",
		},
		{
			name: "After spaces",
			args: args{input: "hello world     "},
			want: "hello world",
		},
		{
			name: "Middle spaces",
			args: args{input: "hello     world"},
			want: "hello world",
		},
		{
			name: "Many spaces",
			args: args{input: "     hello     world     "},
			want: "hello world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := clearInput(tt.args.input); got != tt.want {
				t.Errorf("clearInput() = %v, want %v", got, tt.want)
			}
		})
	}
}
