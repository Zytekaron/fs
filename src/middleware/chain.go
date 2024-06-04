package middleware

import "net/http"

type Middleware func(h http.Handler) http.Handler

// Chain chains middleware in order, returning a new Middleware.
func Chain(middleware ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middleware) - 1; i >= 0; i-- {
			mw := middleware[i]
			next = mw(next)
		}
		return next
	}
}
