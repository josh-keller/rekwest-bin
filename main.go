package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
)

func binHandler(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Path[len("/r/"):]

	if r.URL.RawQuery == "inspect" {
		rekwest, err := loadRequest(hash)

		if err != nil {
			http.NotFound(w, r)
			return
		}

		fmt.Fprintf(w, "<h1>Here is the last rekwest:</h1><p>%s</p>", rekwest)
	} else {
		var buf bytes.Buffer
		r.Write(&buf)
		requestString := buf.Bytes()

		saveRequest(hash, requestString)

		fmt.Fprintf(w, "<h1>Request saved</h1><p>%s</p>", r.RemoteAddr)
	}
}

func main() {
	http.HandleFunc("/r/", binHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func saveRequest(hash string, request []byte) error {
	filename := hash + ".txt"
	return os.WriteFile("./rekwests/"+filename, request, 0600)
}

func loadRequest(hash string) ([]byte, error) {
	filename := hash + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}
