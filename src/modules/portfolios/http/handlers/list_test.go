package handlers_test

import (
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/Nivl/go-rest-tools/security/auth/testauth"
	matcher "github.com/Nivl/gomock-type-matcher"

	"github.com/Nivl/go-rest-tools/paginator"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard/testguard"
	"github.com/Nivl/go-rest-tools/router/mockrouter"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/Nivl/go-sqldb/implementations/mocksqldb"
	"github.com/cryplio/rest-api/src/modules/portfolios"
	"github.com/cryplio/rest-api/src/modules/portfolios/http/handlers"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListInvalidParams(t *testing.T) {
	t.Parallel()

	testCases := []testguard.InvalidParamsTestCase{
		{
			Description: "Should fail with page = 0",
			MsgMatch:    paginator.ErrMsgNumberBelow1,
			FieldName:   "page",
			Sources: map[string]url.Values{
				"query": url.Values{
					"page": []string{"0"},
				},
			},
		},
		{
			Description: "Should fail with per_page = 0",
			MsgMatch:    paginator.ErrMsgNumberBelow1,
			FieldName:   "per_page",
			Sources: map[string]url.Values{
				"query": url.Values{
					"per_page": []string{"0"},
				},
			},
		},
		{
			Description: "Should fail with per_page > 100",
			MsgMatch:    "cannot be > 100",
			FieldName:   "per_page",
			Sources: map[string]url.Values{
				"query": url.Values{
					"per_page": []string{"101"},
				},
			},
		},
	}

	g := handlers.Endpoints[handlers.EndpointList].Guard
	testguard.InvalidParams(t, g, testCases)
}

func TestListValidParams(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		description string
		sources     map[string]url.Values
	}{
		{
			"Should work with nothing",
			map[string]url.Values{
				"query": url.Values{},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			endpts := handlers.Endpoints[handlers.EndpointList]
			_, err := endpts.Guard.ParseParams(tc.sources, nil)
			assert.NoError(t, err)
		})
	}
}

func TestListAccess(t *testing.T) {
	t.Parallel()

	testCases := []testguard.AccessTestCase{
		{
			Description: "Should fail for anonymous users",
			User:        nil,
			ErrCode:     http.StatusUnauthorized,
		},
		{
			Description: "Should fail for logged users",
			User:        &auth.User{ID: "48d0c8b8-d7a3-4855-9d90-29a06ef474b0"},
			ErrCode:     0,
		},
		{
			Description: "Should work for admin users",
			User:        &auth.User{ID: "48d0c8b8-d7a3-4855-9d90-29a06ef474b0", IsAdmin: true},
			ErrCode:     0,
		},
	}

	g := handlers.Endpoints[handlers.EndpointList].Guard
	testguard.AccessTest(t, g, testCases)
}

func TestListNoBDCon(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	user := testauth.NewUser()

	// Mock the database & add expectations
	mockDB := mocksqldb.NewMockConnection(mockCtrl)
	call := mockDB.QEXPECT().Select(matcher.Interface(&portfolios.Portfolios{}), matcher.String(), user.ID, 0, matcher.Int())
	call.Return(errors.New("no connection"))
	call.Times(1)

	// Mock the request & add expectations
	req := mockrouter.NewMockHTTPRequest(mockCtrl)
	req.EXPECT().Params().Return(&handlers.ListParams{})
	req.EXPECT().User().Return(user)

	// call the handler
	err := handlers.List(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	require.Error(t, err, "the handler should have fail")
	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusInternalServerError, httpErr.HTTPStatus())
}

func TestListInvalidSort(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	user := testauth.NewUser()

	// Mock the request & add expectations
	req := mockrouter.NewMockHTTPRequest(mockCtrl)
	req.EXPECT().Params().Return(&handlers.ListParams{Sort: "not_a_field"})
	req.EXPECT().User().Return(user)

	// call the handler
	err := handlers.List(req, &router.Dependencies{DB: nil})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")

	httpErr := apierror.Convert(err)
	assert.Equal(t, http.StatusBadRequest, httpErr.HTTPStatus())
	assert.Equal(t, "sort", httpErr.Field())
}
