package handler

import (
	"net/http"
	"slices"
)

type Router struct {
	globalChain chain
	routeChain  chain
	isSubRoute  bool
	*http.ServeMux
}

func NewRouter() *Router {
	return &Router{ServeMux: http.NewServeMux()}
}

func (r *Router) Use(mw ...func(http.Handler) http.Handler) {
	if r.isSubRoute {
		r.routeChain = append(r.routeChain, mw...)
	}
	r.globalChain = append(r.globalChain, mw...)
}

func (r *Router) Group(fn func(*Router)) {
	subRouter := &Router{
		routeChain: slices.Clone(r.routeChain),
		isSubRoute: true,
		ServeMux:   r.ServeMux,
	}
	fn(subRouter)
}

func (r *Router) HandleFunc(pattern string, h http.HandlerFunc) {
	r.Handle(pattern, h)
}

func (r *Router) Handle(pattern string, h http.Handler) {
	h = r.routeChain.Then(h)
	r.ServeMux.Handle(pattern, h)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, rq *http.Request) {
	var h http.Handler = r.ServeMux
	h = r.globalChain.Then(h)
	h.ServeHTTP(w, rq)
}
