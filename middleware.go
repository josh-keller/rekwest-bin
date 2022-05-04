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

func (s *server) withLogging(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		uri := r.RequestURI
		method := r.Method
		h.ServeHTTP(w, r) // serve the original request

		duration := time.Since(start)

		// log request details
		fmt.Println(uri, method, duration)
	}
}
