package users

import (
	"net/http"

	"github.com/Nivl/go-hasher"
	"github.com/Nivl/go-hasher/implementations/bcrypt"

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

// UpdateParams represents the params accepted Update to update a user
type UpdateParams struct {
	ID              string `from:"url" json:"id"  params:"uuid,required"`
	Name            string `from:"form" json:"name" params:"trim" maxlen:"255"`
	Email           string `from:"form" json:"email" params:"trim,email" maxlen:"255"`
	CurrentPassword string `from:"form" json:"current_password" params:"trim" maxlen:"255"`
	NewPassword     string `from:"form" json:"new_password" params:"trim" maxlen:"255"`
}

// Update is a HTTP handler used to update a user
func Update(req router.HTTPRequest, deps *router.Dependencies) error {
	params := req.Params().(*UpdateParams)
	currentUser := req.User()
	hashr := &bcrypt.Bcrypt{} // todo (melvin): Move to DI

	// Admin are allowed to update any users
	if !currentUser.IsAdmin && params.ID != currentUser.ID {
		return apierror.NewForbidden()
	}

	// Retreive the user and the attached profile
	profile, err := GetByIDWithProfile(deps.DB, params.ID)
	if err != nil {
		return err
	}

	// To change the email or the password we require the current password
	if params.NewPassword != "" || params.Email != "" {
		if !hashr.IsValid(profile.User.Password, params.CurrentPassword) {
			return apierror.NewUnauthorized()
		}
	}

	// Copy the data from the params to the profile
	if err := updateCopyData(profile, params, hashr); err != nil {
		return err
	}

	// Create a transaction to keep the user and the profile in sync
	tx, err := deps.DB.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Save the new data
	if err := profile.User.Save(tx); err != nil {
		return err
	}
	if err := profile.Save(tx); err != nil {
		return err
	}

	// Persist the changes
	if err := tx.Commit(); err != nil {
		return err
	}
	return req.Response().Ok(profile.ExportPrivate())
}

func updateCopyData(profile *Profile, params *UpdateParams, hashr hasher.Hasher) error {
	// Update the User Object
	if params.Name != "" {
		profile.User.Name = params.Name
	}
	if params.Email != "" {
		profile.User.Email = params.Email
	}
	if params.NewPassword != "" {
		hashedPassword, err := hashr.Hash(params.NewPassword)
		if err != nil {
			return err
		}
		profile.User.Password = hashedPassword
	}
	return nil
}
