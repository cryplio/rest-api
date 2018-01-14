package handlers

import (
	"net/http"

	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard"
	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/cryplio/rest-api/src/modules/portfolios"
)

var deleteEndpoint = &router.Endpoint{
	Verb:    http.MethodDelete,
	Path:    "/portfolio/{id}",
	Handler: Delete,
	Guard: &guard.Guard{
		ParamStruct: &DeleteParams{},
		Auth:        guard.LoggedUserAccess,
	},
}

// DeleteParams represent the request params accepted by HandlerDelete
type DeleteParams struct {
	ID string `from:"url" json:"id" params:"uuid,required"`
}

// Delete represent an API handler to delete a portfolio
func Delete(req router.HTTPRequest, deps *router.Dependencies) error {
	params := req.Params().(*DeleteParams)

	portfolio, err := portfolios.GetPortfolioByID(deps.DB, params.ID)
	if err != nil {
		return err
	}

	if req.User().ID != portfolio.UserID {
		// If a user tries to delete the portfolio of someone else
		// we return a NotFound instead of Forbidden to add a small
		// layer of security
		return apierror.NewNotFound()
	}

	if err := portfolio.Delete(deps.DB); err != nil {
		return err
	}

	req.Response().NoContent()
	return nil
}
