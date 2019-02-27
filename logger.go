package main

import (
	"log"
	"net/http"
	"time"
)


type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}


func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}


func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}


func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := NewLoggingResponseWriter(w)
		inner.ServeHTTP(lrw, r)
		statusCode := lrw.statusCode
		log.Printf(
			"%s\t%s\t%s\t%d\t%s",
			r.Method,
			r.RequestURI,
			name,
			statusCode,
			time.Since(start),
		)
	})
}
