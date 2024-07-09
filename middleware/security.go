package middleware

import (
	"net/http"
	"os"
	"strings"
)

const CSP = "default-src 'self'; script-src 'self' cdn.jsdelivr.net 'unsafe-inline'; style-src 'self' 'unsafe-inline' fonts.googleapis.com; font-src 'self' fonts.googleapis.com fonts.gstatic.com;"

var ALLOWED_METHODS = []string{"GET", "POST", "PATCH", "PUT", "DELETE"}

func contains(slices []string, search string) bool {
	for _, value := range slices {
		if strings.ToUpper(value) == strings.ToUpper(search) {
			return true
		}
	}
	return false
}

func Security(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		goEnv := os.Getenv("GO_ENV")

		if goEnv == "staging" || goEnv == "production" {
			method := r.Method
			// Check if the request method is allowed
			if !contains(ALLOWED_METHODS, method) {
				w.WriteHeader(http.StatusMethodNotAllowed)
				w.Write([]byte("405 Method Not Allowed"))
				return
			}

			w.Header().Set("X-XSS-Protection", "1; mode=block")
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("Content-Security-Policy", CSP)
			w.Header().Set("Referrer-Policy", "no-referrer")
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		}

		next.ServeHTTP(w, r)
	})
}
