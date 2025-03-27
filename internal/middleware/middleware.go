package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

func Apply(h http.Handler, mw ...Middleware) http.Handler {
	for _, middleware := range mw {
		h = middleware(h)
	}
	return h
}

func MiddlewareStack(ms ...Middleware) Middleware {
	return Middleware(func(next http.Handler) http.Handler {
		for i := len(ms) - 1; i >= 0; i-- {
			m := ms[i]
			next = m(next)
		}
		return next
	})
}