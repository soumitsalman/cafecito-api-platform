package db

func ConcatArray[T any](arrays ...[]T) []T {
	var result []T
	for _, array := range arrays {
		if array != nil {
			result = append(result, array...)
		}
	}
	return result
}
