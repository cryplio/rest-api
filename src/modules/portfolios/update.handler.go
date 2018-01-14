package portfolios

import (
	"net/http"

	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard"
	"github.com/Nivl/go-rest-tools/types/apierror"
)

var updateEndpoint = &router.Endpoint{
	Verb:    http.MethodPatch,
	Path:    "/users/{id}",
	Handler: Update,
	Guard: &guard.Guard{
		ParamStruct: &UpdateParams{},
		Auth:        guard.LoggedUserAccess,
	},
}

// UpdateParams represents the params accepted Update to update a portfolio
type UpdateParams struct {
	ID   string `from:"url" json:"id"  params:"uuid,required"`
	Name string `from:"form" json:"name" params:"trim" maxlen:"50"`
}

// Update is a HTTP handler used to update a user
func Update(req router.HTTPRequest, deps *router.Dependencies) error {
	params := req.Params().(*UpdateParams)

	portfolio, err := GetPortfolioByID(deps.DB, params.ID)
	if err != nil {
		return err
	}

	// To avoid discovery by brute force we return a not found when someone has
	// no permissions
	if portfolio.UserID != req.User().ID {
		return apierror.NewNotFound()
	}

	if params.Name != "" {
		portfolio.Name = params.Name
	}

	// Save the new data
	if err := portfolio.Save(deps.DB); err != nil {
		return err
	}

	return req.Response().Ok(portfolio.ExportPublic())
}
