package list

func Keys[K comparable, V any](provide map[K]V) []K {
	keys := make([]K, 0, len(provide))
	for key := range provide {
		keys = append(keys, key)
	}
	return keys
}

func Values[K comparable, V any](provide map[K]V) []V {
	values := make([]V, 0, len(provide))
	for _, value := range provide {
		values = append(values, value)
	}
	return values
}
