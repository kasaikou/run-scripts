package trace

import (
	"fmt"
	"log/slog"
	"runtime"
	"time"
)

// TraceLog contains the information issued to the log.
type TraceLog struct {
	Span         *Span
	Level        slog.Level
	Message      string
	JSONAttrs    []byte
	programCount uintptr
	CreatedAt    time.Time
}

// Source indicates the location on the source file where the log was issued.
func (tl TraceLog) Source() string {
	fs := runtime.CallersFrames([]uintptr{tl.programCount})
	f, _ := fs.Next()
	return fmt.Sprintf("%s:%d", f.File, f.Line)
}
