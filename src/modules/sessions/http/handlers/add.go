package handlers

import (
	"database/sql"
	"net/http"

	"github.com/Nivl/go-hasher/implementations/bcrypt"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/cryplio/rest-api/src/modules/sessions"
)

var addEndpoint = &router.Endpoint{
	Verb:    http.MethodPost,
	Path:    "/sessions",
	Handler: Add,
	Guard: &guard.Guard{
		ParamStruct: &AddParams{},
	},
}

// AddParams represent the request params accepted by HandlerAdd
type AddParams struct {
	Email    string `from:"form" json:"email" params:"required,trim"`
	Password string `from:"form" json:"password" params:"required,trim"`
}

// Add represents an API handler to create a new user session
func Add(req router.HTTPRequest, deps *router.Dependencies) error {
	params := req.Params().(*AddParams)

	var user auth.User
	stmt := "SELECT * FROM users WHERE email=$1 LIMIT 1"
	err := deps.DB.Get(&user, stmt, params.Email)
	// in case of not found we don't wan to return an error here
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	hashr := bcrypt.Bcrypt{}
	if user.ID == "" || !hashr.IsValid(user.Password, params.Password) {
		return apierror.NewBadRequest("email/password", "Bad email/password")
	}

	s := &auth.Session{
		UserID: user.ID,
	}
	if err := s.Save(deps.DB); err != nil {
		return err
	}

	req.Response().Created(sessions.NewPayload(s))
	return nil
}
