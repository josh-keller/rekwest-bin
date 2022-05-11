package main

func (s *server) routes() {
	s.mux.HandleFunc("/", s.withLogging(s.handleRoot()))
	s.mux.HandleFunc("/new/", s.withLogging(s.handleBinNew()))
	s.mux.HandleFunc("/r/", s.withLogging(s.fixIPAddress(s.handleRequest())))
	s.mux.HandleFunc("/inspect/", s.withLogging(s.handleBinInspect()))
}
