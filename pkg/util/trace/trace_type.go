package trace

import (
	"context"
	"time"

	apitrace "go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"
)

type traceSpan struct {
	spanContext apitrace.SpanContext
}

// Tracer returns tracer used to create this span. Tracer cannot be nil.
func (span traceSpan) Tracer() apitrace.Tracer {
	return nil
}

// End completes the span. No updates are allowed to span after it
// ends. The only exception is setting status of the span.
func (span traceSpan) End(options ...apitrace.SpanOption) {
	return
}

// AddEvent adds an event to the span.
func (span traceSpan) AddEvent(ctx context.Context, name string, attrs ...label.KeyValue) {
	return
}

// AddEventWithTimestamp adds an event with a custom timestamp
// to the span.
func (span traceSpan) AddEventWithTimestamp(ctx context.Context, timestamp time.Time, name string, attrs ...label.KeyValue) {
	return
}

// IsRecording returns true if the span is active and recording events is enabled.
func (span traceSpan) IsRecording() bool {
	return false
}

// RecordError records an error as a span event.
func (span traceSpan) RecordError(ctx context.Context, err error, opts ...apitrace.ErrorOption) {
	return
}

// SpanContext returns span context of the span. Returned SpanContext is usable
// even after the span ends.
func (span traceSpan) SpanContext() apitrace.SpanContext {
	return span.spanContext
}

// SetStatus sets the status of the span in the form of a code
// and a message.  SetStatus overrides the value of previous
// calls to SetStatus on the Span.
//
// The default span status is OK, so it is not necessary to
// explicitly set an OK status on successful Spans unless it
// is to add an OK message or to override a previous status on the Span.
func (span traceSpan) SetStatus(code codes.Code, msg string) {
	return
}

// SetName sets the name of the span.
func (span traceSpan) SetName(name string) {
	return
}

// SetAttributes set span attributes
func (span traceSpan) SetAttributes(kv ...label.KeyValue) {
	return
}

// SetAttribute set singular span attribute, with type inference.
func (span traceSpan) SetAttribute(k string, v interface{}) {
	return
}
