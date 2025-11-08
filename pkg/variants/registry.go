package variants

import (
	"fmt"
	"sort"
	"sync"
)

// VariantFactory is a function that creates a new variant instance.
type VariantFactory func() Interface

// VariantMetadata contains information about a registered variant.
type VariantMetadata struct {
	Name        string
	Description string
	Category    string // e.g., "rand", "best", "current-to-best", "pbest"
}

// Registry manages variant registrations and creation.
type Registry struct {
	factories map[string]VariantFactory
	metadata  map[string]VariantMetadata
	mu        sync.RWMutex
}

// NewRegistry creates a new variant registry.
func NewRegistry() *Registry {
	return &Registry{
		factories: make(map[string]VariantFactory),
		metadata:  make(map[string]VariantMetadata),
	}
}

// Register adds a variant factory to the registry.
func (r *Registry) Register(name string, factory VariantFactory, meta VariantMetadata) {
	r.mu.Lock()
	defer r.mu.Unlock()

	meta.Name = name // Ensure name is set
	r.factories[name] = factory
	r.metadata[name] = meta
}

// Create instantiates a variant by name.
func (r *Registry) Create(name string) (Interface, error) {
	r.mu.RLock()
	factory, ok := r.factories[name]
	r.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("unknown variant: %s", name)
	}

	return factory(), nil
}

// List returns all registered variant names.
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

// Get returns metadata for a specific variant.
func (r *Registry) Get(name string) (VariantMetadata, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	meta, ok := r.metadata[name]
	return meta, ok
}

// ListMetadata returns metadata for all registered variants.
func (r *Registry) ListMetadata() []VariantMetadata {
	r.mu.RLock()
	defer r.mu.RUnlock()

	metas := make([]VariantMetadata, 0, len(r.metadata))
	for _, meta := range r.metadata {
		metas = append(metas, meta)
	}

	// Sort by name for consistent ordering
	sort.Slice(metas, func(i, j int) bool {
		return metas[i].Name < metas[j].Name
	})

	return metas
}

// DefaultRegistry is the global variant registry.
var DefaultRegistry = NewRegistry()
