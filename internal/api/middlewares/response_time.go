package middlewares

import (
	"fmt"
	"net/http"
	"time"
)

func ResponseTimeMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		wrappedWriter := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(wrappedWriter, r)

		duration := time.Since(start)
		fmt.Printf("Request to %s took %v\n", r.URL.Path, duration)
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
