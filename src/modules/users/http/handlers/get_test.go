package handlers_test

import (
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
	"github.com/cryplio/rest-api/src/modules/users"
	"github.com/cryplio/rest-api/src/modules/users/http/handlers"
	"github.com/cryplio/rest-api/src/modules/users/testusers"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetInvalidParams(t *testing.T) {
	t.Parallel()

	testCases := []testguard.InvalidParamsTestCase{
		{
			Description: "Should fail on missing ID",
			MsgMatch:    params.ErrMsgMissingParameter,
			FieldName:   "id",
			Sources: map[string]url.Values{
				"url": url.Values{},
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
			},
		},
	}

	g := handlers.Endpoints[handlers.EndpointGet].Guard
	testguard.InvalidParams(t, g, testCases)
}

func TestGetValidParams(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		description string
		sources     map[string]url.Values
	}{
		{
			"Should work with a valid ID",
			map[string]url.Values{
				"url": url.Values{
					"id": []string{"48d0c8b8-d7a3-4855-9d90-29a06ef474b0"},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			endpts := handlers.Endpoints[handlers.EndpointGet]
			data, err := endpts.Guard.ParseParams(tc.sources, nil)
			assert.NoError(t, err)

			if data != nil {
				p := data.(*handlers.GetParams)
				assert.Equal(t, tc.sources["url"].Get("id"), p.ID)
			}
		})
	}
}

func TestGetOthersData(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userToGet := testusers.NewProfile()
	handlerParams := &handlers.GetParams{
		ID: userToGet.ID,
	}
	requester := &auth.User{
		ID: "48d0c8b8-d7a3-4855-9d90-29a06ef474b0",
	}

	// Mock the database & add expectations
	mockDB := mocksqldb.NewMockConnection(mockCtrl)
	mockDB.QEXPECT().GetSuccess(&users.Profile{}, func(user *users.Profile, stmt, id string) {
		*user = *userToGet
	})

	// Mock the response & add expectations
	res := mockrouter.NewMockHTTPResponse(mockCtrl)
	res.EXPECT().OkSuccess(&users.ProfilePayload{}, func(pld *users.ProfilePayload) {
		assert.Equal(t, userToGet.User.ID, pld.ID, "The user ID should not have changed")
		assert.Equal(t, userToGet.Name, pld.Name, "Name should not have changed")
		assert.Empty(t, pld.Email, "the email should not be returned to anyone")
		assert.False(t, pld.IsAdmin, "user should not be an admin")
	})

	// Mock the request & add expectations
	req := mockrouter.NewMockHTTPRequest(mockCtrl)
	req.EXPECT().Response().Return(res)
	req.EXPECT().Params().Return(handlerParams)
	req.EXPECT().User().Return(requester).Times(2)

	// call the handler
	err := handlers.Get(req, &router.Dependencies{DB: mockDB})
	assert.NoError(t, err, "the handler should not have fail")
}

func TestGetOwnData(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	handlerParams := &handlers.GetParams{
		ID: "0c2f0713-3f9b-4657-9cdd-2b4ed1f214e9",
	}
	requester := &auth.User{
		ID:      handlerParams.ID,
		Name:    "user name",
		Email:   "email@domain.tld",
		IsAdmin: false,
	}

	// Mock the database & add expectations
	mockDB := mocksqldb.NewMockConnection(mockCtrl)
	mockDB.QEXPECT().GetSuccess(&users.Profile{}, func(profile *users.Profile, stmt, id string) {
		*profile = *(testusers.NewProfile())
		profile.User = requester
		profile.UserID = requester.ID
	})

	// Mock the response & add expectations
	res := mockrouter.NewMockHTTPResponse(mockCtrl)
	res.EXPECT().OkSuccess(&users.ProfilePayload{}, func(pld *users.ProfilePayload) {
		assert.Equal(t, requester.ID, pld.ID, "ID should have not changed")
		assert.Equal(t, requester.Name, pld.Name, "Name should have not changed")
		assert.Equal(t, requester.Email, pld.Email, "the email should be returned")
		assert.False(t, pld.IsAdmin, "user should not be an admin")
	})

	// Mock the request & add expectations
	req := mockrouter.NewMockHTTPRequest(mockCtrl)
	req.EXPECT().Response().Return(res)
	req.EXPECT().Params().Return(handlerParams)
	req.EXPECT().User().Return(requester).Times(2)

	// call the handler
	err := handlers.Get(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.NoError(t, err, "the handler should not have fail")
}

func TestGetUnexistingUser(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	handlerParams := &handlers.GetParams{
		ID: "0c2f0713-3f9b-4657-9cdd-2b4ed1f214e9",
	}

	// Mock the database & add expectations
	mockDB := mocksqldb.NewMockConnection(mockCtrl)
	mockDB.QEXPECT().GetNotFound(&users.Profile{})

	// Mock the request & add expectations
	req := mockrouter.NewMockHTTPRequest(mockCtrl)
	req.EXPECT().Params().Return(handlerParams)

	// call the handler
	err := handlers.Get(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusNotFound, httpErr.HTTPStatus())
}
