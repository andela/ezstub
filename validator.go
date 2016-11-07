package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// Validator validates a request
type Validator interface {
	// Valid checks if the http request should go through.
	Valid(*http.Request) bool
}

// ValidatorFunc makes a func a valid Validator.
type ValidatorFunc func(*http.Request) bool

// Valid satisfies Validator interface.
func (v ValidatorFunc) Valid(r *http.Request) bool {
	return v(r)
}

// ParamValidator is a validator for query and form params.
type ParamValidator KV

// Valid satisfies Validator interface.
func (p ParamValidator) Valid(r *http.Request) bool {
	return compare(r.FormValue(p.Key), KV(p).ValueStr())
}

// HeaderValidator is a validator for request headers.
type HeaderValidator KV

// Valid satisfies Validator interface.
func (h HeaderValidator) Valid(r *http.Request) bool {
	return compare(r.Header.Get(h.Key), KV(h).ValueStr())
}

// JSONValidator is a validtor for JSON body.
type JSONValidator KV

// Valid satisfies validator interface.
func (j JSONValidator) Valid(r *http.Request) bool {
	// read body
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return false
	}

	// refill request body
	r.Body = ioutil.NopCloser(bytes.NewBuffer(b))

	var current interface{}
	if err := json.Unmarshal(b, &current); err != nil {
		log.Println(err)
		return false
	}

	// Retrieve value
	chain := strings.Split(j.Key, ".")
	var ok bool
	for _, c := range chain {
		if num, err := strconv.Atoi(c); err == nil {
			// list
			var array []interface{}
			if array, ok = current.([]interface{}); !ok {
				return false
			}
			if num >= len(array) {
				return false
			}
			current = array[num]
		} else {
			// map
			var m map[string]interface{}
			if m, ok = current.(map[string]interface{}); !ok {
				return false
			}
			if current, ok = m[c]; !ok {
				return false
			}
		}
	}
	return current == j.Value
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

const regexPrefix = "/r/"

// compare compares a and b. b will be regarded as a regex if prefixed with
// /r/.
func compare(a, b string) bool {
	if !strings.HasPrefix(b, regexPrefix) {
		return a == b
	}
	b = strings.TrimPrefix(b, regexPrefix)
	matches, _ := regexp.MatchString(b, a)
	return matches
}
