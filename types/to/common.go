package to

// Ptr convert any value to pointer to this value
func Ptr[T any](value T) *T {
	return &value
}
