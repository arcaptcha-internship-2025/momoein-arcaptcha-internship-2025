package middleware

import (
	"net/http"

	"github.com/arcaptcha-internship-2025/momoein-apartment/api/handler/router"
	"github.com/arcaptcha-internship-2025/momoein-apartment/app"
	appctx "github.com/arcaptcha-internship-2025/momoein-apartment/pkg/context"
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
