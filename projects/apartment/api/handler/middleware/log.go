package middleware

import (
	"bytes"
	"net/http"
	"time"

	"github.com/arcaptcha-internship-2025/momoein-apartment/api/handler/router"
	"go.uber.org/zap"
)

type responseRecorder struct {
	http.ResponseWriter
	status int
	body   bytes.Buffer
}

func (rw *responseRecorder) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseRecorder) Write(b []byte) (int, error) {
	rw.body.Write(b) // Capture body
	return rw.ResponseWriter.Write(b)
}

func GetLogRequest(log *zap.Logger) router.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			log.Info("Incoming request",
				zap.String("method", r.Method),
				zap.String("url", r.URL.String()),
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("user_agent", r.UserAgent()),
			)

			rw := &responseRecorder{
				ResponseWriter: w,
				status:         http.StatusOK, // default if WriteHeader not called
			}
			next.ServeHTTP(rw, r)

			log.Info("Completed request",
				zap.String("method", r.Method),
				zap.String("url", r.URL.String()),
				zap.Int("status", rw.status),
				zap.Duration("duration", time.Since(start)),
			)
		})
	}
}
