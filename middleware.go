package main

import (
	"net/http"
)

type Middleware []http.Handler

type MiddlewareResponsWriter struct {
	http.ResponseWriter
	written bool
}

func (middleware *Middleware) Add(handler http.Handler) {
	*middleware = append(*middleware, handler)
}

func NewMiddlewareResponseWriter(responseWriter http.ResponseWriter) *MiddlewareResponsWriter {
	return &MiddlewareResponsWriter{
		ResponseWriter: responseWriter,
	}
}

func (middlewareResponseWriter *MiddlewareResponsWriter) Write(bytes []byte) (int, error) {
	middlewareResponseWriter.written = true

	return middlewareResponseWriter.ResponseWriter.Write(bytes)
}

func (middlewareResponseWriter *MiddlewareResponsWriter) WriteHeader(code int) {
	middlewareResponseWriter.written = true

	middlewareResponseWriter.ResponseWriter.WriteHeader(code)
}

func (middleware Middleware) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	// Wrap the supplied responseWriter object
	middlewareResponseWriter := NewMiddlewareResponseWriter(responseWriter)

	// loop through all the registered handlers
	for _, handler := range middleware {
		// call the handler with the wrapped MiddlewareResponsWriter
		handler.ServeHTTP(middlewareResponseWriter, request)

		// stop processing, if any response was written
		if middlewareResponseWriter.written {
			return
		}
	}

	// if no handlers wrote a response, return a 404
	http.NotFound(responseWriter, request)
}
