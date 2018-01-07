package users

import (
	"net/http"

	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard"
)

var getEndpoint = &router.Endpoint{
	Verb:    http.MethodGet,
	Path:    "/users/{id}",
	Handler: Get,
	Guard: &guard.Guard{
		ParamStruct: &GetParams{},
	},
}

// GetParams represent the request params accepted by HandlerGet
type GetParams struct {
	ID string `from:"url" json:"id" params:"uuid,required"`
}

// Get represent an API handler to get a user
func Get(req router.HTTPRequest, deps *router.Dependencies) error {
	params := req.Params().(*GetParams)

	profile, err := GetByIDWithProfile(deps.DB, params.ID)
	if err != nil {
		return err
	}

	var pld *ProfilePayload

	if req.User() != nil && req.User().ID == params.ID {
		pld = profile.ExportPrivate()
	} else {
		pld = profile.ExportPublic()
	}

	return req.Response().Ok(pld)
}
