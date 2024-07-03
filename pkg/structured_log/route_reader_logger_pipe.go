package structured_log

import (
	"bufio"
	"context"
	"io"
	"log/slog"
)

func RouteReaderLoggerPipe(ctx context.Context, level slog.Level, src io.Reader, dest *slog.Logger) {
	routeReaderLoggerPipe(src, func(content []byte) { dest.Log(ctx, level, string(content)) })
}

func routeReaderLoggerPipe(src io.Reader, bytesLoggerFn func(content []byte)) {

	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		bytesLoggerFn(scanner.Bytes())
	}

	if err := scanner.Err(); err != nil {
		panic("routeReaderLoggerPipe() ended with non-EOL error: " + err.Error())
	}
}
