package main

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
)

// Config holds the configuration.
type Config struct {
	Port      int              `yaml:"port"`
	Host      string           `yaml:"host"`
	Title     string           `yaml:"title"`
	Endpoints []EndpointConfig `yaml:"endpoints"`
}

// EndpointConfig is config for endpoints.
type EndpointConfig struct {
	URL         string           `yaml:"url"`
	Description string           `yaml:"description"`
	Method      Method           `yaml:"method"`
	Validation  ValidationConfig `yaml:"validation"`
	Response    ResponseConfig   `yaml:"response"`
}

// Endpoint converts EndpointConfig into an Endpoint.
func (e EndpointConfig) Endpoint() (Endpoint, error) {
	response, err := e.Response.Response()
	if err != nil {
		return Endpoint{}, err
	}
	if e.Method == "" {
		e.Method = "GET"
	}
	// validators
	var validators Validators
	for _, kv := range e.Validation.Headers {
		validators.Add(HeaderValidator(kv))
	}
	for _, kv := range e.Validation.Params {
		validators.Add(ParamValidator(kv))
	}

	return Endpoint{
		description: e.Description,
		method:      e.Method,
		response:    response,
		validators:  validators,
	}, nil
}

// ResponseConfig is config for responses.
type ResponseConfig struct {
	Data    string `yaml:"data"`
	File    string `yaml:"file"`
	Status  int    `yaml:"status"`
	Headers []KV   `yaml:"headers"`
}

// Response converts config into a Response.
func (r ResponseConfig) Response() (response Response, err error) {
	if r.Status == 0 {
		r.Status = http.StatusOK
	}
	response.status = r.Status
	response.headers = r.Headers
	switch {
	case r.File != "":
		response.body, err = ioutil.ReadFile(r.File)
	case r.Data != "":
		response.body, err = base64.StdEncoding.DecodeString(r.Data)
	}
	return
}

// ValidationConfig is config for key value pairs.
type ValidationConfig struct {
	Headers []KV `yaml:"headers"`
	Params  []KV `yaml:"params"`
}

// KV is key value.
type KV struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

func initStub(config Config) (*Stub, error) {
	stub := new(Stub)
	stub.title = config.Title
	stub.routes = make(map[string]Route)
	for _, e := range config.Endpoints {
		var route Route
		if _, ok := stub.routes[e.URL]; ok {
			route = stub.routes[e.URL]
		}
		endpoint, err := e.Endpoint()
		if err != nil {
			return nil, err
		}
		route.AddEndpoint(endpoint)
		stub.routes[e.URL] = route
	}
	return stub, nil
}
