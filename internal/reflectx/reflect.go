package reflectx

func IsPointer(x any) bool {
	if x == nil {
		return true
	}

	_, ok := x.(*any)
	return ok
}
