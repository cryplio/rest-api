package handlers_test

import (
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/dchest/uniuri"

	"github.com/Nivl/go-params"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard/testguard"
	"github.com/Nivl/go-rest-tools/router/mockrouter"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/security/auth/testauth"
	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/Nivl/go-sqldb/implementations/mocksqldb"
	"github.com/cryplio/rest-api/src/modules/portfolios"
	"github.com/cryplio/rest-api/src/modules/portfolios/http/handlers"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAddInvalidParams(t *testing.T) {
	t.Parallel()

	testCases := []testguard.InvalidParamsTestCase{
		{
			Description: "Should fail on missing name",
			MsgMatch:    params.ErrMsgMissingParameter,
			FieldName:   "name",
			Sources: map[string]url.Values{
				"form": url.Values{},
			},
		},
		{
			Description: "Should fail on name too long",
			MsgMatch:    params.ErrMsgMaxLen,
			FieldName:   "name",
			Sources: map[string]url.Values{
				"form": url.Values{
					"name": []string{uniuri.NewLen(51)},
				},
			},
		},
	}

	g := handlers.Endpoints[handlers.EndpointAdd].Guard
	testguard.InvalidParams(t, g, testCases)
}

func TestAddAccess(t *testing.T) {
	t.Parallel()

	testCases := []testguard.AccessTestCase{
		{
			Description: "Should fail for anonymous users",
			User:        nil,
			ErrCode:     http.StatusUnauthorized,
		},
		{
			Description: "Should work for logged users",
			User:        &auth.User{ID: "48d0c8b8-d7a3-4855-9d90-29a06ef474b0"},
			ErrCode:     0,
		},
	}

	g := handlers.Endpoints[handlers.EndpointAdd].Guard
	testguard.AccessTest(t, g, testCases)
}

func TestAddHappyPath(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	user := testauth.NewUser()
	handlerParams := &handlers.AddParams{
		Name: "my portfolio",
	}

	// Mock the database & add expectations
	mockDB := mocksqldb.NewMockConnection(mockCtrl)
	mockDB.QEXPECT().InsertSuccess(&portfolios.Portfolio{})

	// Mock the response & add expectations
	res := mockrouter.NewMockHTTPResponse(mockCtrl)
	res.EXPECT().CreatedSuccess(&portfolios.Payload{}, func(p *portfolios.Payload) {
		assert.Equal(t, handlerParams.Name, p.Name)
		assert.NotEmpty(t, p.ID)
	})

	// Mock the request & add expectations
	req := mockrouter.NewMockHTTPRequest(mockCtrl)
	req.EXPECT().Response().Return(res)
	req.EXPECT().Params().Return(handlerParams)
	req.EXPECT().User().Return(user).Times(2)

	// call the handler
	err := handlers.Add(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Nil(t, err, "the handler should not have fail")
}

func TestAddCreateError(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	user := testauth.NewUser()
	handlerParams := &handlers.AddParams{
		Name: "my portfolio",
	}

	// Mock the database & add expectations
	mockDB := mocksqldb.NewMockConnection(mockCtrl)
	mockDB.QEXPECT().InsertError(&portfolios.Portfolio{}, errors.New("could not create portfolio"))

	// Mock the request & add expectations
	req := mockrouter.NewMockHTTPRequest(mockCtrl)
	req.EXPECT().Params().Return(handlerParams)
	req.EXPECT().User().Return(user).Times(2)

	// call the handler
	err := handlers.Add(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")

	apiError := apierror.Convert(err)
	assert.Equal(t, http.StatusInternalServerError, apiError.HTTPStatus())
}
