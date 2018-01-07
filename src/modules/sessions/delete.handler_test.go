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

func TestDeleteInvalidParams(t *testing.T) {
	t.Parallel()

	testCases := []testguard.InvalidParamsTestCase{
		{
			Description: "Should fail on missing token",
			MsgMatch:    params.ErrMsgMissingParameter,
			FieldName:   "token",
			Sources: map[string]url.Values{
				"url": url.Values{
					"token": []string{""},
				},
				"form": url.Values{},
			},
		},
		{
			Description: "Should fail on invalid token",
			MsgMatch:    params.ErrMsgInvalidUUID,
			FieldName:   "token",
			Sources: map[string]url.Values{
				"url": url.Values{
					"token": []string{"xxx-yyyy"},
				},
				"form": url.Values{},
			},
		},
	}

	g := sessions.Endpoints[sessions.EndpointDelete].Guard
	testguard.InvalidParams(t, g, testCases)
}

func TestDeleteValidParams(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		description string
		sources     map[string]url.Values
	}{
		{
			"Should work with only a valid ID",
			map[string]url.Values{
				"url": url.Values{
					"token": []string{"48d0c8b8-d7a3-4855-9d90-29a06ef474b0"},
				},
				"form": url.Values{},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			endpts := sessions.Endpoints[sessions.EndpointDelete]
			data, err := endpts.Guard.ParseParams(tc.sources, nil)
			assert.NoError(t, err)

			if data != nil {
				p := data.(*sessions.DeleteParams)
				assert.Equal(t, tc.sources["url"].Get("token"), p.Token)
			}
		})
	}
}

func TestDeleteAccess(t *testing.T) {
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

	g := sessions.Endpoints[sessions.EndpointDelete].Guard
	testguard.AccessTest(t, g, testCases)
}

// TestDeleteHappyPath test a user loging out (removing the current session)
func TestDeleteHappyPath(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	user := &auth.User{ID: "3e916798-a090-4f22-b1d1-04a63fbed6ef"}
	session := &auth.Session{ID: "3642e0e6-788e-4161-92dd-6c52ea823da9", UserID: user.ID}
	handlerParams := &sessions.DeleteParams{
		Token: session.ID,
	}

	// Mock the database & add expectations
	mockDB := mocksqldb.NewMockConnection(mockCtrl)
	mockDB.QEXPECT().GetSuccess(&auth.Session{}, func(sess *auth.Session, stmt string, token string) {
		sess.ID = session.ID
		sess.UserID = session.UserID
	})
	// delete call
	mockDB.QEXPECT().DeletionSuccess()

	// Mock the response & add expectations
	res := mockrouter.NewMockHTTPResponse(mockCtrl)
	res.EXPECT().NoContent().Return()

	// Mock the request & add expectations
	req := mockrouter.NewMockHTTPRequest(mockCtrl)
	req.EXPECT().Response().Return(res)
	req.EXPECT().Params().Return(handlerParams)
	req.EXPECT().User().Return(user)
	req.EXPECT().Session().Return(session)

	// call the handler
	err := sessions.Delete(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.NoError(t, err, "the handler should not have fail")
}

func TestDeleteOtherSession(t *testing.T) {
	t.Parallel()

	hashr := &bcrypt.Bcrypt{}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	handlerParams := &sessions.DeleteParams{
		Token:           "d1ba3e60-d674-47f0-9694-59fbda0fc659",
		CurrentPassword: "valid password",
	}

	// Generate a password for the user
	userPassword, err := hashr.Hash(handlerParams.CurrentPassword)
	assert.NoError(t, err)

	// defined the request data
	user := &auth.User{ID: "3e916798-a090-4f22-b1d1-04a63fbed6ef", Password: userPassword}
	session := &auth.Session{ID: "3642e0e6-788e-4161-92dd-6c52ea823da9", UserID: user.ID}

	// Mock the database & add expectations
	mockDB := mocksqldb.NewMockConnection(mockCtrl)
	mockDB.QEXPECT().GetSuccess(&auth.Session{}, func(sess *auth.Session, stmt string, token string) {
		sess.ID = handlerParams.Token
		sess.UserID = session.UserID
	})
	// delete call
	mockDB.QEXPECT().DeletionSuccess()

	// Mock the response & add expectations
	res := mockrouter.NewMockHTTPResponse(mockCtrl)
	res.EXPECT().NoContent().Return()

	// Mock the request & add expectations
	req := mockrouter.NewMockHTTPRequest(mockCtrl)
	req.EXPECT().Response().Return(res)
	req.EXPECT().Params().Return(handlerParams)
	req.EXPECT().User().Return(user).Times(2)
	req.EXPECT().Session().Return(session)

	// call the handler
	err = sessions.Delete(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.NoError(t, err, "the handler should not have fail")
}

func TestDeleteOtherSessionWrongPassword(t *testing.T) {
	t.Parallel()

	hashr := &bcrypt.Bcrypt{}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	handlerParams := &sessions.DeleteParams{
		Token:           "d1ba3e60-d674-47f0-9694-59fbda0fc659",
		CurrentPassword: "invalid password",
	}

	// Generate a password for the user
	userPassword, err := hashr.Hash("Valid password")
	assert.NoError(t, err)

	// defined the request data
	user := &auth.User{ID: "3e916798-a090-4f22-b1d1-04a63fbed6ef", Password: userPassword}
	session := &auth.Session{ID: "3642e0e6-788e-4161-92dd-6c52ea823da9", UserID: user.ID}

	// Mock the request & add expectations
	req := mockrouter.NewMockHTTPRequest(mockCtrl)
	req.EXPECT().Params().Return(handlerParams)
	req.EXPECT().User().Return(user)
	req.EXPECT().Session().Return(session)

	// call the handler
	err = sessions.Delete(req, &router.Dependencies{DB: nil})

	// Assert everything
	assert.Error(t, err, "the handler should not have fail")
	assert.Equal(t, http.StatusUnauthorized, apierror.Convert(err).HTTPStatus(), "Should have fail with a 401")
}

func TestDeleteSomeonesSession(t *testing.T) {
	t.Parallel()

	hashr := &bcrypt.Bcrypt{}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	handlerParams := &sessions.DeleteParams{
		Token:           "d1ba3e60-d674-47f0-9694-59fbda0fc659",
		CurrentPassword: "valid password",
	}

	// Generate a password for the user
	userPassword, err := hashr.Hash(handlerParams.CurrentPassword)
	assert.NoError(t, err)

	// defined the request data
	user := &auth.User{ID: "3e916798-a090-4f22-b1d1-04a63fbed6ef", Password: userPassword}
	session := &auth.Session{ID: "3642e0e6-788e-4161-92dd-6c52ea823da9", UserID: user.ID}

	// Mock the database & add expectations
	mockDB := mocksqldb.NewMockConnection(mockCtrl)
	mockDB.QEXPECT().GetSuccess(&auth.Session{}, func(sess *auth.Session, stmt string, token string) {
		sess.ID = handlerParams.Token
		sess.UserID = "d15e8b30-69ad-405b-a0f0-0e298b994d89"
	})

	// Mock the request & add expectations
	req := mockrouter.NewMockHTTPRequest(mockCtrl)
	req.EXPECT().Params().Return(handlerParams)
	req.EXPECT().User().Return(user).Times(2)
	req.EXPECT().Session().Return(session)

	// call the handler
	err = sessions.Delete(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should not have fail")
	assert.Equal(t, http.StatusNotFound, apierror.Convert(err).HTTPStatus(), "Should have fail with a 404")
}

func TestDeleteUnexistingSession(t *testing.T) {
	t.Parallel()

	hashr := &bcrypt.Bcrypt{}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	handlerParams := &sessions.DeleteParams{
		Token:           "d1ba3e60-d674-47f0-9694-59fbda0fc659",
		CurrentPassword: "valid password",
	}

	// Generate a password for the user
	userPassword, err := hashr.Hash(handlerParams.CurrentPassword)
	assert.NoError(t, err)

	// defined the request data
	user := &auth.User{ID: "3e916798-a090-4f22-b1d1-04a63fbed6ef", Password: userPassword}
	session := &auth.Session{ID: "3642e0e6-788e-4161-92dd-6c52ea823da9", UserID: user.ID}

	// Mock the database & add expectations
	mockDB := mocksqldb.NewMockConnection(mockCtrl)
	mockDB.QEXPECT().GetNotFound(&auth.Session{})

	// Mock the request & add expectations
	req := mockrouter.NewMockHTTPRequest(mockCtrl)
	req.EXPECT().Params().Return(handlerParams)
	req.EXPECT().User().Return(user)
	req.EXPECT().Session().Return(session)

	// call the handler
	err = sessions.Delete(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should not have fail")
	assert.Equal(t, http.StatusNotFound, apierror.Convert(err).HTTPStatus(), "Should have fail with a 404")
}

func TestDeleteNoDBConOnDelete(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	user := &auth.User{ID: "3e916798-a090-4f22-b1d1-04a63fbed6ef"}
	session := &auth.Session{ID: "3642e0e6-788e-4161-92dd-6c52ea823da9", UserID: user.ID}

	handlerParams := &sessions.DeleteParams{
		Token: session.ID,
	}

	// Mock the database & add expectations
	mockDB := mocksqldb.NewMockConnection(mockCtrl)
	mockDB.QEXPECT().GetSuccess(&auth.Session{}, func(sess *auth.Session, stmt string, token string) {
		sess.ID = session.ID
		sess.UserID = session.UserID
	})
	// delete call
	mockDB.QEXPECT().DeletionError(errors.New("could not delete"))

	// Mock the request & add expectations
	req := mockrouter.NewMockHTTPRequest(mockCtrl)
	req.EXPECT().Params().Return(handlerParams)
	req.EXPECT().User().Return(user)
	req.EXPECT().Session().Return(session)

	// call the handler
	err := sessions.Delete(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusInternalServerError, httpErr.HTTPStatus())
}

func TestDeleteNoDBConOnGet(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	user := &auth.User{ID: "3e916798-a090-4f22-b1d1-04a63fbed6ef"}
	session := &auth.Session{ID: "3642e0e6-788e-4161-92dd-6c52ea823da9", UserID: user.ID}

	handlerParams := &sessions.DeleteParams{
		Token: session.ID,
	}

	// Mock the database & add expectations
	mockDB := mocksqldb.NewMockConnection(mockCtrl)
	mockDB.QEXPECT().GetError(&auth.Session{}, errors.New("could not get"))

	// Mock the request & add expectations
	req := mockrouter.NewMockHTTPRequest(mockCtrl)
	req.EXPECT().Params().Return(handlerParams)
	req.EXPECT().Session().Return(session)

	// call the handler
	err := sessions.Delete(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusInternalServerError, httpErr.HTTPStatus())
}
