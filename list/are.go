package list

// AreUnique compares slice for unique elements by provided func
func AreUnique[T any](source []T, fn func(a, b T) bool) bool {
	return len(Unique(source, fn)) == len(source)
}

// AreUniqueComparable compares slice of comparable types for unique elements.
//
// Uses AreUnique function with default provided func
func AreUniqueComparable[T comparable](source []T) bool {
	return AreUnique(source, func(a, b T) bool {
		return a == b
	})
}

// AreEqual compares two slices by using provided func
func AreEqual[T any](source []T, against []T, fn func(T, T) bool) bool {
	if len(source) != len(against) {
		return false
	}

	for idx := range source {
		if !fn(source[idx], against[idx]) {
			return false
		}
	}

	return true
}

// AreEqualComparable compares two slices of comparable types.
//
// Uses AreEqual function with default provided func
func AreEqualComparable[T comparable](source []T, against []T) bool {
	return AreEqual(source, against, func(a, b T) bool {
		return a == b
	})
}
