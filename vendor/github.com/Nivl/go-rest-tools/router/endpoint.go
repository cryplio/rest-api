package router

import "github.com/Nivl/go-rest-tools/router/guard"

// RouteHandler is the function signature we nee
type RouteHandler func(HTTPRequest, *Dependencies) error

// Endpoint represents an HTTP endpoint
type Endpoint struct {
	Verb string

	// Path is the path for the current component
	Path string

	// Handler is the handler to call
	Handler RouteHandler

	// Guard is the security system of an endpoint
	Guard *guard.Guard
}
