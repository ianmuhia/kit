package slices

// All returns true if all elements satisfy the predicate
func All[T any](slice []T, predicate func(T) bool) bool {
	for _, v := range slice {
		if !predicate(v) {
			return false
		}
	}
	return true
}

// Any returns true if at least one element satisfies the predicate
func Any[T any](slice []T, predicate func(T) bool) bool {
	for _, v := range slice {
		if predicate(v) {
			return true
		}
	}
	return false
}

// Chunk splits a slice into chunks of specified size
func Chunk[T any](slice []T, size int) [][]T {
	if size <= 0 {
		return nil
	}

	chunks := make([][]T, 0, (len(slice)+size-1)/size)
	for i := 0; i < len(slice); i += size {
		end := i + size
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}

// Flatten flattens a slice of slices into a single slice
func Flatten[T any](slices [][]T) []T {
	result := make([]T, 0)
	for _, slice := range slices {
		result = append(result, slice...)
	}
	return result
}
