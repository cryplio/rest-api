package apierror

// Error represents an error with a code attached.
type Error interface {
	error

	// HTTPStatus return the HTTP code of the error
	HTTPStatus() int

	// Field returns the http/sql/etc. field associated to the error
	Field() string

	// Origin returns the original error if there's one
	Origin() error
}
