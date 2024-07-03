package trace

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// Tracer contains Span and the logs generated with in it.
type Tracer struct {
	name    string
	id      uuid.UUID
	spans   map[uuid.UUID]*Span
	beginAt time.Time
	endAt   sql.NullTime
}

// NewTracerAndRootSpan creates a Tracer and its immediate Span.
//
// The Span directly under the Tracer can only be created in this function.
//
// The created Span is stored in context.
//
// closer closes Tracer and Span created in this function.
func NewTracerAndRootSpan(ctx context.Context, name string, attrs ...slog.Attr) (c context.Context, tracer *Tracer, closer func()) {

	tracer = &Tracer{
		name:    name,
		id:      uuid.Must(uuid.NewV7()),
		spans:   make(map[uuid.UUID]*Span),
		beginAt: time.Now(),
		endAt:   sql.NullTime{},
	}
	span := newSpan(tracer, nil, name, attrs...)
	tracer.registerSpan(span)

	return withSpan(ctx, span), tracer, func() {
		span.close()
		tracer.endAt = sql.NullTime{Valid: true, Time: time.Now()}
	}
}

func (t *Tracer) registerSpan(span *Span) { t.spans[span.id] = span }
