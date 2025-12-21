package slices

// Reduce applies a function against an accumulator and each element to reduce it to a single value
func Reduce[T any, U any](slice []T, initial U, fn func(U, T) U) U {
	result := initial
	for _, v := range slice {
		result = fn(result, v)
	}
	return result
}

// Sum returns the sum of all numeric elements in a slice
func Sum[T int | int64 | float64](slice []T) T {
	var sum T
	for _, v := range slice {
		sum += v
	}
	return sum
}

// GroupBy groups slice elements by a key function
func GroupBy[T any, K comparable](slice []T, keyFn func(T) K) map[K][]T {
	groups := make(map[K][]T)
	for _, v := range slice {
		key := keyFn(v)
		groups[key] = append(groups[key], v)
	}
	return groups
}
