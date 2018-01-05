package apierror

import (
	"net/http"
)

// Convert takes an error an turns it into an APIError
func Convert(e error) *APIError {
	err, casted := e.(*APIError)
	if !casted {
		err = NewServerError(e.Error())
		err.ErrorOrigin = e
	}

	return err
}

// APIError represents an error with an HTTP code
type APIError struct {
	error
	ErrorStatus int
	ErrorField  string
	ErrorOrigin error
}

// HTTPStatus returns the HTTP code associated to the error
func (err *APIError) HTTPStatus() int {
	if err == nil {
		return http.StatusInternalServerError
	}
	return err.ErrorStatus
}

// Field returns the HTTP param associated to the error
func (err *APIError) Field() string {
	return err.ErrorField
}

// Origin returns the original error
func (err *APIError) Origin() error {
	return err.ErrorOrigin
}
