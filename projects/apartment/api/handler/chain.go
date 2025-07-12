package handler

import (
	"net/http"
	"slices"
)

type chain []func(http.Handler) http.Handler

func (c *chain) Append(mws ...func(http.Handler) http.Handler) {
	*c = append(*c, mws...)
}

func (c *chain) Then(h http.Handler) http.Handler {
	for _, mw := range slices.Backward(*c) {
		h = mw(h)
	}
	return h
}

func Chain(middlewares ...func(http.Handler) http.Handler) *chain {
	c := make(chain, 0)
	c = append(c, middlewares...)
	return &c
}
