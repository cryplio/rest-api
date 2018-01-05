package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Nivl/go-rest-tools/types/apierror"
)

// ResponseError represents the data sent the client when an error occurs
type ResponseError struct {
	Error string `json:"error,omitempty"`
	Field string `json:"field,omitempty"`
}

// HTTPResponse represents an http response
type HTTPResponse interface {
	// Header returns the header map that will be sent by WriteHeader
	Header() http.Header

	// NoContent sends a http.StatusNoContent response
	NoContent()

	// Created sends a http.StatusCreated response with a JSON object attached
	Created(obj interface{}) error
	Ok(obj interface{}) error
}

// Response is a basic implementation of the HTTPResponse that uses a ResponseWriter
type Response struct {
	writer http.ResponseWriter
	deps   *Dependencies
}

// NewResponse creates a new response
func NewResponse(writer http.ResponseWriter, deps *Dependencies) *Response {
	return &Response{
		writer: writer,
		deps:   deps,
	}
}

// Header returns the header object of the response
func (res *Response) Header() http.Header {
	return res.writer.Header()
}

// NoContent sends a http.StatusNoContent response
func (res *Response) NoContent() {
	res.writer.WriteHeader(http.StatusNoContent)
}

// Created sends a http.StatusCreated response with a JSON object attached
func (res *Response) Created(obj interface{}) error {
	return res.renderJSON(http.StatusCreated, obj)
}

// Ok sends a http.StatusOK response with a JSON object attached
func (res *Response) Ok(obj interface{}) error {
	return res.renderJSON(http.StatusOK, obj)
}

// renderJSON attaches a json object to the response
func (res *Response) renderJSON(code int, obj interface{}) error {
	res.setJSON(code)

	if obj != nil {
		return json.NewEncoder(res.writer).Encode(obj)
	}
	return nil
}

// Error sends an error to the client
// If the error is an instance of HTTPError, the returned code will
// match HTTPError.HTTPStatus(). It returns a 500 if no code has been set.
func (res *Response) Error(e error, req HTTPRequest) {
	err := apierror.Convert(e)
	res.errorJSON(err)

	// if the error has a field attached we log it
	field := ""
	if err.Field() != "" {
		field = fmt.Sprintf(`, field: "%s"`, err.Field())
	}
	req.Logger().Errorf(`code: "%d"%s, message: "%s", %s`, err.HTTPStatus(), field, err.Error(), req)

	// We send a report for all server errors
	if err.HTTPStatus() == http.StatusInternalServerError {
		if req.Reporter() != nil {
			req.Reporter().ReportError(err)
		}
	}
}

// errorJSON set the request content to the specified error message and HTTP code.
// The error message should be valid json.
func (res *Response) errorJSON(err apierror.Error) {
	if err.Error() == "" {
		res.writer.WriteHeader(err.HTTPStatus())
		return
	}
	resError := &ResponseError{
		Error: err.Error(),
		Field: err.Field(),
	}

	if err.HTTPStatus() == http.StatusInternalServerError {
		resError.Error = "Something went wrong"
		resError.Field = ""
	}
	res.renderJSON(err.HTTPStatus(), resError)
}

// setJSON set the response to JSON and with the specify HTTP code.
func (res *Response) setJSON(code int) {
	res.writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.writer.Header().Set("X-Content-Type-Options", "nosniff")
	res.writer.WriteHeader(code)
}
