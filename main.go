package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

func main() {
	if err := run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)
	var (
		port = flags.Int("port", 8080, "port to listen on")
	)
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}
	addr := fmt.Sprintf("0.0.0.0:%d", *port)
	srv, err := newServer()
	if err != nil {
		return err
	}

	fmt.Printf("Rekwest Bin listening on :%d\n", *port)
	return http.ListenAndServe(addr, srv)
}

type server struct {
	mux  *http.ServeMux
	tmpl *template.Template
	db   *Database
}

func newServer() (*server, error) {
	srv := &server{
		mux:  http.NewServeMux(),
		tmpl: template.Must(template.ParseGlob("templates/*.html")),
		db:   NewDatabase("rekwest-bin", "bins"),
	}

	srv.routes()
	return srv, nil
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.db.Connect()
	defer s.db.Disconnect()
	s.mux.ServeHTTP(w, r)
}

func (s *server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("public", "index.html"))
	}
}

/*
func main() {
	db_controller.Connect()
	defer db_controller.Disconnect()

	http.HandleFunc("/r/", binHandler)
	http.HandleFunc("/", rootHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "")
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
}

func renderTemplate(writer http.ResponseWriter, tmpl string, bin *db_controller.Bin) {
	err := templates.ExecuteTemplate(writer, tmpl+".html", bin)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}
*/
