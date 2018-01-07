package api

import (
	"net/http"

	"github.com/gorilla/handlers"
)

// AllowedOrigins is a list containing all origins allowed to hit the API
var AllowedOrigins = handlers.AllowedOrigins([]string{
	"https://www.crypl.io/",       // prod
	"https://swan.crypl.io",       // staging
	"http://orchid.crypl.io:4200", // local
})

// AllowedMethods is a list containing all HTTP verb accepted by the API
var AllowedMethods = handlers.AllowedMethods([]string{
	http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete,
})

// AllowedHeaders is a list custom headers accepted by the API
var AllowedHeaders = handlers.AllowedHeaders([]string{
	"Content-Type", "Authorization",
})
