package router

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	logger "github.com/Nivl/go-logger"
	reporter "github.com/Nivl/go-reporter"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/gorilla/mux"
)

const (
	// ContentTypeJSON represents the content type of a JSON request
	ContentTypeJSON = "application/json"
	// ContentTypeMultipartForm represents the content type of a multipart request
	ContentTypeMultipartForm = "multipart/form-data"
	// ContentTypeForm represents the content type of a POST/PUT/PATCH request
	ContentTypeForm = "application/x-www-form-urlencoded"
)

// HTTPRequest represents an http request
type HTTPRequest interface {
	String() string

	// Logger returns the logger used by the request
	Logger() logger.Logger

	// Reporter returns the reporter used by the request
	Reporter() reporter.Reporter

	// Signature returns the signature of the request
	// Ex. POST /users
	Signature() string

	// ID returns the ID of the request
	ID() string

	// Response returns the response of the request
	Response() HTTPResponse

	// Params returns the params needed by the endpoint
	Params() interface{}

	// User returns the user that made the request
	User() *auth.User

	// Session returns the session used to make the request
	Session() *auth.Session
}

// Request represent a client request
type Request struct {
	id           string
	res          *Response
	http         *http.Request
	params       interface{}
	user         *auth.User
	session      *auth.Session
	_contentType string
	logger       logger.Logger
	reporter     reporter.Reporter
}

// User returns the user that made the request
func (req *Request) User() *auth.User {
	return req.user
}

// Session returns the session used to make the request
func (req *Request) Session() *auth.Session {
	return req.session
}

// ID returns the ID of the request
func (req *Request) ID() string {
	return req.id
}

// Signature returns the signature of the request
// Ex. POST /users
func (req *Request) Signature() string {
	return fmt.Sprintf("%s %s", req.http.Method, req.http.RequestURI)
}

// Response returns the response of the request
func (req *Request) Response() HTTPResponse {
	return req.res
}

// Logger returns the logger used by the request
func (req *Request) Logger() logger.Logger {
	return req.logger
}

// Reporter returns the reporter used by the request
func (req *Request) Reporter() reporter.Reporter {
	return req.reporter
}

func (req *Request) String() string {
	user := "anonymous"
	userID := "0"
	if req.user != nil {
		user = req.user.Name
		userID = req.user.ID
	}

	return fmt.Sprintf(`req_id: "%s", user: "%s", user_id: "%s", endpoint: "%s", params: %#v`,
		req.id, user, userID, req.Signature(), req.params)
}

// Params returns the params needed by the endpoint
func (req *Request) Params() interface{} {
	return req.params
}

// muxVariables returns the URL variables associated to the request
func (req *Request) muxVariables() url.Values {
	output := url.Values{}

	if req == nil {
		return output
	}

	vars := mux.Vars(req.http)
	for k, v := range vars {
		output.Set(k, v)
	}

	return output
}

// contentType returns the content type of the current request
func (req *Request) contentType() string {
	if req == nil {
		return ""
	}

	if req._contentType == "" {
		contentType := req.http.Header.Get("Content-Type")

		if contentType != "" {
			req._contentType = strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
		}
	}

	return req._contentType
}

// parseJSONBody parses and returns the body of the request
func (req *Request) parseJSONBody() (url.Values, error) {
	output := url.Values{}

	if req.contentType() != ContentTypeJSON {
		return output, nil
	}

	vars := map[string]interface{}{}
	if err := json.NewDecoder(req.http.Body).Decode(&vars); err != nil && err != io.EOF {
		return nil, err
	}

	for k, v := range vars {
		output.Set(k, fmt.Sprintf("%v", v))
	}

	return output, nil
}

// httpParamsBySource returns a map of all http params ordered by their source (url, query, form, ...)
func (req *Request) httpParamsBySource() (map[string]url.Values, error) {
	params := map[string]url.Values{
		"url":   req.muxVariables(),
		"query": req.http.URL.Query(),
		"form":  url.Values{},
	}

	if req.contentType() == ContentTypeJSON {
		form, err := req.parseJSONBody()
		if err != nil {
			return nil, err
		}
		params["form"] = form
	} else if req.contentType() == ContentTypeForm || req.contentType() == ContentTypeMultipartForm {
		if err := req.http.ParseForm(); err != nil {
			return nil, err
		}
		params["form"] = req.http.PostForm
	}

	return params, nil
}

// handlePanic will recover a panic an log what happen
func (req *Request) handlePanic() {
	if rec := recover(); rec != nil {
		// The recovered panic may not be an error
		var err error
		switch val := rec.(type) {
		case error:
			err = val
		default:
			err = fmt.Errorf("%v", val)
		}
		err = fmt.Errorf("panic: %v", err)

		req.res.Error(err, req)
	}
}
