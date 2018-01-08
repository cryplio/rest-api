package portfolios_test

import (
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/Nivl/go-rest-tools/security/auth/testauth"
	"github.com/Nivl/go-rest-tools/types/apierror"

	"github.com/Nivl/go-params"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard/testguard"
	"github.com/Nivl/go-rest-tools/router/mockrouter"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-sqldb/implementations/mocksqldb"
	"github.com/cryplio/rest-api/src/modules/portfolios"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDeleteInvalidParams(t *testing.T) {
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
			},
		},
	}

	g := portfolios.Endpoints[portfolios.EndpointDelete].Guard
	testguard.InvalidParams(t, g, testCases)
}

func TestDeleteValidParams(t *testing.T) {
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

			endpts := portfolios.Endpoints[portfolios.EndpointDelete]
			data, err := endpts.Guard.ParseParams(tc.sources, nil)
			assert.NoError(t, err)

			if data != nil {
				p := data.(*portfolios.DeleteParams)
				assert.Equal(t, tc.sources["url"].Get("id"), p.ID)
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

	g := portfolios.Endpoints[portfolios.EndpointDelete].Guard
	testguard.AccessTest(t, g, testCases)
}

func TestDeleteHappyPath(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	user := testauth.NewUser()
	handlerParams := &portfolios.DeleteParams{
		ID: "48d0c8b8-d7a3-4855-9d90-29a06ef474b0",
	}

	// Mock the database & add expectations
	mockDB := mocksqldb.NewMockConnection(mockCtrl)
	mockDB.QEXPECT().GetSuccess(&portfolios.Portfolio{}, func(p *portfolios.Portfolio, stmt, id string) {
		p.ID = handlerParams.ID
		p.UserID = user.ID
	})
	mockDB.QEXPECT().DeletionSuccess()

	// Mock the response & add expectations
	res := mockrouter.NewMockHTTPResponse(mockCtrl)
	res.EXPECT().NoContent().Return()

	// Mock the request & add expectations
	req := mockrouter.NewMockHTTPRequest(mockCtrl)
	req.EXPECT().Response().Return(res)
	req.EXPECT().Params().Return(handlerParams)
	req.EXPECT().User().Return(user)

	// call the handler
	err := portfolios.Delete(req, &router.Dependencies{DB: mockDB})
	assert.NoError(t, err, "the handler should not have fail")
}

func TestDeleteDeletionError(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	user := testauth.NewUser()
	handlerParams := &portfolios.DeleteParams{
		ID: "48d0c8b8-d7a3-4855-9d90-29a06ef474b0",
	}

	// Mock the database & add expectations
	mockDB := mocksqldb.NewMockConnection(mockCtrl)
	mockDB.QEXPECT().GetSuccess(&portfolios.Portfolio{}, func(p *portfolios.Portfolio, stmt, id string) {
		p.ID = handlerParams.ID
		p.UserID = user.ID
	})
	mockDB.QEXPECT().DeletionError(errors.New("could not delete"))

	// Mock the request & add expectations
	req := mockrouter.NewMockHTTPRequest(mockCtrl)
	req.EXPECT().Params().Return(handlerParams)
	req.EXPECT().User().Return(user)

	// call the handler
	err := portfolios.Delete(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusInternalServerError, httpErr.HTTPStatus())
}

func TestDeleteGetPortfolioError(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	handlerParams := &portfolios.DeleteParams{
		ID: "48d0c8b8-d7a3-4855-9d90-29a06ef474b0",
	}

	// Mock the database & add expectations
	mockDB := mocksqldb.NewMockConnection(mockCtrl)
	mockDB.QEXPECT().GetError(&portfolios.Portfolio{}, errors.New("could not get portfolio"))

	// Mock the request & add expectations
	req := mockrouter.NewMockHTTPRequest(mockCtrl)
	req.EXPECT().Params().Return(handlerParams)

	// call the handler
	err := portfolios.Delete(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusInternalServerError, httpErr.HTTPStatus())
}

func TestDeleteGetPortfolioNotFound(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	handlerParams := &portfolios.DeleteParams{
		ID: "48d0c8b8-d7a3-4855-9d90-29a06ef474b0",
	}

	// Mock the database & add expectations
	mockDB := mocksqldb.NewMockConnection(mockCtrl)
	mockDB.QEXPECT().GetNotFound(&portfolios.Portfolio{})

	// Mock the request & add expectations
	req := mockrouter.NewMockHTTPRequest(mockCtrl)
	req.EXPECT().Params().Return(handlerParams)

	// call the handler
	err := portfolios.Delete(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusNotFound, httpErr.HTTPStatus())
}

func TestDeleteSomeonesPortfolio(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	user := testauth.NewUser()
	handlerParams := &portfolios.DeleteParams{
		ID: "48d0c8b8-d7a3-4855-9d90-29a06ef474b0",
	}

	// Mock the database & add expectations
	mockDB := mocksqldb.NewMockConnection(mockCtrl)
	mockDB.QEXPECT().GetSuccess(&portfolios.Portfolio{}, func(p *portfolios.Portfolio, stmt, id string) {
		p.ID = handlerParams.ID
		// we use a different user ID as owner of the portfolio
		p.UserID = "9cef9a62-f320-4efe-b52c-6038c1e4668c"
	})

	// Mock the request & add expectations
	req := mockrouter.NewMockHTTPRequest(mockCtrl)
	req.EXPECT().Params().Return(handlerParams)
	req.EXPECT().User().Return(user)

	// call the handler
	err := portfolios.Delete(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusNotFound, httpErr.HTTPStatus())
}
