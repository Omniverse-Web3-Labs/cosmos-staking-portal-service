package otelutil

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("utils/otelutil")

func Tracer() trace.Tracer {
	return tracer
}

func StartSpan(ctx context.Context, spanName string, callback func(ctx context.Context, span trace.Span)) {
	ctx, span := tracer.Start(ctx, spanName)
	defer span.End()
	callback(ctx, span)
}

func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}
