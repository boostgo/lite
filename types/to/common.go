package to

func Ptr[T any](value T) *T {
	return &value
}
