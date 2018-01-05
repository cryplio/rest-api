package apierror

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/Nivl/go-params/perror"

	"github.com/lib/pq"
)

const (
	// ErrDup contains the errcode of a unique constraint violation
	ErrDup = "23505"
)

// NewFromSQL returns an http error based on a pq.Error
// the provided error will be returned if it's not a pq.Error instance,
// or if the error cannot be matched to
func NewFromSQL(err error) error {
	if err != nil && err == sql.ErrNoRows {
		return NewNotFound()
	}

	if pqErr, ok := err.(*pq.Error); ok {
		switch pqErr.Code {
		case ErrDup:
			// because it's a constraint issue, the column name won't be stored in
			// pqErr.Column. Fortunately we can find it in detail.
			// Example of detail: "Key (name)=(Google) already exists."
			r := regexp.MustCompile(`^Key \(([a-z_]+)\).*\.$`)
			matches := r.FindStringSubmatch(pqErr.Detail)
			fieldName := "unknown"
			if len(matches) > 1 {
				fieldName = matches[1]
			}
			return NewConflict(fieldName)
		default:
			return err
		}
	}
	return err
}

// NewFromError returns an api error based on an error
// the provided error will be returned if it doesn't match any known error
func NewFromError(err error) error {
	err = NewFromSQL(err)

	switch e := err.(type) {
	case perror.Error:
		return NewBadRequest(e.Field(), e.Error())
	}

	return err
}

// NewError returns an error with an associated code
func NewError(code int, field string, message string, args ...interface{}) *APIError {
	fullMessage := fmt.Sprintf(message, args...)
	return &APIError{errors.New(fullMessage), code, field, nil}
}

// NewServerError returns an Internal Error.
func NewServerError(message string, args ...interface{}) *APIError {
	return NewError(http.StatusInternalServerError, "", message, args...)
}

// NewBadRequest returns an error caused by a user. Example: A missing param
func NewBadRequest(field string, message string, args ...interface{}) *APIError {
	return NewError(http.StatusBadRequest, field, message, args...)
}

// NewInvalidParam is an alias for NewBadRequest
func NewInvalidParam(field string, message string, args ...interface{}) *APIError {
	return NewBadRequest(field, message, args)
}

// NewConflict returns an error caused by a conflict with the current state
// of the app. Example: A duplicate slug
func NewConflict(field string) *APIError {
	return NewError(http.StatusConflict, field, "already exists")
}

// NewConflictR returns an error caused by a conflict with the current state
// of the app. A reason is sent back to the user.
func NewConflictR(field string, message string, args ...interface{}) *APIError {
	return NewError(http.StatusConflict, field, message, args...)
}

// NewUnauthorized returns an error caused by a anonymous user trying to access
// a protected resource
func NewUnauthorized() *APIError {
	return NewUnauthorizedR(http.StatusText(http.StatusUnauthorized))
}

// NewUnauthorizedR returns an error caused by a anonymous user trying to access
// a protected resource. A reason is sent back to the user.
func NewUnauthorizedR(reason string) *APIError {
	return NewError(http.StatusUnauthorized, "", reason)
}

// NewForbidden returns an error caused by a user trying to access
// a protected resource.
func NewForbidden() *APIError {
	return NewForbiddenR(http.StatusText(http.StatusForbidden))
}

// NewForbiddenR returns an error caused by a user trying to access
// a protected resource. A reason is sent back to the user.
func NewForbiddenR(reason string) *APIError {
	return NewError(http.StatusForbidden, "", reason)
}

// NewNotFound returns an error caused by a user trying to access
// a resource that does not exists
func NewNotFound() *APIError {
	return NewNotFoundR(http.StatusText(http.StatusNotFound))
}

// NewNotFoundR returns an error caused by a user trying to access
// a resource that does not exists. A reason is sent back to the user.
func NewNotFoundR(reason string) *APIError {
	return NewError(http.StatusNotFound, "", reason)
}

// NewNotFoundField returns an error caused by a user trying to access
// a resource that does not exists. A reason is sent back to the user.
func NewNotFoundField(field string, reason string) *APIError {
	return NewError(http.StatusNotFound, field, reason)
}
