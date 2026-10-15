package db

const (
	MIN_CANDIDATE_LIMIT = 2048
)

func searchCandidateLimit(limit int) int {
	return max(limit*4, MIN_CANDIDATE_LIMIT)
}

func ConcatArray[T any](arrays ...[]T) []T {
	var result []T
	for _, array := range arrays {
		if array != nil {
			result = append(result, array...)
		}
	}
	return result
}
