package trace

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSmallUnit_RootSpanNoAttrsExample watches that
// log information is output correctly with no attributes.
func TestSmallUnit_RootSpanNoAttrsExample(t *testing.T) {

	const (
		expectedLevel        = slog.LevelInfo
		expectedSpanName     = "test/smallUnit/NoAttrsExample"
		expectedSourceSuffix = "/pkg/structured_log/trace/logger_test.go:33"
		expectedMessage      = "Test log."
		expectedJSONAttrs    = "{}"
	)

	queue := make(chan TraceLog, 1)
	ctx, _, close := NewTracerAndRootSpan(context.Background(), "test/smallUnit/NoAttrsExample")
	defer close()
	logger := slog.New(NewTracingHandler(slog.LevelInfo, queue))
	logger.InfoContext(ctx, expectedMessage)
	require.Equal(t, 1, len(queue))

	result := <-queue

	assert.Equal(t, expectedLevel, result.Level)
	assert.Equal(t, expectedSpanName, result.Span.Name())
	assert.Equal(t, expectedMessage, result.Message)
	assert.Equal(t, expectedJSONAttrs, string(result.JSONAttrs))
	assert.Truef(t, strings.HasSuffix(result.Source(), expectedSourceSuffix), "result.Source() '%s' do not have expected suffix '%s'", result.Source(), expectedSourceSuffix)
}

func toAnySlice[S []T, T any](from S) []any {
	s := make([]any, 0, len(from))
	for _, f := range from {
		s = append(s, f)
	}
	return s
}

// TestSmallUnit_RootSpanAttrs watches that
// log information is output correctly with some attributes.
func TestSmallUnit_RootSpanAttrs(t *testing.T) {

	t.Parallel()

	const (
		expectedLevel    = slog.LevelInfo
		expectedSpanName = "test/smallUnit/Attrs"
		expectedMessage  = "Test log with attributes."
	)

	testCases := []struct {
		loggerWith        func(logger *slog.Logger) *slog.Logger
		attrs             []slog.Attr
		expectedJSONAttrs string
	}{
		{
			loggerWith:        func(logger *slog.Logger) *slog.Logger { return logger },
			attrs:             []slog.Attr{},
			expectedJSONAttrs: "{}",
		},
		{
			loggerWith:        func(logger *slog.Logger) *slog.Logger { return logger.With(slog.String("key", "use logger.With()")) },
			attrs:             []slog.Attr{},
			expectedJSONAttrs: `{"key":"use logger.With()"}`,
		},
		{
			loggerWith:        func(logger *slog.Logger) *slog.Logger { return logger },
			attrs:             []slog.Attr{slog.String("key", "use attrs")},
			expectedJSONAttrs: `{"key":"use attrs"}`,
		},
		{
			loggerWith:        func(logger *slog.Logger) *slog.Logger { return logger },
			attrs:             []slog.Attr{slog.Any("error", errors.New("test error"))},
			expectedJSONAttrs: `{"error":"test error"}`,
		},
		{
			loggerWith:        func(logger *slog.Logger) *slog.Logger { return logger },
			attrs:             []slog.Attr{slog.Any("uuid", uuid.MustParse("6b02b9af-0e3f-4f49-83fd-39ab2eb38efc"))},
			expectedJSONAttrs: `{"uuid":"6b02b9af-0e3f-4f49-83fd-39ab2eb38efc"}`,
		},
		{
			loggerWith:        func(logger *slog.Logger) *slog.Logger { return logger },
			attrs:             []slog.Attr{slog.String("primary", "primary"), slog.String("secondary", "secondary")},
			expectedJSONAttrs: `{"primary":"primary","secondary":"secondary"}`,
		},
	}

	ctx, _, close := NewTracerAndRootSpan(context.Background(), expectedSpanName)
	defer close()
	wg := sync.WaitGroup{}

	for _, testCase := range testCases {
		testCase := testCase

		wg.Add(1)
		t.Run(fmt.Sprintf("expectedJSONAttrs=%s", testCase.expectedJSONAttrs), func(t *testing.T) {
			defer wg.Done()

			queue := make(chan TraceLog, 1)
			logger := testCase.loggerWith(slog.New(NewTracingHandler(slog.LevelInfo, queue)))
			logger.InfoContext(ctx, expectedMessage, toAnySlice(testCase.attrs)...)

			result := <-queue

			assert.Equal(t, expectedLevel, result.Level)
			assert.Equal(t, expectedSpanName, result.Span.Name())
			assert.Equal(t, expectedMessage, result.Message)
			assert.Equal(t, testCase.expectedJSONAttrs, string(result.JSONAttrs))
		})
	}

	defer wg.Wait()
}
