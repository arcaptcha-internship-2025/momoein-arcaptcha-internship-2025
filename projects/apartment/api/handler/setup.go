package handler

import (
	"fmt"
	"net/http"

	"github.com/arcaptcha-internship-2025/momoein-apartment/app"
)

func Run(app app.App) error {
	r := NewRouter()
	r.Use(firstMiddleware)

	addr := fmt.Sprintf(":%d", app.Config().HTTP.Port)
	app.Logger().Info("listen on: " + addr)
	return http.ListenAndServe(addr, r)
}

func RegisterAPI(r *Router, app app.App) {
	r.Group(func(r *Router) {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("group first middleware\n"))
				next.ServeHTTP(w, r)
			})
		})
		r.HandleFunc("GET /group", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("hello group\n"))
		}))
	})

	r.Group(func(r *Router) {
		r.Use(secondMiddleware)

		r.HandleFunc("GET /admin", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("hello group\n"))
		}))
	})

	mwChain := chain{firstMiddleware, secondMiddleware}
	r.Handle(fmt.Sprintf("%s %s", http.MethodGet, "/"), mwChain.Then(myHandler()))
}

func myHandler() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("hello world"))
	}
	return http.HandlerFunc(fn)
}

func firstMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Im first middleware\n"))
		next.ServeHTTP(w, req)
	}
	return http.HandlerFunc(fn)
}

func secondMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Im second middleware\n"))
		next.ServeHTTP(w, req)
	}
	return http.HandlerFunc(fn)
}
