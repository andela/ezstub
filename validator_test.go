package main

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
)

func TestValidator(t *testing.T) {
	tests := []KV{
		{"key", "val"},
		{"another_key", "another val"},
	}

	for i, test := range tests {
		url := fmt.Sprintf("/?%s=%s", test.Key, test.Value)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			t.Fatal(err)
		}
		// Param Validator
		if !ParamValidator(test).Valid(req) {
			t.Errorf("Test %d: params should be valid", i)
		}

		// Header Validator
		req.Header.Add(test.Key, test.ValueStr())
		if !HeaderValidator(test).Valid(req) {
			t.Errorf("Test %d: headers should be valid", i)
		}

		// Validators
		var v Validators
		v.Add(ParamValidator(test))
		v.Add(HeaderValidator(test))
		if !v.Valid(req) {
			t.Errorf("Test %d: params and headers should be valid", i)
		}
	}
}

func TestJSONValidator(t *testing.T) {
	tests := []KV{
		{"id", float64(1)},
		{"user.name", "John Doe"},
		{"user.titles.1", "Dr"},
		{"user.schools.0.name", "Queens College"},
	}
	jsonStr := `{
		"id" : 1,
		"user" : {
			"name" : "John Doe",
			"titles" : ["Mr", "Dr"],
			"schools" : [
				{ "name" : "Queens College" }
			]
		}
	}`
	req, err := http.NewRequest("POST", "/", bytes.NewBufferString(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	for i, test := range tests {
		if !JSONValidator(test).Valid(req) {
			t.Errorf("Test %d: validation failed for %v. Expected %#v", i, test.Key, test.Value)
		}
	}
}
