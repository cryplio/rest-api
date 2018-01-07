package sessions_test

import (
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/Nivl/go-hasher/implementations/bcrypt"
	"github.com/Nivl/go-params"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard/testguard"
	"github.com/Nivl/go-rest-tools/router/mockrouter"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/Nivl/go-sqldb/implementations/mocksqldb"
	"github.com/cryplio/rest-api/src/modules/sessions"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestInvalidParams(t *testing.T) {
	t.Parallel()

	testCases := []testguard.InvalidParamsTestCase{
		{
			Description: "Should fail on missing email",
			MsgMatch:    params.ErrMsgMissingParameter,
			FieldName:   "email",
			Sources: map[string]url.Values{
				"form": url.Values{
					"password": []string{"password"},
				},
			},
		},
		{
			Description: "Should fail on missing password",
			MsgMatch:    params.ErrMsgMissingParameter,
			FieldName:   "password",
			Sources: map[string]url.Values{
				"form": url.Values{
					"email": []string{"email@valid.tld"},
				},
			},
		},
	}

	g := sessions.Endpoints[sessions.EndpointAdd].Guard
	testguard.InvalidParams(t, g, testCases)
}

func TestAddValidData(t *testing.T) {
	t.Parallel()

	hashr := &bcrypt.Bcrypt{}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	handlerParams := &sessions.AddParams{
		Email:    "email@domain.tld",
		Password: "valid password",
	}

	// Mock the database & add expectations
	mockDB := mocksqldb.NewMockConnection(mockCtrl)
	mockDB.QEXPECT().InsertSuccess(&auth.Session{})
	mockDB.QEXPECT().GetSuccess(&auth.User{}, func(u *auth.User, query string, email string) {
		u.ID = "0c2f0713-3f9b-4657-9cdd-2b4ed1f214e9"

		var err error
		u.Password, err = hashr.Hash(handlerParams.Password)
		assert.NoError(t, err)
	})

	// Mock the response & add expectations
	res := mockrouter.NewMockHTTPResponse(mockCtrl)
	res.EXPECT().CreatedSuccess(&sessions.Payload{}, func(s *sessions.Payload) {
		assert.Equal(t, "0c2f0713-3f9b-4657-9cdd-2b4ed1f214e9", s.UserID)
		assert.NotEmpty(t, s.Token)
	})

	// Mock the request & add expectations
	req := mockrouter.NewMockHTTPRequest(mockCtrl)
	req.EXPECT().Response().Return(res)
	req.EXPECT().Params().Return(handlerParams)

	// call the handler
	err := sessions.Add(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.NoError(t, err, "the handler should not have fail")
}

func TestAddUnexistingEmail(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	handlerParams := &sessions.AddParams{
		Email:    "email@domain.tld",
		Password: "valid password",
	}

	// Mock the database & add expectations
	mockDB := mocksqldb.NewMockConnection(mockCtrl)
	mockDB.QEXPECT().GetNotFound(&auth.User{})

	// Mock the request & add expectations
	req := mockrouter.NewMockHTTPRequest(mockCtrl)
	req.EXPECT().Params().Return(handlerParams)

	// call the handler
	err := sessions.Add(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err)

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusBadRequest, httpErr.HTTPStatus())
	assert.Equal(t, "email/password", httpErr.Field())
}

func TestAddWrongPassword(t *testing.T) {
	t.Parallel()

	hashr := &bcrypt.Bcrypt{}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	handlerParams := &sessions.AddParams{
		Email:    "email@domain.tld",
		Password: "invalid password",
	}

	// Mock the database & add expectations
	mockDB := mocksqldb.NewMockConnection(mockCtrl)
	mockDB.QEXPECT().GetSuccess(&auth.User{}, func(u *auth.User, query string, email string) {
		u.ID = "0c2f0713-3f9b-4657-9cdd-2b4ed1f214e9"

		var err error
		u.Password, err = hashr.Hash("valid password")
		assert.NoError(t, err)
	})

	// Mock the request & add expectations
	req := mockrouter.NewMockHTTPRequest(mockCtrl)
	req.EXPECT().Params().Return(handlerParams)

	// call the handler
	err := sessions.Add(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err)

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusBadRequest, httpErr.HTTPStatus())
	assert.Equal(t, "email/password", httpErr.Field())
}

func TestAddNoDbConOnGet(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	handlerParams := &sessions.AddParams{
		Email:    "email@domain.tld",
		Password: "invalid password",
	}

	// Mock the database & add expectations
	mockDB := mocksqldb.NewMockConnection(mockCtrl)
	mockDB.QEXPECT().GetError(&auth.User{}, errors.New("system error"))

	// Mock the request & add expectations
	req := mockrouter.NewMockHTTPRequest(mockCtrl)
	req.EXPECT().Params().Return(handlerParams)

	// call the handler
	err := sessions.Add(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err)

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusInternalServerError, httpErr.HTTPStatus())
}

func TestAddNoDBConOnSave(t *testing.T) {
	t.Parallel()

	hashr := &bcrypt.Bcrypt{}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	handlerParams := &sessions.AddParams{
		Email:    "email@domain.tld",
		Password: "valid password",
	}

	// Mock the database & add expectations
	mockDB := mocksqldb.NewMockConnection(mockCtrl)
	mockDB.QEXPECT().InsertError(&auth.Session{}, errors.New("system error"))
	mockDB.QEXPECT().GetSuccess(&auth.User{}, func(u *auth.User, query string, email string) {
		u.ID = "0c2f0713-3f9b-4657-9cdd-2b4ed1f214e9"

		var err error
		u.Password, err = hashr.Hash(handlerParams.Password)
		assert.NoError(t, err)
	})

	// Mock the request & add expectations
	req := mockrouter.NewMockHTTPRequest(mockCtrl)
	req.EXPECT().Params().Return(handlerParams)

	// call the handler
	err := sessions.Add(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusInternalServerError, httpErr.HTTPStatus())
}
