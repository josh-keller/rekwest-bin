package main

import (
	"fmt"
	"net/http"
)

func (s *server) handleRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("handleRequest")
		binID := r.URL.Path[len("/r/"):]
		fmt.Println("Request made: ", binID)

		if err := s.db.AddRekwest(binID, r); err == nil {
			fmt.Fprintf(w, "Request saved. (ip: %s)", r.RemoteAddr)
		} else {
			http.NotFound(w, r)
		}
	}
}
