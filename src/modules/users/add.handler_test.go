package users_test

import (
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/Nivl/go-params"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard/testguard"
	"github.com/Nivl/go-rest-tools/router/mockrouter"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/Nivl/go-sqldb/implementations/mocksqldb"
	"github.com/cryplio/rest-api/src/helpers"
	"github.com/cryplio/rest-api/src/modules/users"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestInvalidParams(t *testing.T) {
	t.Parallel()

	testCases := []testguard.InvalidParamsTestCase{
		{
			Description: "Should fail on missing name",
			MsgMatch:    params.ErrMsgMissingParameter,
			FieldName:   "name",
			Sources: map[string]url.Values{
				"form": url.Values{
					"email":    []string{"email@valid.tld"},
					"password": []string{"password"},
				},
			},
		},
		{
			Description: "Should fail on missing email",
			MsgMatch:    params.ErrMsgMissingParameter,
			FieldName:   "email",
			Sources: map[string]url.Values{
				"form": url.Values{
					"name":     []string{"username"},
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
					"name":  []string{"username"},
					"email": []string{"email@valid.tld"},
				},
			},
		},
	}

	g := users.Endpoints[users.EndpointAdd].Guard
	testguard.InvalidParams(t, g, testCases)
}

func TestAddHappyPath(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	handlerParams := &users.AddParams{
		Name:     "username",
		Email:    "email@domain.tld",
		Password: "valid password",
	}

	// Mock the database & add expectations
	mockDB := mocksqldb.NewMockConnection(mockCtrl)
	tx, _ := mockDB.EXPECT().TransactionSuccess(mockCtrl)
	tx.QEXPECT().InsertSuccess(&auth.User{})
	tx.QEXPECT().InsertSuccess(&users.Profile{})
	tx.EXPECT().CommitSuccess()
	tx.EXPECT().RollbackSuccess()

	// Mock the response & add expectations
	res := mockrouter.NewMockHTTPResponse(mockCtrl)
	res.EXPECT().CreatedSuccess(&users.ProfilePayload{}, func(user *users.ProfilePayload) {
		assert.Equal(t, handlerParams.Name, user.Name)
		assert.Equal(t, handlerParams.Email, user.Email)
		assert.NotEmpty(t, user.ID)
		assert.False(t, user.IsAdmin)
	})

	// Mock the request & add expectations
	req := mockrouter.NewMockHTTPRequest(mockCtrl)
	req.EXPECT().Response().Return(res)
	req.EXPECT().Params().Return(handlerParams)

	// call the handler
	err := users.Add(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Nil(t, err, "the handler should not have fail")
}

func TestAddConflict(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	handlerParams := &users.AddParams{
		Name:     "username",
		Email:    "email@domain.tld",
		Password: "valid password",
	}

	// Mock the database & add expectations
	mockDB := mocksqldb.NewMockConnection(mockCtrl)
	tx, _ := mockDB.EXPECT().TransactionSuccess(mockCtrl)
	tx.QEXPECT().InsertError(&auth.User{}, helpers.SQLConflictError("email", handlerParams.Email))
	tx.EXPECT().RollbackSuccess()

	// Mock the request & add expectations
	req := mockrouter.NewMockHTTPRequest(mockCtrl)
	req.EXPECT().Params().Return(handlerParams)

	// call the handler
	err := users.Add(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")

	apiError := apierror.Convert(err)
	assert.Equal(t, http.StatusConflict, apiError.HTTPStatus())
	assert.Equal(t, "email", apiError.Field())
}

func TestAddProfileError(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	handlerParams := &users.AddParams{
		Name:     "username",
		Email:    "email@domain.tld",
		Password: "valid password",
	}

	// Mock the database & add expectations
	mockDB := mocksqldb.NewMockConnection(mockCtrl)
	tx, _ := mockDB.EXPECT().TransactionSuccess(mockCtrl)
	tx.QEXPECT().InsertSuccess(&auth.User{})
	tx.QEXPECT().InsertError(&users.Profile{}, errors.New("could not create profile"))
	tx.EXPECT().RollbackSuccess()

	// Mock the request & add expectations
	req := mockrouter.NewMockHTTPRequest(mockCtrl)
	req.EXPECT().Params().Return(handlerParams)

	// call the handler
	err := users.Add(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")

	apiError := apierror.Convert(err)
	assert.Equal(t, http.StatusInternalServerError, apiError.HTTPStatus())
}

func TestAddCommitError(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	handlerParams := &users.AddParams{
		Name:     "username",
		Email:    "email@domain.tld",
		Password: "valid password",
	}

	// Mock the database & add expectations
	mockDB := mocksqldb.NewMockConnection(mockCtrl)
	tx, _ := mockDB.EXPECT().TransactionSuccess(mockCtrl)
	tx.QEXPECT().InsertSuccess(&auth.User{})
	tx.QEXPECT().InsertSuccess(&users.Profile{})
	tx.EXPECT().CommitError()
	tx.EXPECT().RollbackSuccess()

	// Mock the request & add expectations
	req := mockrouter.NewMockHTTPRequest(mockCtrl)
	req.EXPECT().Params().Return(handlerParams)

	// call the handler
	err := users.Add(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")

	apiError := apierror.Convert(err)
	assert.Equal(t, http.StatusInternalServerError, apiError.HTTPStatus())
}

func TestAddTransactionError(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	handlerParams := &users.AddParams{
		Name:     "username",
		Email:    "email@domain.tld",
		Password: "valid password",
	}

	// Mock the database & add expectations
	mockDB := mocksqldb.NewMockConnection(mockCtrl)
	mockDB.EXPECT().TransactionError()

	// Mock the request & add expectations
	req := mockrouter.NewMockHTTPRequest(mockCtrl)
	req.EXPECT().Params().Return(handlerParams)

	// call the handler
	err := users.Add(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")

	apiError := apierror.Convert(err)
	assert.Equal(t, http.StatusInternalServerError, apiError.HTTPStatus())
}
