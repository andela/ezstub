package main

import "net/http"

// ParamValidator is a validator for query and form params.
type ParamValidator KV

// Valid satisfies Validator interface.
func (p ParamValidator) Valid(r *http.Request) bool {
	return r.FormValue(p.Key) == p.Value
}

// HeaderValidator is a validator for request headers.
type HeaderValidator KV

// Valid satisfies Validator interface.
func (h HeaderValidator) Valid(r *http.Request) bool {
	return r.Header.Get(h.Key) == h.Value
}

// Validator validates a request
type Validator interface {
	// Valid checks if the http request should go through.
	Valid(*http.Request) bool
}

// Validators is a chain of validator.
type Validators []Validator

// Valid satisfies Validator interface.
func (v Validators) Valid(r *http.Request) bool {
	for i := range v {
		if !v[i].Valid(r) {
			return false
		}
	}
	return true
}

// Add adds a new validator to the chain.
func (v *Validators) Add(validator Validator) {
	*v = append(*v, validator)
}
