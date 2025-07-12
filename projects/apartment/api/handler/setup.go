package handler

import (
	"fmt"
	"net/http"

	"github.com/arcaptcha-internship-2025/momoein-apartment/app"
)

func Run(app app.App) error {
	mux := http.NewServeMux()

	mwChain := Chain(firstMiddleware, secondMiddleware)

	mux.Handle(fmt.Sprintf("%s %s", http.MethodGet, "/"), mwChain.Then(myHandler()))

	addr := fmt.Sprintf(":%d", app.Config().HTTP.Port)
	app.Logger().Info("listen on: " + addr)
	return http.ListenAndServe(addr, mux)
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