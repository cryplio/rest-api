package guard

import (
	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/Nivl/go-rest-tools/security/auth"
)

// RouteAuth represents a middleware used to allow/block the access to an endpoint
type RouteAuth func(*auth.User) apierror.Error

// LoggedUserAccess is a auth middleware that filters out anonymous users
func LoggedUserAccess(u *auth.User) apierror.Error {
	if u == nil || u.ID == "" {
		return apierror.NewUnauthorized()
	}
	return nil
}

// AdminAccess is a auth middleware that filters out non admin users
func AdminAccess(u *auth.User) apierror.Error {
	if err := LoggedUserAccess(u); err != nil {
		return apierror.NewUnauthorized()
	}
	if !u.IsAdmin {
		return apierror.NewForbidden()
	}
	return nil
}
