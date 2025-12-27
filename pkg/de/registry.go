package de

import (
	"fmt"
	"sort"
	"sync"
)

// AlgorithmMetadata contains information about a registered algorithm.
type AlgorithmMetadata struct {
	Name        string
	Description string
}

// Registry manages algorithm registrations.
type Registry struct {
	metadata map[string]AlgorithmMetadata
	mu       sync.RWMutex
}

// NewRegistry creates a new algorithm registry.
func NewRegistry() *Registry {
	return &Registry{
		metadata: make(map[string]AlgorithmMetadata),
	}
}

// Register adds an algorithm to the registry.
func (r *Registry) Register(name string, meta AlgorithmMetadata) {
	r.mu.Lock()
	defer r.mu.Unlock()

	meta.Name = name // Ensure name is set
	r.metadata[name] = meta
}

// IsSupported checks if an algorithm is registered.
func (r *Registry) IsSupported(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, ok := r.metadata[name]
	return ok
}

// Get returns metadata for a specific algorithm.
func (r *Registry) Get(name string) (AlgorithmMetadata, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	meta, ok := r.metadata[name]
	if !ok {
		return AlgorithmMetadata{}, fmt.Errorf("algorithm not found: %s", name)
	}
	return meta, nil
}

// List returns all registered algorithm names.
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.metadata))
	for name := range r.metadata {
		names = append(names, name)
	}

	sort.Strings(names)
	return names
}

// ListMetadata returns metadata for all registered algorithms.
func (r *Registry) ListMetadata() []AlgorithmMetadata {
	r.mu.RLock()
	defer r.mu.RUnlock()

	metas := make([]AlgorithmMetadata, 0, len(r.metadata))
	for _, meta := range r.metadata {
		metas = append(metas, meta)
	}

	// Sort by name for consistent ordering
	sort.Slice(metas, func(i, j int) bool {
		return metas[i].Name < metas[j].Name
	})

	return metas
}

// DefaultRegistry is the global algorithm registry.
var DefaultRegistry = NewRegistry()
