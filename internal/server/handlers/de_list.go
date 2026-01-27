package handlers

import (
	"context"

	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/nicholaspcr/GoDE/pkg/de"
	"github.com/nicholaspcr/GoDE/pkg/problems"
	"github.com/nicholaspcr/GoDE/pkg/variants"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ListSupportedAlgorithms returns the list of supported DE algorithms.
func (deh *deHandler) ListSupportedAlgorithms(
	ctx context.Context, _ *emptypb.Empty,
) (*api.ListSupportedAlgorithmsResponse, error) {
	return &api.ListSupportedAlgorithmsResponse{
		Algorithms: de.DefaultRegistry.List(),
	}, nil
}

// ListSupportedVariants returns the list of supported DE variants.
func (deh *deHandler) ListSupportedVariants(
	ctx context.Context, _ *emptypb.Empty,
) (*api.ListSupportedVariantsResponse, error) {
	metas := variants.DefaultRegistry.ListMetadata()
	apiVariants := make([]*api.Variant, len(metas))
	for i, meta := range metas {
		apiVariants[i] = &api.Variant{
			Name:        meta.Name,
			Description: meta.Description,
		}
	}
	return &api.ListSupportedVariantsResponse{Variants: apiVariants}, nil
}

// ListSupportedProblems returns the list of supported optimization problems.
func (deh *deHandler) ListSupportedProblems(
	ctx context.Context, _ *emptypb.Empty,
) (*api.ListSupportedProblemsResponse, error) {
	metas := problems.DefaultRegistry.ListMetadata()
	apiProblems := make([]*api.Problem, len(metas))
	for i, meta := range metas {
		apiProblems[i] = &api.Problem{
			Name:        meta.Name,
			Description: meta.Description,
		}
	}
	return &api.ListSupportedProblemsResponse{Problems: apiProblems}, nil
}
