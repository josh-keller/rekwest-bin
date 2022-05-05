package main

import (
	"net/http"
)

func (s *server) handleBinInspect() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			binID := r.URL.Path[len("/inspect/"):]

			bin, exists := s.db.FindBin(binID)

			if !exists {
				http.NotFound(w, r)
				return
			}

			bin.Host = r.Host

			s.renderTemplate(w, "inspect", &bin)
			return

		default:
			http.NotFound(w, r)
			return
		}
	}
}
