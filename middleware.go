package main

import (
	"fmt"
	"net/http"
	"time"
)

func (s *server) fixIPAddress(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ipAddress string
		var ipSources = []string{
			r.Header.Get("True-Client-IP"),
			r.Header.Get("True-Real-IP"),
			r.Header.Get("X-Forwarded-For"),
			r.Header.Get("X-Originating-IP"),
		}

		for _, ip := range ipSources {
			if ip != "" {
				ipAddress = ip
				break
			}
		}

		if ipAddress != "" {
			r.RemoteAddr = ipAddress
		}

		h(w, r)
	}
}

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func (s *server) withLogging(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}

		lrw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		h.ServeHTTP(lrw, r) // serve the original request

		uri := r.RequestURI
		method := r.Method

		duration := time.Since(start)

		// log request details
		fmt.Printf("%s %s (%d) - %d bytes, %v\n", method, uri, lrw.responseData.status, lrw.responseData.size, duration)
	}
}
