package users_test

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
	"github.com/cryplio/rest-api/src/helpers"
	"github.com/cryplio/rest-api/src/modules/users"
	"github.com/cryplio/rest-api/src/modules/users/testusers"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUpdateInvalidParams(t *testing.T) {
	t.Parallel()

	testCases := []testguard.InvalidParamsTestCase{
		{
			Description: "Should fail on missing ID",
			MsgMatch:    params.ErrMsgMissingParameter,
			FieldName:   "id",
			Sources: map[string]url.Values{
				"url":  url.Values{},
				"form": url.Values{},
			},
		},
		{
			Description: "Should fail on invalid ID",
			MsgMatch:    params.ErrMsgInvalidUUID,
			FieldName:   "id",
			Sources: map[string]url.Values{
				"url": url.Values{
					"id": []string{"not-a-uuid"},
				},
				"form": url.Values{},
			},
		},
		{
			Description: "Should fail on invalid Email",
			MsgMatch:    params.ErrMsgInvalidEmail,
			FieldName:   "email",
			Sources: map[string]url.Values{
				"url": url.Values{
					"id": []string{"48d0c8b8-d7a3-4855-9d90-29a06ef474b0"},
				},
				"form": url.Values{
					"email": []string{"not-an-email"},
				},
			},
		},
	}

	g := users.Endpoints[users.EndpointUpdate].Guard
	testguard.InvalidParams(t, g, testCases)
}

func TestUpdateValidParams(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		description string
		sources     map[string]url.Values
	}{
		{
			"Should work with only a valid ID",
			map[string]url.Values{
				"url": url.Values{
					"id": []string{"48d0c8b8-d7a3-4855-9d90-29a06ef474b0"},
				},
				"form": url.Values{},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			endpts := users.Endpoints[users.EndpointUpdate]
			data, err := endpts.Guard.ParseParams(tc.sources, nil)
			assert.NoError(t, err)

			if data != nil {
				p := data.(*users.UpdateParams)
				assert.Equal(t, tc.sources["url"].Get("id"), p.ID)
			}
		})
	}
}

func TestUpdateAccess(t *testing.T) {
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

	g := users.Endpoints[users.EndpointUpdate].Guard
	testguard.AccessTest(t, g, testCases)
}

func TestUpdateHappyPath(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	profile := testusers.NewProfile()
	handlerParams := &users.UpdateParams{
		ID:              profile.User.ID,
		CurrentPassword: "fake",
		Email:           "new_email@domain.tld",
	}

	// Mock the database & add expectations
	mockDB := mocksqldb.NewMockConnection(mockCtrl)
	mockDB.QEXPECT().GetSuccess(&users.Profile{}, func(p *users.Profile, stmt, id string) {
		*p = *profile
	})
	// mock the transaction
	tx, _ := mockDB.EXPECT().TransactionSuccess(mockCtrl)
	tx.QEXPECT().UpdateSuccess(&auth.User{})
	tx.QEXPECT().UpdateSuccess(&users.Profile{})
	tx.EXPECT().CommitSuccess()
	tx.EXPECT().RollbackSuccess()

	// Mock the response & add expectati ons
	res := mockrouter.NewMockHTTPResponse(mockCtrl)
	res.EXPECT().OkSuccess(&users.ProfilePayload{}, func(data *users.ProfilePayload) {
		assert.Equal(t, profile.User.Name, data.Name, "the name should have not changed")
		assert.Equal(t, handlerParams.Email, data.Email, "email should have been updated")
	})

	// Mock the request & add expectations
	req := mockrouter.NewMockHTTPRequest(mockCtrl)
	req.EXPECT().Response().Return(res)
	req.EXPECT().Params().Return(handlerParams)
	req.EXPECT().User().Return(profile.User)

	// call the handler
	err := users.Update(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.NoError(t, err, "the handler should not have fail")
}

func TestUpdateInvalidPassword(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	profile := testusers.NewProfile()
	handlerParams := &users.UpdateParams{
		ID:              profile.UserID,
		CurrentPassword: "invalid password",
		NewPassword:     "new password",
	}

	// Mock the database & add expectations
	mockDB := mocksqldb.NewMockConnection(mockCtrl)
	mockDB.QEXPECT().GetSuccess(&users.Profile{}, func(p *users.Profile, stmt, id string) {
		*p = *profile
	})

	// Mock the request & add expectations
	req := mockrouter.NewMockHTTPRequest(mockCtrl)
	req.EXPECT().Params().Return(handlerParams)
	req.EXPECT().User().Return(profile.User)

	// call the handler
	err := users.Update(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusUnauthorized, httpErr.HTTPStatus())
}

func TestUpdateInvalidUser(t *testing.T) {
	t.Parallel()

	hashr := bcrypt.Bcrypt{}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	handlerParams := &users.UpdateParams{
		ID:              "48d0c8b8-d7a3-4855-9d90-29a06ef474b0",
		CurrentPassword: "valid password",
	}

	userPassword, err := hashr.Hash("valid password")
	assert.NoError(t, err)
	user := &auth.User{
		ID:       "0c2f0713-3f9b-4657-9cdd-2b4ed1f214e9",
		Password: userPassword,
	}

	// Mock the request & add expectations
	req := mockrouter.NewMockHTTPRequest(mockCtrl)
	req.EXPECT().Params().Return(handlerParams)
	req.EXPECT().User().Return(user)

	// call the handler
	err = users.Update(req, &router.Dependencies{})

	// Assert everything
	assert.Error(t, err, "the handler should not have fail")

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusForbidden, httpErr.HTTPStatus())
}

func TestUpdateUnexistingUser(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	handlerParams := &users.UpdateParams{
		ID: "48d0c8b8-d7a3-4855-9d90-29a06ef474b0",
	}
	requester := &auth.User{
		ID:      "0c2f0713-3f9b-4657-9cdd-2b4ed1f214e9",
		IsAdmin: true,
	}

	// Mock the database & add expectations
	mockDB := mocksqldb.NewMockConnection(mockCtrl)
	mockDB.QEXPECT().GetNotFound(&users.Profile{})

	// Mock the request & add expectations
	req := mockrouter.NewMockHTTPRequest(mockCtrl)
	req.EXPECT().Params().Return(handlerParams)
	req.EXPECT().User().Return(requester)

	// call the handler
	err := users.Update(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusNotFound, httpErr.HTTPStatus())
}

func TestUpdateAllTheFields(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	profile := testusers.NewProfile()
	handlerParams := &users.UpdateParams{
		ID:              profile.User.ID,
		CurrentPassword: "fake",
		Email:           "new_email@domain.tld",
	}

	// Mock the database & add expectations
	mockDB := mocksqldb.NewMockConnection(mockCtrl)
	mockDB.QEXPECT().GetSuccess(&users.Profile{}, func(p *users.Profile, stmt, id string) {
		*p = *profile
	})
	// mock the transaction
	tx, _ := mockDB.EXPECT().TransactionSuccess(mockCtrl)
	tx.QEXPECT().UpdateSuccess(&auth.User{})
	tx.QEXPECT().UpdateSuccess(&users.Profile{})
	tx.EXPECT().CommitSuccess()
	tx.EXPECT().RollbackSuccess()

	// Mock the response & add expectati ons
	res := mockrouter.NewMockHTTPResponse(mockCtrl)
	res.EXPECT().OkSuccess(&users.ProfilePayload{}, func(data *users.ProfilePayload) {
		assert.Equal(t, profile.User.Name, data.Name, "the name should have not changed")
		assert.Equal(t, handlerParams.Email, data.Email, "email should have been updated")
	})

	// Mock the request & add expectations
	req := mockrouter.NewMockHTTPRequest(mockCtrl)
	req.EXPECT().Response().Return(res)
	req.EXPECT().Params().Return(handlerParams)
	req.EXPECT().User().Return(profile.User)

	// call the handler
	err := users.Update(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.NoError(t, err, "the handler should not have fail")
}

func TestUpdateTransactionError(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	profile := testusers.NewProfile()
	handlerParams := &users.UpdateParams{
		ID:              profile.UserID,
		CurrentPassword: "fake",
		NewPassword:     "new password",
	}

	// Mock the database & add expectations
	mockDB := mocksqldb.NewMockConnection(mockCtrl)
	mockDB.QEXPECT().GetSuccess(&users.Profile{}, func(p *users.Profile, stmt, id string) {
		*p = *profile
	})
	mockDB.EXPECT().TransactionError()

	// Mock the request & add expectations
	req := mockrouter.NewMockHTTPRequest(mockCtrl)
	req.EXPECT().Params().Return(handlerParams)
	req.EXPECT().User().Return(profile.User)

	// call the handler
	err := users.Update(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusInternalServerError, httpErr.HTTPStatus())
}

func TestUpdateConflict(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	profile := testusers.NewProfile()
	handlerParams := &users.UpdateParams{
		ID:              profile.UserID,
		CurrentPassword: "fake",
		Name:            "new Name",
	}

	// Mock the database & add expectations
	mockDB := mocksqldb.NewMockConnection(mockCtrl)
	mockDB.QEXPECT().GetSuccess(&users.Profile{}, func(p *users.Profile, stmt, id string) {
		*p = *profile
	})
	// mock the transaction
	tx, _ := mockDB.EXPECT().TransactionSuccess(mockCtrl)
	tx.QEXPECT().UpdateError(&auth.User{}, helpers.SQLConflictError("email", "email@gmail.com"))
	tx.EXPECT().RollbackSuccess()

	// Mock the request & add expectations
	req := mockrouter.NewMockHTTPRequest(mockCtrl)
	req.EXPECT().Params().Return(handlerParams)
	req.EXPECT().User().Return(profile.User)

	// call the handler
	err := users.Update(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")

	apiError := apierror.Convert(err)
	assert.Equal(t, http.StatusConflict, apiError.HTTPStatus())
	assert.Equal(t, "email", apiError.Field())
}

func TestUpdateProfileError(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	profile := testusers.NewProfile()
	handlerParams := &users.UpdateParams{
		ID:              profile.UserID,
		CurrentPassword: "fake",
		Name:            "new Name",
	}

	// Mock the database & add expectations
	mockDB := mocksqldb.NewMockConnection(mockCtrl)
	mockDB.QEXPECT().GetSuccess(&users.Profile{}, func(p *users.Profile, stmt, id string) {
		*p = *profile
	})
	// mock the transaction
	tx, _ := mockDB.EXPECT().TransactionSuccess(mockCtrl)
	tx.QEXPECT().UpdateSuccess(&auth.User{})
	tx.QEXPECT().UpdateError(&users.Profile{}, errors.New("could not update profile"))
	tx.EXPECT().RollbackSuccess()

	// Mock the request & add expectations
	req := mockrouter.NewMockHTTPRequest(mockCtrl)
	req.EXPECT().Params().Return(handlerParams)
	req.EXPECT().User().Return(profile.User)

	// call the handler
	err := users.Update(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")

	apiError := apierror.Convert(err)
	assert.Equal(t, http.StatusInternalServerError, apiError.HTTPStatus())
}

func TestUpdateCommitError(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	profile := testusers.NewProfile()
	handlerParams := &users.UpdateParams{
		ID:              profile.UserID,
		CurrentPassword: "fake",
		Name:            "new Name",
	}

	// Mock the database & add expectations
	mockDB := mocksqldb.NewMockConnection(mockCtrl)
	mockDB.QEXPECT().GetSuccess(&users.Profile{}, func(p *users.Profile, stmt, id string) {
		*p = *profile
	})
	// mock the transaction
	tx, _ := mockDB.EXPECT().TransactionSuccess(mockCtrl)
	tx.QEXPECT().UpdateSuccess(&auth.User{})
	tx.QEXPECT().UpdateSuccess(&users.Profile{})
	tx.EXPECT().CommitError()
	tx.EXPECT().RollbackSuccess()

	// Mock the request & add expectations
	req := mockrouter.NewMockHTTPRequest(mockCtrl)
	req.EXPECT().Params().Return(handlerParams)
	req.EXPECT().User().Return(profile.User)

	// call the handler
	err := users.Update(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")

	apiError := apierror.Convert(err)
	assert.Equal(t, http.StatusInternalServerError, apiError.HTTPStatus())
}
