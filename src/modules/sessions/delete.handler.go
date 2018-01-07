package sessions

import (
	"database/sql"
	"net/http"

	"github.com/Nivl/go-hasher/implementations/bcrypt"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/types/apierror"
)

var deleteEndpoint = &router.Endpoint{
	Verb:    http.MethodDelete,
	Path:    "/sessions/{token}",
	Handler: Delete,
	Guard: &guard.Guard{
		ParamStruct: &DeleteParams{},
		Auth:        guard.LoggedUserAccess,
	},
}

// DeleteParams represent the request params accepted by HandlerDelete
type DeleteParams struct {
	Token           string `from:"url" json:"token" params:"uuid,required"`
	CurrentPassword string `from:"form" json:"current_password" params:"trim"`
}

// Delete represent an API handler to remove a session
func Delete(req router.HTTPRequest, deps *router.Dependencies) error {
	params := req.Params().(*DeleteParams)

	// If a user tries to delete their current session, then no need for the
	// password (that's just a logout)
	if params.Token != req.Session().ID {
		hashr := bcrypt.Bcrypt{}
		if !hashr.IsValid(req.User().Password, params.CurrentPassword) {
			return apierror.NewUnauthorized()
		}
	}

	var session auth.Session
	stmt := "SELECT * FROM sessions WHERE id=$1 AND deleted_at IS NULL LIMIT 1"
	err := deps.DB.Get(&session, stmt, params.Token)
	// in case of not found we don't wan to return an error here
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	// We always return a 404 in case of a user error to avoid brute-force
	if session.ID == "" || session.UserID != req.User().ID {
		return apierror.NewNotFound()
	}

	if err := session.Delete(deps.DB); err != nil {
		return err
	}

	req.Response().NoContent()
	return nil
}
