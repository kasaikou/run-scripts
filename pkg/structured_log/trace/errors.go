package trace

import "errors"

var (
	// ErrSpanIsInvalid shows Span is invalid or not found, including nil.
	ErrSpanIsInvalid = errors.New("span is invalid or not found")
)
