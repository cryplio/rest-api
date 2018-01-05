package guard

import (
	"net/url"
	"reflect"

	"github.com/Nivl/go-params"
	"github.com/Nivl/go-params/formfile"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/types/apierror"
)

// Guard represents a security access system for routes
type Guard struct {
	// ParamStruct is an instance of a struct that describes the http params
	// accepted by an endpoint
	ParamStruct interface{}

	// Auth is used to add a auth middleware
	Auth RouteAuth
}

// ParseParams parses and returns the list of params needed
// Returns an error if a required param is missing, or if a type is wrong
func (g *Guard) ParseParams(sources map[string]url.Values, fileHolder formfile.FileHolder) (interface{}, error) {
	// It's ok not to have a guard provided, as well as not having params
	if g == nil || g.ParamStruct == nil {
		return nil, nil
	}

	// We give p the same type as g.ParamStruct
	p := reflect.New(reflect.TypeOf(g.ParamStruct).Elem()).Interface()
	err := params.New(p).Parse(sources, fileHolder)
	if err != nil {
		return nil, apierror.NewFromError(err)
	}
	return p, nil
}

// HasAccess check if a given user has access to the
func (g *Guard) HasAccess(u *auth.User) (bool, apierror.Error) {
	// It's ok not to have a guard provided, as well as not having an auth check
	if g == nil || g.Auth == nil {
		return true, nil
	}

	err := g.Auth(u)
	return err == nil, err
}
