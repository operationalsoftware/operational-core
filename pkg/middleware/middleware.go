package middleware

import (
	"net/http"
)

// Middleware type definition
type Middleware func(http.Handler) http.Handler

// ChainMiddleware function to chain multiple middleware
func CreateStack(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}

		return next
	}
}
