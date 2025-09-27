package otelutil

import (
	"log/slog"
	"sync"

	"go.opentelemetry.io/contrib/bridges/otelslog"
)

var name string = "utils/otelutil"
var once sync.Once

var logger *slog.Logger

func Logger() *slog.Logger {
	once.Do(func() {
		logger = otelslog.NewLogger(name)
	})
	return logger
}
