package interfaces

import "net/http"

//go:generate mockgen -destination mocks/http.go -package mocks github.com/cryplio/rest-api/src/interfaces HTTPGetter

// HTTPGetter is an interface to Get data over http
type HTTPGetter interface {
	Get(url string) (resp *http.Response, err error)
}
