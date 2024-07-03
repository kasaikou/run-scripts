package structured_log

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func must[V any](v V, err error) V {
	if err != nil {
		panic(err)
	}
	return v
}

// TestSmallUnit_RouteReaderLoggerPipeSplitsWithLF watches that
// RouteReaderLoggerPipe splits the string line by line and outputs it to the logger.
func TestSmallUnit_RouteReaderLoggerPipeSplitsWithLF(t *testing.T) {

	testCases := []struct {
		SourceString    string
		ExpectedJsonMap []map[string]any
	}{
		{
			SourceString: "abcdefg",
			ExpectedJsonMap: []map[string]any{
				{
					"level": "INFO",
					"msg":   "abcdefg",
				},
			},
		}, {
			SourceString: "abcd\nefgh\n",
			ExpectedJsonMap: []map[string]any{
				{
					"level": "INFO",
					"msg":   "abcd",
				}, {
					"level": "INFO",
					"msg":   "efgh",
				},
			},
		}, {
			SourceString: "abcd\nefgh",
			ExpectedJsonMap: []map[string]any{
				{
					"level": "INFO",
					"msg":   "abcd",
				}, {
					"level": "INFO",
					"msg":   "efgh",
				},
			},
		}, {
			SourceString:    "",
			ExpectedJsonMap: []map[string]any{},
		},
	}

	for _, testCase := range testCases {

		testCase := testCase
		t.Run(fmt.Sprintf("parseText=%s", must(json.Marshal(testCase.SourceString))), func(t *testing.T) {

			t.Parallel()

			// Initialize resources.
			reader := bytes.NewBufferString(testCase.SourceString)
			buffer := bytes.NewBuffer([]byte{})
			logger := slog.New(slog.NewJSONHandler(buffer, &slog.HandlerOptions{}))

			// Logging from io.Reader.
			RouteReaderLoggerPipe(context.Background(), slog.LevelInfo, reader, logger)

			// Read JSON line format.
			scanner := bufio.NewScanner(bytes.NewReader(buffer.Bytes()))
			results := make([]map[string]any, 0, len(testCase.ExpectedJsonMap))

			for scanner.Scan() {
				dest := map[string]any{}
				assert.NoError(t, json.Unmarshal(scanner.Bytes(), &dest))

				// Intentionally remove time: not needed for testing as it is included in the output in slog.Logger.
				delete(dest, "time")

				results = append(results, dest)
			}

			assert.EqualValues(t, testCase.ExpectedJsonMap, results)
		})
	}

}
