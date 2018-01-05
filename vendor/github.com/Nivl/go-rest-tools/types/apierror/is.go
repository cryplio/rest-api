package apierror

import "net/http"

// IsNotFound checks if an error is the NotFound type
func IsNotFound(e error) bool {
	err, casted := e.(*APIError)
	if !casted {
		return false
	}
	return err.HTTPStatus() == http.StatusNotFound
}

// IsConflict checks if an error is caused by a conflict
func IsConflict(e error) bool {
	err, casted := e.(*APIError)
	if !casted {
		return false
	}
	return err.HTTPStatus() == http.StatusConflict
}

// IsInternalServerError checks if an error is caused by an internal error
func IsInternalServerError(e error) bool {
	err := Convert(e)
	return err.HTTPStatus() == http.StatusInternalServerError
}

// IsBadRequest checks if an error is caused by a bad request
func IsBadRequest(e error) bool {
	err, casted := e.(*APIError)
	if !casted {
		return false
	}
	return err.HTTPStatus() == http.StatusBadRequest
}

// IsForbidden checks if an error is caused by a forbidden access
func IsForbidden(e error) bool {
	err, casted := e.(*APIError)
	if !casted {
		return false
	}
	return err.HTTPStatus() == http.StatusForbidden
}

// IsUnauthorized checks if an error is caused by an Unauthorized access
func IsUnauthorized(e error) bool {
	err, casted := e.(*APIError)
	if !casted {
		return false
	}
	return err.HTTPStatus() == http.StatusUnauthorized
}

// IsInvalidParam checks if an error is caused by an invalid param
func IsInvalidParam(e error) bool {
	return IsBadRequest(e)
}
