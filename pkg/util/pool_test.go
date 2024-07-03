package util

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// BenchmarkPoolIsNotRecycled shows speed of (*SlicePool).NewSlice and (*SlicePool).ReleaseSlice.
//
// Used to verify that the allocation count (per opts) is 0.
func BenchmarkPoolIsRecycled(b *testing.B) {

	bp := NewSlicePool[[]int]()

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			ints := bp.NewSlice()
			*ints = append(*ints, 100)
			bp.ReleaseSlice(ints)
		}
	})
}

// TestSmallUnit_IsAllocatedSlice watches that
// the length of slice as 0 when there is already an unused slice in the pool.
func TestSmallUnit_IsRecyecledIsResetted(t *testing.T) {

	const (
		Length      = 100
		expectedLen = 0
		expectedCap = Length
	)

	// Initialize
	bp := NewSlicePool[[]int]()
	used := make([]int, Length)
	bp.ReleaseSlice(&used)

	// Get slice pointer
	ints := bp.NewSlice()
	assert.Equal(t, expectedLen, len(*ints))
	assert.Equal(t, expectedCap, cap(*ints))
}

// TestSmallUnit_AllocatedSlice watches that
// the capacity of slice, created by NewSlicePoolAllocate varies with cap argument,
// and that length of slice is 0.
func TestSmallUnit_AllocatedSlice(t *testing.T) {

	for i := 0; i < 5; i++ {

		expectedCap := int(math.Pow10(i))

		t.Run(fmt.Sprintf("cap=%d", expectedCap), func(t *testing.T) {

			t.Parallel()

			const expectedLen = 0

			// Initialize
			bp := NewSlicePoolAllocate[[]int](expectedCap)

			// Get slice pointer
			ints := bp.NewSlice()

			assert.Equal(t, expectedLen, len(*ints))
			assert.Equal(t, expectedCap, cap(*ints))
		})
	}
}
