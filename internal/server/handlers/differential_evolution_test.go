package handlers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/emptypb"

	_ "github.com/nicholaspcr/GoDE/pkg/problems/many/dtlz"       // Register DTLZ problems
	_ "github.com/nicholaspcr/GoDE/pkg/problems/many/wfg"        // Register WFG problems
	_ "github.com/nicholaspcr/GoDE/pkg/problems/multi"           // Register ZDT and VNT problems
	_ "github.com/nicholaspcr/GoDE/pkg/variants/best"            // Register best/* variants
	_ "github.com/nicholaspcr/GoDE/pkg/variants/current-to-best" // Register current-to-best/* variants
	_ "github.com/nicholaspcr/GoDE/pkg/variants/pbest"           // Register pbest/* variants
	_ "github.com/nicholaspcr/GoDE/pkg/variants/rand"            // Register rand/* variants
)

func TestDEHandler_ListSupportedAlgorithms(t *testing.T) {
	handler, _ := setupTestHandler()

	resp, err := handler.ListSupportedAlgorithms(context.Background(), &emptypb.Empty{})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, []string{"gde3"}, resp.Algorithms)
}

func TestDEHandler_ListSupportedVariants(t *testing.T) {
	handler, _ := setupTestHandler()

	resp, err := handler.ListSupportedVariants(context.Background(), &emptypb.Empty{})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Variants, 6)

	// Verify we have the expected variants
	variantNames := make(map[string]bool)
	for _, v := range resp.Variants {
		variantNames[v.Name] = true
		assert.NotEmpty(t, v.Description)
	}

	// Verify expected variants exist
	assert.Contains(t, variantNames, "rand1")
	assert.Contains(t, variantNames, "rand2")
	assert.Contains(t, variantNames, "best1")
	assert.Contains(t, variantNames, "best2")
	assert.Contains(t, variantNames, "pbest")
	assert.Contains(t, variantNames, "currToBest1")
}

func TestDEHandler_ListSupportedProblems(t *testing.T) {
	handler, _ := setupTestHandler()

	resp, err := handler.ListSupportedProblems(context.Background(), &emptypb.Empty{})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Problems, 22) // 6 ZDT/VNT + 7 DTLZ + 9 WFG

	// Verify we have the expected problem families
	problemNames := make(map[string]bool)
	for _, p := range resp.Problems {
		problemNames[p.Name] = true
		assert.NotEmpty(t, p.Description)
	}

	// Check for ZDT problems
	assert.Contains(t, problemNames, "zdt1")
	assert.Contains(t, problemNames, "zdt2")
	assert.Contains(t, problemNames, "zdt3")
	assert.Contains(t, problemNames, "zdt4")
	assert.Contains(t, problemNames, "zdt6")

	// Check for VNT problem
	assert.Contains(t, problemNames, "vnt1")

	// Check for DTLZ problems
	assert.Contains(t, problemNames, "dtlz1")
	assert.Contains(t, problemNames, "dtlz2")

	// Check for WFG problems
	assert.Contains(t, problemNames, "wfg1")
	assert.Contains(t, problemNames, "wfg2")
}
