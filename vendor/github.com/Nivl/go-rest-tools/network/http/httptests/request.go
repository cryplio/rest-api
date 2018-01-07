package httptests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

// DefaultRouter represents the router to be used if none are specified
var DefaultRouter *mux.Router

// NewRequest simulates a new http request executed against the api
func NewRequest(t *testing.T, info *RequestInfo) *httptest.ResponseRecorder {
	info.ParseParams()

	// Parse the body as a JSON object
	mime, body, err := info.Body()
	if err != nil {
		t.Fatalf("could not create request: %s", err)
	}

	req, err := http.NewRequest(info.Endpoint.Verb, info.URL(), body)
	if err != nil {
		t.Fatalf("could not execute request %s", err)
	}

	// Attach the query string
	qs := req.URL.Query()
	info.PopulateQuery(qs)
	req.URL.RawQuery = qs.Encode()

	if info.Auth != nil {
		req.Header.Add("Authorization", info.Auth.ToBasicAuth())
	}

	req.Header.Add("Content-Type", mime)

	// If no router is provided we assume that we want to execute a regular endpoint
	if info.Router == nil {
		if DefaultRouter == nil {
			t.Fatalf("no router provided and DefaultRouter not set")
		}
		info.Router = DefaultRouter
	}

	rec := httptest.NewRecorder()
	info.Router.ServeHTTP(rec, req)
	return rec
}
