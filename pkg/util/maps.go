package util

import (
	"sort"
	"sync"
)

// SortedMapValues returns the values of a string-keyed map sorted by key.
// Acquires a read lock on mu for thread safety.
func SortedMapValues[V any](mu *sync.RWMutex, m map[string]V) []V {
	mu.RLock()
	defer mu.RUnlock()

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	result := make([]V, 0, len(m))
	for _, k := range keys {
		result = append(result, m[k])
	}
	return result
}

// SortedMapKeys returns the keys of a string-keyed map in sorted order.
// Acquires a read lock on mu for thread safety.
func SortedMapKeys[V any](mu *sync.RWMutex, m map[string]V) []string {
	mu.RLock()
	defer mu.RUnlock()

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
