package handlers

import (
	"net/http"

	"github.com/Nivl/go-hasher/implementations/bcrypt"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard"
	"github.com/Nivl/go-rest-tools/types/apierror"
)

var deleteEndpoint = &router.Endpoint{
	Verb:    http.MethodDelete,
	Path:    "/users/{id}",
	Handler: Delete,
	Guard: &guard.Guard{
		ParamStruct: &DeleteParams{},
		Auth:        guard.LoggedUserAccess,
	},
}

// DeleteParams represent the request params accepted by HandlerDelete
type DeleteParams struct {
	ID              string `from:"url" json:"id" params:"uuid,required"`
	CurrentPassword string `from:"form" json:"current_password" params:"required,trim"`
}

// Delete represent an API handler to remove a user
func Delete(req router.HTTPRequest, deps *router.Dependencies) error {
	params := req.Params().(*DeleteParams)
	user := req.User()
	hashr := &bcrypt.Bcrypt{} // todo (melvin): Move to DI

	if params.ID != user.ID {
		return apierror.NewForbidden()
	}

	if !hashr.IsValid(user.Password, params.CurrentPassword) {
		return apierror.NewUnauthorized()
	}

	if err := user.Delete(deps.DB); err != nil {
		return err
	}

	req.Response().NoContent()
	return nil
}
