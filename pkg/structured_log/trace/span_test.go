package trace

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSmallUnit_CreateChildSpan(t *testing.T) {

	ctxRoot, trace, close := NewTracerAndRootSpan(context.Background(), "test/smallUnit/CreateChildSpan")
	defer close()
	ctxChildA, close := WithSpan(ctxRoot, "test/smallUnit/CreateChildSpan/childSpanA")
	defer close()
	ctxChildB, close := WithSpan(ctxRoot, "test/smallUnit/CreateChildSpan/childSpanB")
	defer close()
	ctxChildAA, close := WithSpan(ctxChildA, "test/smallUnit/CreateChildSpan/childSpanA/childSpanA")
	defer close()

	testCases := []struct {
		Span               *Span
		ExpectedSpanName   string
		ExpectedParentSpan *Span
		ExpectedTrace      *Tracer
	}{
		{
			Span:               SpanFromContext(ctxRoot),
			ExpectedSpanName:   "test/smallUnit/CreateChildSpan",
			ExpectedParentSpan: nil,
			ExpectedTrace:      trace,
		},
		{
			Span:               SpanFromContext(ctxChildA),
			ExpectedSpanName:   "test/smallUnit/CreateChildSpan/childSpanA",
			ExpectedParentSpan: SpanFromContext(ctxRoot),
			ExpectedTrace:      trace,
		},
		{
			Span:               SpanFromContext(ctxChildB),
			ExpectedSpanName:   "test/smallUnit/CreateChildSpan/childSpanB",
			ExpectedParentSpan: SpanFromContext(ctxRoot),
			ExpectedTrace:      trace,
		},
		{
			Span:               SpanFromContext(ctxChildAA),
			ExpectedSpanName:   "test/smallUnit/CreateChildSpan/childSpanA/childSpanA",
			ExpectedParentSpan: SpanFromContext(ctxChildA),
			ExpectedTrace:      trace,
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("ExpectedSpanName=%s", testCase.ExpectedSpanName), func(t *testing.T) {
			if assert.NotNil(t, testCase.Span) {
				assert.Equal(t, testCase.ExpectedSpanName, testCase.Span.Name())
				assert.Equal(t, testCase.ExpectedParentSpan, testCase.Span.Parent())
				assert.Equal(t, testCase.ExpectedTrace, testCase.Span.tracer)
			}
		})
	}
}
