package middleware

import (
	"go-banking/pkg/response"
	"log"
	"net/http"
	"time"
)

//middlware is used to log the incoming requests and the time taken to process them.
//It wraps the http.ResponseWriter to capture the status code of the response and logs the
//HTTP method, request path, status code, and duration of the request processing. This helps
//in monitoring and debugging the API by providing insights into the performance and any
//potential issues with the requests being handled by the server.

//In simple words, this middleware acts acts as a logger for all incoming HTTP requests,
//recording the method, path, status code, and how long it took to process each request.

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		wrappedWriter := newResponseWriter(w)

		next.ServeHTTP(wrappedWriter, r)

		duration := time.Since(startTime)
		log.Printf("%s %s %d %s", r.Method, r.URL.Path, wrappedWriter.statusCode, duration)
	})
}

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %v", err)
				response.WriteError(w, http.StatusInternalServerError, "internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}
