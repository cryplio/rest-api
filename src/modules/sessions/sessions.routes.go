package sessions

import (
	"github.com/Nivl/go-rest-tools/dependencies"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/gorilla/mux"
)

// Contains the index of all Endpoints
const (
	EndpointAdd = iota
	EndpointDelete
)

// Endpoints is a list of endpoints for this components
var Endpoints = router.Endpoints{
	EndpointAdd:    addEndpoint,
	EndpointDelete: deleteEndpoint,
}

// SetRoutes is used to set all the sessions routes
func SetRoutes(r *mux.Router, deps dependencies.Dependencies) {
	Endpoints.Activate(r, deps)
}
