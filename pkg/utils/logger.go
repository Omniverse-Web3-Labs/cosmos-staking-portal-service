package utils

import (
	"app/pkg/utils/otelutil"
	"context"
	"log/slog"
)

var Logger logger

type logger struct {
}

func (*logger) Debug(ctx context.Context, module string, funcName string, message string, attributes ...any) {
	attrs := make([]any, 0)
	attrs = append(attrs, slog.Attr{
		Key:   "module",
		Value: slog.StringValue(module),
	})

	attrs = append(attrs, slog.Attr{
		Key:   "funcName",
		Value: slog.StringValue(funcName),
	})

	attrs = append(attrs, attributes...)

	otelutil.Logger().Log(ctx, slog.LevelDebug, message, attrs...)
}

func (*logger) Info(ctx context.Context, module string, funcName string, message string, attributes ...any) {
	attrs := make([]any, 0)
	attrs = append(attrs, slog.Attr{
		Key:   "module",
		Value: slog.StringValue(module),
	})

	attrs = append(attrs, slog.Attr{
		Key:   "funcName",
		Value: slog.StringValue(funcName),
	})

	attrs = append(attrs, attributes...)

	otelutil.Logger().Log(ctx, slog.LevelInfo, message, attrs...)
}

func (*logger) Warn(ctx context.Context, module string, funcName string, message string, attributes ...any) {
	attrs := make([]any, 0)
	attrs = append(attrs, slog.Attr{
		Key:   "module",
		Value: slog.StringValue(module),
	})

	attrs = append(attrs, slog.Attr{
		Key:   "funcName",
		Value: slog.StringValue(funcName),
	})

	attrs = append(attrs, attributes...)

	otelutil.Logger().Log(ctx, slog.LevelWarn, message, attrs...)
}

func (*logger) Error(ctx context.Context, module string, funcName string, message string, attributes ...any) {
	attrs := make([]any, 0)
	attrs = append(attrs, slog.Attr{
		Key:   "module",
		Value: slog.StringValue(module),
	})

	attrs = append(attrs, slog.Attr{
		Key:   "funcName",
		Value: slog.StringValue(funcName),
	})

	attrs = append(attrs, attributes...)

	otelutil.Logger().Log(ctx, slog.LevelError, message, attrs...)
}
