package main

import "testing"
import "net/http"

import "fmt"

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
		req.Header.Add(test.Key, test.Value)
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
