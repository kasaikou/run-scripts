package trace

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type contextSpanKey struct{}

var ctxSpanKey = contextSpanKey{}

// SpanFromContext gets the Span in context and returns nil if it is invalid or not found it.
//
// Normally, this function is not used and should only be used for testing.
func SpanFromContext(ctx context.Context) *Span {

	span, ok := ctx.Value(ctxSpanKey).(*Span)
	if ok {
		if span.IsValid() {
			return span
		}
	}

	return nil
}

func withSpan(ctx context.Context, span *Span) context.Context {
	return context.WithValue(ctx, ctxSpanKey, span)
}

// WithSpan createds a Span that is a child of the Span in ctx and adds it to context.
//
// If there is no Span in context, or if the Span in context is closed, it generates panic.
func WithSpan(ctx context.Context, name string, attrs ...slog.Attr) (c context.Context, closer func()) {

	parent := SpanFromContext(ctx)

	if !parent.IsValid() {
		panic("span is invalid or not set")
	} else if isClosed, beginAt, _ := parent.Span(); isClosed {
		panic(fmt.Sprintf("span '%s' (beginAt: %s), in the context has been already ended", parent.name, beginAt.Format(time.RFC3339)))
	}

	tracer := parent.tracer
	span := newSpan(tracer, parent, name, attrs...)

	return withSpan(ctx, span), span.close
}

// Span contains information about Span dependencies and execution periods for logging.
type Span struct {
	name    string
	id      uuid.UUID
	parent  *Span
	tracer  *Tracer
	beginAt time.Time
	endAt   sql.NullTime
	attrs   map[string]slog.Value
}

func newSpan(tracer *Tracer, parent *Span, name string, attrs ...slog.Attr) *Span {

	attrsMap := make(map[string]slog.Value)
	for _, attr := range attrs {
		attrsMap[attr.Key] = attr.Value
	}

	return &Span{
		name:    name,
		id:      uuid.Must(uuid.NewV7()),
		parent:  parent,
		tracer:  tracer,
		beginAt: time.Now(),
		endAt:   sql.NullTime{},
		attrs:   attrsMap,
	}
}

// IsValid indictes whether it is valid or not.
func (s *Span) IsValid() bool { return s != nil }

// Name returns the Span name.
func (s *Span) Name() string { return s.name }

// SpanID returns the Span ID.
func (s *Span) SpanID() string { return s.id.String() }

// TraceID returns the Trace ID.
func (s *Span) TraceID() string { return s.tracer.id.String() }

// Parent returns the parent Span.
func (s *Span) Parent() *Span { return s.parent }

// HasParent indicates whether this has parent Span or not.
func (s *Span) HasParent() bool { return s.parent != nil }

// Attr returns the Value associated with key and Span.
func (s *Span) Attr(key string) (v slog.Value, isExist bool) {
	v, isExist = s.attrs[key]
	return v, isExist
}

// Span returns whether it is closed or not, begin time and end time.
func (s *Span) Span() (isClosed bool, beginAt, endAt time.Time) {
	return s.endAt.Valid, s.beginAt, s.endAt.Time
}

func (s *Span) close() {
	s.endAt = sql.NullTime{Valid: true, Time: time.Now()}
}
