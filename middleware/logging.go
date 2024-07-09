package middleware

import (
	"log"
	"net/http"
	"os"
	"time"
)

// wrapped writer to expost status code for logging
type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	// if header isn't written, write it
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Serve the request
		next.ServeHTTP(wrapped, r)

		// Log the request
		log.SetOutput(os.Stdout)
		log.Println(wrapped.statusCode, r.Method, r.RequestURI, time.Since(start))
	})
}
