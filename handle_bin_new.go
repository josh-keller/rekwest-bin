package main

import (
	"fmt"
	"net/http"
)

func (s *server) handleBinNew() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("handleBinNew")
		switch r.Method {
		case "POST":
			bin, _ := s.db.NewBin()
			http.Redirect(w, r, "/inspect/"+bin.BinId, 302)
			fmt.Println("Redirected")
			return
		default:
			http.NotFound(w, r)
			return
		}
	}
}
