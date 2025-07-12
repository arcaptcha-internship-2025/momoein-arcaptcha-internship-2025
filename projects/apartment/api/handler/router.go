package handler

import (
	"fmt"
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

func (r *Router) Group(fn func(*Router)) *Router {
	subRouter := &Router{
		routeChain: slices.Clone(r.routeChain),
		isSubRoute: true,
		ServeMux:   r.ServeMux,
	}
	fn(subRouter)
	return subRouter
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

func (r *Router) Method(method, pattern string, h http.Handler) {
	r.Handle(fmt.Sprintf("%s %s", method, pattern), h)
}

func (r *Router) Connect(pattern string, h http.Handler) {
	r.Method(http.MethodConnect, pattern, h)
}

func (r *Router) Delete(pattern string, h http.Handler) {
	r.Method(http.MethodDelete, pattern, h)
}

func (r *Router) Get(pattern string, h http.Handler) {
	r.Method(http.MethodGet, pattern, h)
}

func (r *Router) Head(pattern string, h http.Handler) {
	r.Method(http.MethodGet, pattern, h)
}

func (r *Router) Options(pattern string, h http.Handler) {
	r.Method(http.MethodOptions, pattern, h)
}

func (r *Router) Patch(pattern string, h http.Handler) {
	r.Method(http.MethodPatch, pattern, h)
}

func (r *Router) Post(pattern string, h http.Handler) {
	r.Method(http.MethodPost, pattern, h)
}

func (r *Router) Put(pattern string, h http.Handler) {
	r.Method(http.MethodPut, pattern, h)
}

func (r *Router) Trace(pattern string, h http.Handler) {
	r.Method(http.MethodTrace, pattern, h)
}
