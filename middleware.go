package main

import "net/http"

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
