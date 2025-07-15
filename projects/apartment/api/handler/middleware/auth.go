package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/arcaptcha-internship-2025/momoein-apartment/api/handler/router"
	appctx "github.com/arcaptcha-internship-2025/momoein-apartment/pkg/context"
	appjwt "github.com/arcaptcha-internship-2025/momoein-apartment/pkg/jwt"
	"go.uber.org/zap"
)

const (
	UserEmailKey appctx.CtxKey = "UserEmail"
	UserIDKey    appctx.CtxKey = "UserID"
)

func GetAuth(secret []byte) router.Middleware {
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

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UserEmailKey, claims.UserMail)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
