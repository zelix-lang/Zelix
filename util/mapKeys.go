package util

// MapKeys returns the keys of a map
func MapKeys[K comparable, V any](list map[K]V) []K {
	result := make([]K, 0)

	for key := range list {
		result = append(result, key)
	}

	return result
}
