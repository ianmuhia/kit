// Package testutil provides testing utilities and helpers
package testutil

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// HTTPTestCase represents a test case for HTTP handlers
type HTTPTestCase struct {
	Name           string
	Method         string
	URL            string
	Body           interface{}
	Headers        map[string]string
	ExpectedStatus int
	ExpectedBody   interface{}
	Setup          func(*testing.T)
	Teardown       func(*testing.T)
}

// DoHTTPRequest performs an HTTP request for testing
func DoHTTPRequest(t *testing.T, handler http.Handler, method, url string, body interface{}) *httptest.ResponseRecorder {
	t.Helper()
	
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("failed to marshal body: %v", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}
	
	req := httptest.NewRequest(method, url, bodyReader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	
	return rec
}

// AssertStatus asserts that the response has the expected status code
func AssertStatus(t *testing.T, rec *httptest.ResponseRecorder, expected int) {
	t.Helper()
	if rec.Code != expected {
		t.Errorf("expected status %d, got %d", expected, rec.Code)
	}
}

// AssertJSON asserts that the response body matches the expected JSON
func AssertJSON(t *testing.T, rec *httptest.ResponseRecorder, expected interface{}) {
	t.Helper()
	
	var actual interface{}
	if err := json.NewDecoder(rec.Body).Decode(&actual); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	
	expectedJSON, _ := json.Marshal(expected)
	actualJSON, _ := json.Marshal(actual)
	
	if string(expectedJSON) != string(actualJSON) {
		t.Errorf("expected body:\n%s\ngot:\n%s", expectedJSON, actualJSON)
	}
}

// AssertContains asserts that the response body contains the expected substring
func AssertContains(t *testing.T, rec *httptest.ResponseRecorder, expected string) {
	t.Helper()
	
	body := rec.Body.String()
	if !bytes.Contains([]byte(body), []byte(expected)) {
		t.Errorf("expected body to contain %q, got:\n%s", expected, body)
	}
}
