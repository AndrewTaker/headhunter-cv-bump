package handler

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

type loggerKey int

const lk loggerKey = 0

func GetLogger(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(lk).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}
func LogRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		logger := slog.Default().With(
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("query", r.URL.RawQuery),
		)

		ctx := context.WithValue(r.Context(), lk, logger)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

		logger.Info("request completed",
			slog.Duration("latency", time.Since(start)),
		)
	})
}
