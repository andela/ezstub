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
func newServer(config Config) (*Server, error) {
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
	fmt.Println("Listening on", s.Addr())
	fmt.Println()
	fmt.Println(s.stub.title)
	fmt.Println("Routes:")
	s.stub.dump()
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
		errorResponse(w, http.StatusNotFound)
		return
	}
	endpoint, ok := route.endpoints[Method(r.Method)]
	if !ok {
		errorResponse(w, http.StatusNotFound)
		return
	}
	if !endpoint.validators.Valid(r) {
		errorResponse(w, http.StatusForbidden)
		return
	}

	// Response
	endpoint.WriteResponse(w)
}

func errorResponse(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}
