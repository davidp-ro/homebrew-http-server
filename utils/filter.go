package utils

func Filter[T any](s []T, predicate func(T) bool) []T {
	var result []T
	for _, value := range s {
		if predicate(value) {
			result = append(result, value)
		}
	}
	return result
}
