package list

func AreUnique[T any](source []T, fn func(a, b T) bool) bool {
	return len(Unique(source, fn)) == len(source)
}

func AreUniqueComparable[T comparable](source []T) bool {
	return AreUnique(source, func(a, b T) bool {
		return a == b
	})
}

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

func AreEqualComparable[T comparable](source []T, against []T) bool {
	return AreEqual(source, against, func(a, b T) bool {
		return a == b
	})
}
