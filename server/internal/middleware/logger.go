package middleware

import (
	"fmt"
	"log"
	"net/http"
)

type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// Constructs LoggingResponseWriter handler
func NewLoggingResponseWriter(w http.ResponseWriter) *LoggingResponseWriter {
	// Default status code if no WriteHeader is used
	return &LoggingResponseWriter{w, http.StatusOK}
}

func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	// Update the status code in the lrw
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	// Replace the ResponseWriter with the LoggingResponseWriter struct

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := NewLoggingResponseWriter(w)

		// Pass the logging response writer to the request and serve
		next.ServeHTTP(lrw, r)

		// Setup the output
		logString := fmt.Sprintf("[http] [%d] %s User-Agent: %s", lrw.statusCode, r.RequestURI, r.Header.Get("User-Agent"))
		log.Println(logString)
	})
}
