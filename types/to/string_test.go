package to

import (
	"reflect"
	"testing"
)

func TestBool(t *testing.T) {
	type args struct {
		value any
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "True",
			args: args{value: true},
			want: true,
		},
		{
			name: "False",
			args: args{value: false},
			want: false,
		},
		{
			name: "String True + uppercase",
			args: args{value: "TRUE"},
			want: true,
		},
		{
			name: "String False + uppercase",
			args: args{value: "FALSE"},
			want: false,
		},
		{
			name: "Integer True",
			args: args{value: 1},
			want: true,
		},
		{
			name: "Integer False",
			args: args{value: 0},
			want: false,
		},
		{
			name: "Float32 False",
			args: args{value: float32(0.0)},
			want: false,
		},
		{
			name: "Float32 True",
			args: args{value: float32(1.0)},
			want: true,
		},
		{
			name: "Float64 False",
			args: args{value: 0.0},
			want: false,
		},
		{
			name: "Float64 True",
			args: args{value: 1.0},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Bool(tt.args.value); got != tt.want {
				t.Errorf("Bool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytes(t *testing.T) {

	type args struct {
		value any
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "String",
			args: args{value: "hello world"},
			want: []byte("hello world"),
		},
		{
			name: "Integer",
			args: args{value: 1},
			want: []byte{49},
		},
		{
			name: "JSON",
			args: args{value: map[string]any{"foo": "bar"}},
			want: []byte{123, 34, 102, 111, 111, 34, 58, 34, 98, 97, 114, 34, 125},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Bytes(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytesFromString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "String",
			args: args{s: "hello world"},
			want: []byte("hello world"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BytesFromString(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BytesFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat32(t *testing.T) {
	type args struct {
		anyValue any
	}
	tests := []struct {
		name string
		args args
		want float32
	}{
		{
			name: "Float32",
			args: args{anyValue: float32(1.0)},
			want: float32(1.0),
		},
		{
			name: "Float64",
			args: args{anyValue: 1.0},
			want: 1.0,
		},
		{
			name: "String",
			args: args{anyValue: "1"},
			want: float32(1.0),
		},
		{
			name: "String with dot",
			args: args{anyValue: "1.2"},
			want: float32(1.2),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Float32(tt.args.anyValue); got != tt.want {
				t.Errorf("Float32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat64(t *testing.T) {
	type args struct {
		anyValue any
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "Float32",
			args: args{anyValue: float32(1.0)},
			want: 1.0,
		},
		{
			name: "Float64",
			args: args{anyValue: 1.0},
			want: 1.0,
		},
		{
			name: "String",
			args: args{anyValue: "1"},
			want: 1.0,
		},
		{
			name: "String with dot",
			args: args{anyValue: "1.2"},
			want: 1.2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Float64(tt.args.anyValue); got != tt.want {
				t.Errorf("Float64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt(t *testing.T) {
	type args struct {
		anyValue any
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Int",
			args: args{anyValue: 1},
			want: 1,
		},
		{
			name: "String",
			args: args{anyValue: "1"},
			want: 1,
		},
		{
			name: "Int32",
			args: args{anyValue: int32(1)},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Int(tt.args.anyValue); got != tt.want {
				t.Errorf("Int() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString(t *testing.T) {
	type args struct {
		value any
	}
	value := "hello world"
	ptrString := &value

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "String",
			args: args{value: "hello world"},
			want: "hello world",
		},
		{
			name: "Pointer String",
			args: args{value: ptrString},
			want: "hello world",
		},
		{
			name: "Nil value",
			args: args{value: nil},
			want: "",
		},
		{
			name: "Integer",
			args: args{value: 123},
			want: "123",
		},
		{
			name: "Float32",
			args: args{value: float32(1.2)},
			want: "1.200000",
		},
		{
			name: "Float64",
			args: args{value: 1.2},
			want: "1.2",
		},
		{
			name: "Map",
			args: args{value: map[string]any{"foo": "bar"}},
			want: "{\"foo\":\"bar\"}",
		},
		{
			name: "Slice",
			args: args{value: []int{1, 2, 3}},
			want: "[1,2,3]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := String(tt.args.value); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
