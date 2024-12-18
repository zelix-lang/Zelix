package util

// AnyMatch returns true if any element in the list matches the predicate
func AnyMatch[T any](list []T, predicate func(T) bool) bool {
	for _, element := range list {
		if predicate(element) {
			return true
		}
	}

	return false
}
