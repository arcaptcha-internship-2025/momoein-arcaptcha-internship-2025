package handler

import (
	"fmt"
	"net/http"

	"github.com/arcaptcha-internship-2025/momoein-apartment/api/handler/middleware"
	"github.com/arcaptcha-internship-2025/momoein-apartment/api/handler/router"
	"github.com/arcaptcha-internship-2025/momoein-apartment/app"
)

func Run(app app.App) error {
	r := router.NewRouter()
	r.Use(middleware.GetLogRequest(app.Logger()))
	r.Get("/", getRootHandler())

	api := r.Group("/api/v1", nil)
	RegisterAPI(api, app)

	addr := fmt.Sprintf(":%d", app.Config().HTTP.Port)
	app.Logger().Info("listen on " + addr)
	return http.ListenAndServe(addr, r)
}

func RegisterAPI(r *router.Router, app app.App) {
	r.Group("/auth", func(r *router.Router) {
		r.Post("/sing-up", getSignUpHandler())
	})
}

func getRootHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
		w.Write([]byte("hello from arcaptcha apartment api\n"))
	})
}
