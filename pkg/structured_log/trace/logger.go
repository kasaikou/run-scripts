package trace

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/kasaikou/markflow/pkg/util"
)

func must[V any](v V, err error) V {
	if err != nil {
		panic(err.Error())
	}

	return v
}

type groupOrAttrs struct {
	group string
	attrs []slog.Attr
}

// TracingHandler is a log handler for tracing using Span in context.
type TracingHandler struct {
	senderCh     chan<- TraceLog
	Enable       slog.Level
	groupOrAttrs []groupOrAttrs
}

// NewTracingHandler creates an instance of TracingHandler, [slog.Handler] associated with Tracer.
//
// TracingHandler associates logs with Span by getting information about Span from context.
func NewTracingHandler(level slog.Level, dest chan<- TraceLog) *TracingHandler {
	return &TracingHandler{
		senderCh: dest,
		Enable:   level,
	}
}

// Enabled indicates whether or not to output as logging.
func (th *TracingHandler) Enabled(ctx context.Context, level slog.Level) bool {

	span := SpanFromContext(ctx)
	if !span.IsValid() {
		return false
	}

	return level >= th.Enable
}

func appendJSONAttr(buf []byte, attr slog.Attr) []byte {

	attr.Value = attr.Value.Resolve()
	if attr.Equal(slog.Attr{}) {
		return buf
	}

	buf = append(buf, must(json.Marshal(attr.Key))...)
	buf = append(buf, ':')

	switch attr.Value.Kind() {
	case slog.KindGroup:
		attrs := attr.Value.Group()
		buf = appendJSONAttrs(buf, attrs)
	default:
		anyVal := attr.Value.Any()
		switch v := anyVal.(type) {
		case error:
			anyVal = v.Error()
		case time.Time:
			anyVal = v.Format(time.RFC3339Nano)
		case fmt.Stringer:
			anyVal = v.String()
		}

		buf = append(buf, must(json.Marshal(anyVal))...)
	}

	return buf
}

func appendJSONAttrs(buf []byte, attrs []slog.Attr) []byte {

	buf = append(buf, '{')
	for i, attr := range attrs {
		if i > 0 {
			buf = append(buf, ',')
		}
		buf = appendJSONAttr(buf, attr)
	}
	buf = append(buf, '}')

	return buf
}

var attrsPool = util.NewSlicePool[[]slog.Attr]()

// Handle generates a log from the record and context.
func (th *TracingHandler) Handle(ctx context.Context, record slog.Record) error {

	span := SpanFromContext(ctx)
	if !span.IsValid() {
		return ErrSpanIsInvalid
	}

	attrsPtr := attrsPool.NewSlice()
	groupOrAttrs := th.groupOrAttrs

	if record.NumAttrs() == 0 {
		for len(groupOrAttrs) > 0 && groupOrAttrs[len(groupOrAttrs)-1].group != "" {
			groupOrAttrs = groupOrAttrs[:len(groupOrAttrs)-1]
		}
	}

	for _, groupOrAttr := range groupOrAttrs {
		if groupOrAttr.group != "" {
			panic("cannot support WithGroup expression")
		} else {
			*attrsPtr = append(*attrsPtr, groupOrAttr.attrs...)
		}
	}

	record.Attrs(func(a slog.Attr) bool {
		*attrsPtr = append(*attrsPtr, a)
		return true
	})

	jsonAttrsPtr := util.NewBytes()
	*jsonAttrsPtr = appendJSONAttrs(*jsonAttrsPtr, *attrsPtr)

	th.senderCh <- TraceLog{
		Span:         span,
		Level:        record.Level,
		Message:      record.Message,
		JSONAttrs:    *jsonAttrsPtr,
		programCount: record.PC,
		CreatedAt:    record.Time,
	}

	attrsPool.ReleaseSlice(attrsPtr)
	return nil
}

func (th *TracingHandler) withGroupOrAttrs(goa groupOrAttrs) *TracingHandler {
	handler := *th
	handler.groupOrAttrs = make([]groupOrAttrs, 0, len(th.groupOrAttrs)+1)
	handler.groupOrAttrs = append(handler.groupOrAttrs, th.groupOrAttrs...)
	handler.groupOrAttrs = append(handler.groupOrAttrs, goa)

	return &handler
}

// WithAttrs returns a log handler containing attributes.
func (th *TracingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return th
	}

	return th.withGroupOrAttrs(groupOrAttrs{attrs: attrs})
}

// WithGroup returns a log handler containing group.
func (tl *TracingHandler) WithGroup(name string) slog.Handler {
	panic("WithGroup not implemented")
}
