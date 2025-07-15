package router

import (
	"net/http"
	"slices"
)

type Middleware = func(http.Handler) http.Handler

type Chain []Middleware

func (c Chain) Then(h http.Handler) http.Handler {
	for _, mw := range slices.Backward(c) {
		h = mw(h)
	}
	return h
}
