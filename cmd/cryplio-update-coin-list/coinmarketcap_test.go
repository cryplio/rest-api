package updatecoinlist_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	updatecoinlist "github.com/cryplio/rest-api/cmd/cryplio-update-coin-list"
	"github.com/cryplio/rest-api/src/interfaces/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCoinmarketcapTickers(t *testing.T) {
	testCases := []struct {
		description           string
		statusCode            int
		currentPage           int
		resultsPerpage        int
		expectsErrors         bool
		expectedAmountTickers int
	}{
		{
			description:           "request a full page of tickers",
			statusCode:            200,
			currentPage:           1,
			resultsPerpage:        5,
			expectsErrors:         false,
			expectedAmountTickers: 5,
		},
		{
			description:           "request an incomplete page of tickers",
			statusCode:            200,
			currentPage:           1,
			resultsPerpage:        5,
			expectsErrors:         false,
			expectedAmountTickers: 2,
		},
		{
			description:           "request an empty page of tickers",
			statusCode:            404,
			currentPage:           1,
			resultsPerpage:        5,
			expectsErrors:         false,
			expectedAmountTickers: 0,
		},
		{
			description:           "bad request",
			statusCode:            400,
			currentPage:           1,
			resultsPerpage:        5,
			expectsErrors:         true,
			expectedAmountTickers: 0,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			var pageContent []byte
			if tc.expectedAmountTickers > 0 {
				// We get the right tickers from the global list
				from := tc.currentPage - 1
				to := tc.expectedAmountTickers

				var err error
				pageContent, err = json.Marshal(allTickers[from:to])
				if err != nil {
					t.Fatal(err)
				}
			}

			if tc.expectsErrors {
				var err error
				pageContent, err = json.Marshal(&updatecoinlist.CoinmarketcapAPIError{Error: "request error"})
				if err != nil {
					t.Fatal(err)
				}
			}

			// We mock the http request
			mockHTTP := mocks.NewMockHTTPGetter(mockCtrl)
			mockHTTP.
				EXPECT().Get(fmt.Sprintf("https://api.coinmarketcap.com/v1/ticker/?start=%d&limit=%d", (tc.currentPage-1)*tc.resultsPerpage, tc.resultsPerpage)).
				Return(&http.Response{
					StatusCode: tc.statusCode,
					Body:       ioutil.NopCloser(bytes.NewReader(pageContent)),
				}, nil).
				Times(1)

			// Run and assert
			tickers, err := updatecoinlist.GetCoinmarketcapTickers(mockHTTP, tc.currentPage, tc.resultsPerpage)
			if tc.expectsErrors {
				require.Error(t, err, "The http request should have failed")
			} else {
				require.NoError(t, err, "The http request should not have failed")
			}

			assert.Equal(t, tc.expectedAmountTickers, len(tickers))
		})
	}
}

// func TestGetCoinmarketcapTickers(t *testing.T) {
// 	mockCtrl := gomock.NewController(t)
// 	defer mockCtrl.Finish()

// mockHTTP := mocks.NewMockHTTPGetter(mockCtrl)
// resultPerPage := 5

// 	//
// 	// Sets expectations
// 	//

// 	// Page 1 is full: 0 to resultPerPage
// 	pageContent, err := json.Marshal(allTickers[0:resultPerPage])
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	mockHTTP.
// 		EXPECT().Get(fmt.Sprintf("https://api.coinmarketcap.com/v1/ticker/?start=0&limit=%d", resultPerPage)).
// 		Return(&http.Response{
// 			StatusCode: 200,
// 			Body:       ioutil.NopCloser(bytes.NewReader(pageContent)),
// 		}, nil).
// 		Times(1)

// 	// Page 2 only has 2 entries
// 	pageContent, err = json.Marshal(allTickers[resultPerPage:2])
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	mockHTTP.
// 		EXPECT().Get(fmt.Sprintf("https://api.coinmarketcap.com/v1/ticker/?start=%d&limit=%d", resultPerPage, 2)).
// 		Return(&http.Response{
// 			StatusCode: 200,
// 			Body:       ioutil.NopCloser(bytes.NewReader(pageContent)),
// 		}, nil).
// 		Times(1)

// 	// Page 3 is empty and therefore returns a 404 error (that's how the API works)
// 	pageContent, err = json.Marshal(&updatecoinlist.CoinmarketcapAPIError{Error: "id not found"})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	mockHTTP.
// 		EXPECT().Get(fmt.Sprintf("https://api.coinmarketcap.com/v1/ticker/?start=%d&limit=%d", resultPerPage, 1)).
// 		Return(&http.Response{
// 			StatusCode: 404,
// 			Body:       ioutil.NopCloser(bytes.NewReader(pageContent)),
// 		}, nil).
// 		Times(1)

// 	tickers, err := updatecoinlist.GetCoinmarketcapTickers(mockHTTP, 1, resultPerPage)
// 	require.NoError(t, err, "The first page should not have failed")
// 	assert.Equal()
// }
