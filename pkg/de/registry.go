package de

import (
	"fmt"
	"sort"
	"sync"

	api "github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/problems"
	"github.com/nicholaspcr/GoDE/pkg/variants"
)

// AlgorithmMetadata contains information about a registered algorithm.
type AlgorithmMetadata struct {
	Name        string
	Description string
}

// AlgorithmParams holds all parameters needed to instantiate an algorithm.
type AlgorithmParams struct {
	Problem           problems.Interface
	Variant           variants.Interface
	PopulationParams  models.PopulationParams
	InitialPopulation models.Population
	ProgressCallback  ProgressCallback
}

// AlgorithmFactory creates an Algorithm from execution parameters and config.
// The config provides algorithm-specific parameters (e.g. CR, F, P for GDE3).
type AlgorithmFactory func(params AlgorithmParams, config *api.DEConfig) (Algorithm, error)

// Registry manages algorithm registrations.
type Registry struct {
	metadata  map[string]AlgorithmMetadata
	factories map[string]AlgorithmFactory
	mu        sync.RWMutex
}

// NewRegistry creates a new algorithm registry.
func NewRegistry() *Registry {
	return &Registry{
		metadata:  make(map[string]AlgorithmMetadata),
		factories: make(map[string]AlgorithmFactory),
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

// RegisterFactory associates an AlgorithmFactory with a registered algorithm.
// The algorithm must be registered via Register before calling RegisterFactory.
func (r *Registry) RegisterFactory(name string, factory AlgorithmFactory) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.factories[name] = factory
}

// GetFactory returns the AlgorithmFactory for the named algorithm.
func (r *Registry) GetFactory(name string) (AlgorithmFactory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	factory, ok := r.factories[name]
	if !ok {
		return nil, fmt.Errorf("no factory registered for algorithm: %s", name)
	}
	return factory, nil
}

// DefaultRegistry is the global algorithm registry.
var DefaultRegistry = NewRegistry()
