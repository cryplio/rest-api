package handlers_test

import (
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/Nivl/go-rest-tools/security/auth/testauth"

	"github.com/cryplio/rest-api/src/modules/portfolios"
	"github.com/cryplio/rest-api/src/modules/portfolios/http/handlers"
	"github.com/cryplio/rest-api/src/modules/portfolios/testportfolios"
	"github.com/dchest/uniuri"
	gomock "github.com/golang/mock/gomock"

	"github.com/Nivl/go-params"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard/testguard"
	"github.com/Nivl/go-rest-tools/router/mockrouter"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/Nivl/go-sqldb/implementations/mocksqldb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			Description: "Should fail on name too long",
			MsgMatch:    params.ErrMsgMaxLen,
			FieldName:   "name",
			Sources: map[string]url.Values{
				"url": url.Values{
					"id": []string{"48d0c8b8-d7a3-4855-9d90-29a06ef474b0"},
				},
				"form": url.Values{
					"name": []string{uniuri.NewLen(51)},
				},
			},
		},
	}

	g := handlers.Endpoints[handlers.EndpointUpdate].Guard
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
		{
			"Should work with a 50 char long name",
			map[string]url.Values{
				"url": url.Values{
					"id": []string{"48d0c8b8-d7a3-4855-9d90-29a06ef474b0"},
				},
				"form": url.Values{
					"name": []string{uniuri.NewLen(50)},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			endpts := handlers.Endpoints[handlers.EndpointUpdate]
			data, err := endpts.Guard.ParseParams(tc.sources, nil)
			assert.NoError(t, err)

			if data != nil {
				p := data.(*handlers.UpdateParams)
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

	g := handlers.Endpoints[handlers.EndpointUpdate].Guard
	testguard.AccessTest(t, g, testCases)
}

func TestUpdateHappyPath(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	portfolio := testportfolios.NewPortfolio()
	handlerParams := &handlers.UpdateParams{
		ID:   portfolio.ID,
		Name: "my renamed portfolio",
	}

	// Mock the database & add expectations
	mockDB := mocksqldb.NewMockConnection(mockCtrl)
	mockDB.QEXPECT().GetSuccess(&portfolios.Portfolio{}, func(p *portfolios.Portfolio, stmt, id string) {
		*p = *portfolio
	})
	mockDB.QEXPECT().UpdateSuccess(&portfolios.Portfolio{})

	// Mock the response & add expectati ons
	res := mockrouter.NewMockHTTPResponse(mockCtrl)
	res.EXPECT().OkSuccess(&portfolios.Payload{}, func(p *portfolios.Payload) {
		assert.Equal(t, p.Name, handlerParams.Name, "the name should have changed")
	})

	// Mock the request & add expectations
	req := mockrouter.NewMockHTTPRequest(mockCtrl)
	req.EXPECT().Response().Return(res)
	req.EXPECT().Params().Return(handlerParams)
	req.EXPECT().User().Return(portfolio.User)

	// call the handler
	err := handlers.Update(req, &router.Dependencies{DB: mockDB})
	assert.NoError(t, err, "the handler should not have fail")
}

func TestUpdatePortfolioNotFound(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	handlerParams := &handlers.UpdateParams{
		ID:   "9cef9a62-f320-4efe-b52c-6038c1e4668c",
		Name: "my renamed portfolio",
	}

	// Mock the database & add expectations
	mockDB := mocksqldb.NewMockConnection(mockCtrl)
	mockDB.QEXPECT().GetNotFound(&portfolios.Portfolio{})

	// Mock the request & add expectations
	req := mockrouter.NewMockHTTPRequest(mockCtrl)
	req.EXPECT().Params().Return(handlerParams)

	// call the handler
	err := handlers.Update(req, &router.Dependencies{DB: mockDB})
	require.Error(t, err, "the handler should have fail")
	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusNotFound, httpErr.HTTPStatus())
}

func TestUpdateWrongUser(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	requester := testauth.NewUser()
	portfolio := testportfolios.NewPortfolio()
	handlerParams := &handlers.UpdateParams{
		ID:   portfolio.ID,
		Name: "my renamed portfolio",
	}

	// Mock the database & add expectations
	mockDB := mocksqldb.NewMockConnection(mockCtrl)
	mockDB.QEXPECT().GetSuccess(&portfolios.Portfolio{}, func(p *portfolios.Portfolio, stmt, id string) {
		*p = *portfolio
	})

	// Mock the request & add expectations
	req := mockrouter.NewMockHTTPRequest(mockCtrl)
	req.EXPECT().Params().Return(handlerParams)
	req.EXPECT().User().Return(requester)

	// call the handler
	err := handlers.Update(req, &router.Dependencies{DB: mockDB})
	require.Error(t, err, "the handler should have fail")
	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusNotFound, httpErr.HTTPStatus())
}

func TestUpdateUpdateError(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	portfolio := testportfolios.NewPortfolio()
	handlerParams := &handlers.UpdateParams{
		ID:   portfolio.ID,
		Name: "my renamed portfolio",
	}

	// Mock the database & add expectations
	mockDB := mocksqldb.NewMockConnection(mockCtrl)
	mockDB.QEXPECT().GetSuccess(&portfolios.Portfolio{}, func(p *portfolios.Portfolio, stmt, id string) {
		*p = *portfolio
	})
	mockDB.QEXPECT().UpdateError(&portfolios.Portfolio{}, errors.New("could not update portfolio"))

	// Mock the request & add expectations
	req := mockrouter.NewMockHTTPRequest(mockCtrl)
	req.EXPECT().Params().Return(handlerParams)
	req.EXPECT().User().Return(portfolio.User)

	// call the handler
	err := handlers.Update(req, &router.Dependencies{DB: mockDB})
	require.Error(t, err, "the handler should have fail")
	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusInternalServerError, httpErr.HTTPStatus())
}
