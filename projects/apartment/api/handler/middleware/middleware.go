package middleware

import (
	"bytes"
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/arcaptcha-internship-2025/momoein-apartment/api/handler/router"
	"github.com/arcaptcha-internship-2025/momoein-apartment/app"
	appctx "github.com/arcaptcha-internship-2025/momoein-apartment/pkg/context"
	appjwt "github.com/arcaptcha-internship-2025/momoein-apartment/pkg/jwt"
	"go.uber.org/zap"
)

func SetRequestContext(app app.App) router.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = appctx.New(ctx, appctx.WithLogger(app.Logger()))
			req := r.WithContext(ctx)
			next.ServeHTTP(w, req)
		})
	}
}

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

func LogRequest() router.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			log := appctx.Logger(r.Context())
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

func NewAuth(secret []byte) router.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := appctx.Logger(r.Context())

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := appjwt.ParseToken(token, secret)
			if err != nil {
				switch err {
				case appjwt.ErrInvalidToken, appjwt.ErrNilToken:
					log.Warn("parse jwt token", zap.Error(err))
				default:
					log.Error("parse jwt token", zap.Error(err))
				}
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), appjwt.UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, appjwt.UserEmailKey, claims.UserEMail)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
