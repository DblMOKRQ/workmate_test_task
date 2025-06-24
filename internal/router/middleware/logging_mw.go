package middleware

import (
	"net/http"

	"go.uber.org/zap"
)

func LoggingMiddleware(next http.Handler, log *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info("request", zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.String("proto", r.Proto), zap.String("remote_addr", r.RemoteAddr))
		next.ServeHTTP(w, r)
	})
}
