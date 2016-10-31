package main

import (
	"fmt"
	"log"
	"net/http"
)

// Server is an Stub server.
type Server struct {
	host string
	port int
	stub *Stub
}

// NewServer creates a new Server.
func NewServer(config Config) (*Server, error) {
	stub, err := initStub(config)
	if err != nil {
		return nil, err
	}
	return &Server{
		port: config.Port,
		host: config.Host,
		stub: stub,
	}, nil
}

// Addr returns the host:port config for the server.
func (s *Server) Addr() string {
	return fmt.Sprintf("%v:%v", s.host, s.port)
}

// Start starts the server.
func (s *Server) Start() error {
	fmt.Println(s.stub.title)
	fmt.Println("ezstub listening on", s.Addr())
	fmt.Println()
	return http.ListenAndServe(s.Addr(), s)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w = &responseWriter{rw: w}

	// Log
	defer func() {
		rw := w.(*responseWriter)
		log.Println(r.Method, r.URL.Path, rw.status, rw.written)
	}()

	route, ok := s.stub.routes[r.URL.Path]
	if !ok {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	endpoint, ok := route.endpoints[Method(r.Method)]
	if !ok {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if !endpoint.validators.Valid(r) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	// Response
	endpoint.WriteResponse(w)
}
