package util

import "sync"

// SlicePool provides a ppol of slices with elements of T.
type SlicePool[S []T, T any] struct {
	_    struct{}
	pool sync.Pool
}

// NewSlicePool initialize pool.SlicePool instance.
func NewSlicePool[S []T, T any]() *SlicePool[S, T] {
	return NewSlicePoolAllocate[S](0)
}

// NewSlicePoolAllocate initialize pool.SlicePool instance with cap parameter.
//
// cap is used to initialize the slice.
func NewSlicePoolAllocate[S []T, T any](cap int) *SlicePool[S, T] {
	return &SlicePool[S, T]{
		pool: sync.Pool{
			New: func() any {
				s := make(S, 0, cap)
				return &s
			},
		},
	}
}

// NewSlice generates or recycles slice.
func (bp *SlicePool[S, T]) NewSlice() *S {
	return bp.pool.Get().(*S)
}

// ReleaseSlice stores unused slice.
func (bp *SlicePool[S, T]) ReleaseSlice(s *S) {
	(*s) = (*s)[:0]
	bp.pool.Put(s)
}

var bytesPool = NewSlicePool[[]byte]()

// NewBytes generates or recycles byte slice using the SlicePool[[]byte] shared within internal package.
func NewBytes() *[]byte { return bytesPool.NewSlice() }

// ReleaseBytes stores unused slice to SlicePool[[]byte] shared within internal package.
func ReleaseBytes(b *[]byte) { bytesPool.ReleaseSlice(b) }
