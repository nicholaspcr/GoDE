package variants

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateIndices(t *testing.T) {
	tests := []struct {
		name      string
		startInd  int
		NP        int
		r         []int
		wantErr   bool
		expectedR []int
	}{
		{
			name:     "success",
			startInd: 1,
			NP:       10,
			r:        make([]int, 4),
			wantErr:  false,
		},
		{
			name:     "insufficient population",
			startInd: 1,
			NP:       3,
			r:        make([]int, 4),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			random := rand.New(rand.NewSource(1))
			err := GenerateIndices(tt.startInd, tt.NP, tt.r, random)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			// Check for uniqueness
			seen := make(map[int]bool)
			for _, val := range tt.r {
				assert.False(t, seen[val], "should not have duplicate indices")
				seen[val] = true
			}
		})
	}
}
