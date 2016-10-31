package main

import (
	"log"
	"net/http"
)

// Stub is an endpoint stub.
type Stub struct {
	title      string
	routes     map[string]Route
	validators Validators
}

// Route routes requests to different endpoints.
type Route struct {
	url        string
	endpoints  map[Method]Endpoint
	validators Validators
}

// AddEndpoint adds an endpoint to the route.
func (r *Route) AddEndpoint(e Endpoint) {
	if r.endpoints == nil {
		r.endpoints = make(map[Method]Endpoint)
	}
	r.endpoints[e.method] = e
}

// Method is request method.
type Method string

// Endpoint is an API endpoint.
type Endpoint struct {
	method      Method
	description string
	response    Response
	validators  Validators
}

// WriteResponse writes the endpoint response.
func (e Endpoint) WriteResponse(w http.ResponseWriter) {
	// Headers
	for _, header := range e.response.headers {
		w.Header().Add(header.Key, header.Value)
	}

	// Response
	w.WriteHeader(e.response.status)
	w.Write(e.response.body)

	// Log
	log.Println()
	log.Println(e.description)
}

// Response is a request response
type Response struct {
	body    []byte
	headers []KV
	status  int
}
