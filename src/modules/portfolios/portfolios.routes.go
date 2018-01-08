package portfolios

import (
	"github.com/Nivl/go-rest-tools/dependencies"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/gorilla/mux"
)

// Contains the index of all Endpoints
const (
	EndpointAdd = iota
	// EndpointUpdate
	EndpointDelete
	// EndpointList
)

// Endpoints is a list of endpoints for this components
var Endpoints = router.Endpoints{
	EndpointAdd: addEndpoint,
	// EndpointUpdate: updateEndpoint,
	EndpointDelete: deleteEndpoint,
	// EndpointList:   listEndpoint,
}

// SetRoutes is used to set all the users routes
func SetRoutes(r *mux.Router, deps dependencies.Dependencies) {
	Endpoints.Activate(r, deps)
}
