package arrayutil

// Includes returns true if val is in the slice.
func Includes[T comparable](slice []T, val T) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

// IndexOf returns the index of the first occurrence of val in slice, or -1 if not found.
func IndexOf[T comparable](slice []T, val T) int {
	for i, item := range slice {
		if item == val {
			return i
		}
	}
	return -1
}

// Filter returns a new slice with elements matching the predicate.
func Filter[T any](slice []T, predicate func(T) bool) []T {
	result := make([]T, 0)
	for _, item := range slice {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// Map transforms a slice using the mapper function.
func Map[T any, U any](slice []T, mapper func(T) U) []U {
	result := make([]U, len(slice))
	for i, item := range slice {
		result[i] = mapper(item)
	}
	return result
}
