package handler

import (
	"net/http"
	"slices"
)

type chain []func(http.Handler) http.Handler

func (c chain) Then(h http.Handler) http.Handler {
	for _, mw := range slices.Backward(c) {
		h = mw(h)
	}
	return h
}
