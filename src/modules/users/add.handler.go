package users

import (
	"net/http"

	"github.com/Nivl/go-hasher/implementations/bcrypt"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard"
	"github.com/Nivl/go-rest-tools/security/auth"
)

var addEndpoint = &router.Endpoint{
	Verb:    http.MethodPost,
	Path:    "/users",
	Handler: Add,
	Guard: &guard.Guard{
		ParamStruct: &AddParams{},
	},
}

// AddParams represents the params accepted by Add to create a new user
type AddParams struct {
	Name     string `from:"form" json:"name" params:"required,trim"`
	Email    string `from:"form" json:"email" params:"required,trim"`
	Password string `from:"form" json:"password" params:"required,trim"`
}

// Add is a HTTP handler used to add a new user
func Add(req router.HTTPRequest, deps *router.Dependencies) error {
	params := req.Params().(*AddParams)
	hashr := &bcrypt.Bcrypt{} // todo (melvin): Move to DI

	encryptedPassword, err := hashr.Hash(params.Password)
	if err != nil {
		return err
	}

	user := &auth.User{
		Name:     params.Name,
		Email:    params.Email,
		Password: encryptedPassword,
	}

	// Create a transaction to keep the user and the profile in sync
	tx, err := deps.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Create the user
	if err := user.Create(tx); err != nil {
		return err
	}
	// Create the profile
	profile := &Profile{User: user, UserID: user.ID}
	if err := profile.Create(tx); err != nil {
		return err
	}

	// Persist the data
	if err := tx.Commit(); err != nil {
		return err
	}
	return req.Response().Created(profile.ExportPrivate())
}
