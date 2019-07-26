/*
Author: Connor Sanders
MIT License
RESTful API Boilerplate
7/19/2019
*/

package server

import (
	"log"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// NewLoggingResponseWriter is a function used to write API log data to the console
func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

// Writes response headers for API
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// Logger is a function that is used to manage the API logging process
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
		// Todo: Add functionality Write to a log file here in the /log directory here
	})
}
