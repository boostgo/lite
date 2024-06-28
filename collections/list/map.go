package list

func MapKeys[K comparable, V any](provide map[K]V) []K {
	keys := make([]K, 0, len(provide))
	for key, _ := range provide {
		keys = append(keys, key)
	}
	return keys
}
