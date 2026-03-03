// Package util provides generic utility functions for common slice and map operations.
package util

// MapSlice transforms a slice of T into a slice of U using the provided mapper function.
// Returns an error if any mapper call fails.
func MapSlice[T, U any](items []T, mapper func(T) (U, error)) ([]U, error) {
	result := make([]U, len(items))
	for i, item := range items {
		mapped, err := mapper(item)
		if err != nil {
			return nil, err
		}
		result[i] = mapped
	}
	return result, nil
}

// Copyable constrains T to types that have a Copy() T method.
type Copyable[T any] interface {
	Copy() T
}

// CopySlice returns a new slice containing deep copies of all items.
func CopySlice[T Copyable[T]](items []T) []T {
	result := make([]T, len(items))
	for i, item := range items {
		result[i] = item.Copy()
	}
	return result
}
