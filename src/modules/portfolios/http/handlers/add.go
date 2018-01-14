package handlers

import (
	"net/http"

	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard"
	"github.com/cryplio/rest-api/src/modules/portfolios"
)

// addEndpoint represent the endpoint data for the Add handler
var addEndpoint = &router.Endpoint{
	Verb:    http.MethodPost,
	Path:    "/portfolios",
	Handler: Add,
	Guard: &guard.Guard{
		ParamStruct: &AddParams{},
		Auth:        guard.LoggedUserAccess,
	},
}

// AddParams represents the params accepted by Add to create a new portfolio
type AddParams struct {
	Name string `from:"form" json:"name" params:"required,trim" maxlen:"50"`
}

// Add is a HTTP handler used to add a new user
func Add(req router.HTTPRequest, deps *router.Dependencies) error {
	params := req.Params().(*AddParams)

	p := &portfolios.Portfolio{
		Name:   params.Name,
		UserID: req.User().ID,
		User:   req.User(),
	}

	if err := p.Create(deps.DB); err != nil {
		return err
	}

	return req.Response().Created(p.ExportPublic())
}
