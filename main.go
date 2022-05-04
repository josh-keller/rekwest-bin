package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	socketio "github.com/googollee/go-socket.io"
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

	srv.db.Connect()
	defer srv.db.Disconnect()

	go srv.sockets.Serve()
	defer srv.sockets.Close()

	fmt.Printf("Rekwest Bin listening on :%d\n", *port)
	return http.ListenAndServe(addr, srv)
}

type server struct {
	mux     *http.ServeMux
	tmpl    map[string]*template.Template
	db      *Database
	sockets *socketio.Server
}

func newServer() (*server, error) {
	srv := &server{
		mux: http.NewServeMux(),
		tmpl: map[string]*template.Template{
			"inspect": template.Must(template.ParseFiles("templates/inspect.html", "templates/rekwest.html")),
			"rekwest": template.Must(template.ParseFiles("templates/rekwest.html")),
		},
		sockets: socketio.NewServer(nil),
		db:      NewDatabase("rekwest-bin", "bins"),
	}

	srv.socketRoutes()
	srv.routes()
	return srv, nil
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("public", "index.html"))
	}
}

func (s *server) renderTemplate(writer http.ResponseWriter, tmpl string, bin *Bin) {
	err := s.tmpl[tmpl].ExecuteTemplate(writer, tmpl+".html", bin)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}
