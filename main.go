package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/http/httputil"
)

var templates = template.Must(template.ParseFiles("templates/inspect.html"))

var binStore = NewBinStore()

func main() {
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

	for source, ip := range ipSources {
		fmt.Println(source, ": ", ip)
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
	// The POST route creates a new bin and redirects to the inspect page
	if r.Method == "POST" {
		binName := binStore.NewBin()

		http.Redirect(w, r, "/r/"+binName+"?inspect", 302)
		return
	}

	// This grabs the part after /r/ in the path
	binID := r.URL.Path[len("/r/"):]
	// Put the full link together here to be displayed on a landing page
	binAddress := fmt.Sprintf("http://%s/r/%s", r.Host, binID)

	// If there is a query "inspect", show all the requests
	if r.URL.RawQuery == "inspect" {
		rekwests, exists := loadRequest(binID)

		if !exists {
			http.NotFound(w, r)
			return
		}

		requestInfo := make([]RequestInfo, len(rekwests))

		for i, req := range rekwests {
			requestInfo[i] = RequestInfo{req}
		}

		bin := BinInfo{
			BinAddress: binAddress,
			Requests:   requestInfo,
		}

		renderTemplate(w, "inspect", &bin)

	} else {
		dump, err := httputil.DumpRequest(r, true)

		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}

		fixIPAddress(r)

		if saveRequest(binID, dump) {
			fmt.Fprintf(w, "<h1>Request saved</h1><p>%s</p>", r.RemoteAddr)
			fmt.Fprintf(w, "<p><a href=%s>View requests</a>", binAddress+"?inspect")
		} else {
			http.NotFound(w, r)
		}
	}
}

func renderTemplate(writer http.ResponseWriter, tmpl string, bin *BinInfo) {
	err := templates.ExecuteTemplate(writer, tmpl+".html", bin)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}
