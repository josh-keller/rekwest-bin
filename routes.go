package main

func (s *server) routes() {
	s.mux.Handle("/socket.io/", s.sockets)
	s.mux.HandleFunc("/", s.handleIndex())
	s.mux.HandleFunc("/new/", s.handleBinNew())
	s.mux.HandleFunc("/r/", s.fixIPAddress(s.handleRequest()))
	s.mux.HandleFunc("/inspect/", s.handleBinInspect())
}
