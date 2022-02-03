package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/wboard82/rekwest-bin/db_controller"
)

var templates = template.Must(template.ParseFiles("templates/inspect.html"))

func main() {
	db_controller.Connect()
	defer db_controller.Disconnect()

	http.HandleFunc("/r/", binHandler)
	http.HandleFunc("/", rootHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Welcome to Rekwest Bin</h1><form method='POST' action='/r/'><button type='submit'>Create a bin</button></form>")
}

func fixIPAddress(r *http.Request) {
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
}

func binHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		bin, _ := db_controller.NewBin()

		http.Redirect(w, r, "/r/"+bin.BinId+"?inspect", 302)
		return
	}

	binID := r.URL.Path[len("/r/"):]
	binAddress := fmt.Sprintf("http://%s/r/%s", r.Host, binID)

	if r.URL.RawQuery == "inspect" {
		bin, exists := db_controller.FindBin(binID)

		if !exists {
			http.NotFound(w, r)
			return
		}

		renderTemplate(w, "inspect", &bin)
	} else {
		dump, err := httputil.DumpRequest(r, true)

		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}

		rekwest := db_controller.Rekwest{string(dump)}

		fixIPAddress(r)

		if db_controller.AddRekwest(binID, rekwest) {
			fmt.Fprintf(w, "<h1>Request saved</h1><p>%s</p>", r.RemoteAddr)
			fmt.Fprintf(w, "<p><a href=%s>View requests</a>", binAddress+"?inspect")
		} else {
			http.NotFound(w, r)
		}
	}
}

func renderTemplate(writer http.ResponseWriter, tmpl string, bin *db_controller.Bin) {
	err := templates.ExecuteTemplate(writer, tmpl+".html", bin)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}
