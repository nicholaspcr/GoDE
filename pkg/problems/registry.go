package problems

import (
	"fmt"
	"sort"
	"sync"
)

// ProblemFactory is a function that creates a new problem instance.
type ProblemFactory func(dim, objs int) (Interface, error)

// ProblemMetadata contains information about a registered problem.
type ProblemMetadata struct {
	Name        string
	Description string
	MinDim      int
	MaxDim      int
	NumObjs     int
	Category    string // "multi" or "many"
}

// Registry manages problem registrations and creation.
type Registry struct {
	factories map[string]ProblemFactory
	metadata  map[string]ProblemMetadata
	mu        sync.RWMutex
}

// NewRegistry creates a new problem registry.
func NewRegistry() *Registry {
	return &Registry{
		factories: make(map[string]ProblemFactory),
		metadata:  make(map[string]ProblemMetadata),
	}
}

// Register adds a problem factory to the registry.
func (r *Registry) Register(name string, factory ProblemFactory, meta ProblemMetadata) {
	r.mu.Lock()
	defer r.mu.Unlock()

	meta.Name = name // Ensure name is set
	r.factories[name] = factory
	r.metadata[name] = meta
}

// Create instantiates a problem by name.
func (r *Registry) Create(name string, dim, objs int) (Interface, error) {
	r.mu.RLock()
	factory, ok := r.factories[name]
	r.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("problem does not exist")
	}

	return factory(dim, objs)
}

// List returns all registered problem names.
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.factories))
	for name := range r.factories {
		names = append(names, name)
	}

	sort.Strings(names)
	return names
}

// Get returns metadata for a specific problem.
func (r *Registry) Get(name string) (ProblemMetadata, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	meta, ok := r.metadata[name]
	return meta, ok
}

// ListMetadata returns metadata for all registered problems.
func (r *Registry) ListMetadata() []ProblemMetadata {
	r.mu.RLock()
	defer r.mu.RUnlock()

	metas := make([]ProblemMetadata, 0, len(r.metadata))
	for _, meta := range r.metadata {
		metas = append(metas, meta)
	}

	// Sort by name for consistent ordering
	sort.Slice(metas, func(i, j int) bool {
		return metas[i].Name < metas[j].Name
	})

	return metas
}

// DefaultRegistry is the global problem registry.
var DefaultRegistry = NewRegistry()
