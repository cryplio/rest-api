package updatecoinlist_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"testing"

	"github.com/cryplio/rest-api/src/modules/tickers"

	"github.com/stretchr/testify/assert"

	"github.com/Nivl/go-sqldb/implementations/mocksqldb"
	matcher "github.com/Nivl/gomock-type-matcher"
	updatecoinlist "github.com/cryplio/rest-api/cmd/cryplio-update-coin-list"
	"github.com/cryplio/rest-api/src/interfaces/mocks"
	"github.com/golang/mock/gomock"
)

var allTickers []*updatecoinlist.CoinmarketcapTicker

func init() {
	f, err := os.Open(path.Join("testdata", "100_tickers.json"))
	if err != nil {
		panic(err)
	}
	if err := json.NewDecoder(f).Decode(&allTickers); err != nil {
		panic(err)
	}
}

func TestGetExistingTickersHappyPath(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	existingSymbols := []string{"BTC", "LTC", "ETH"}

	mockDB := mocksqldb.NewMockConnection(mockCtrl)
	mockDB.QEXPECT().Select(matcher.Interface(&[]string{}), matcher.String()).
		Do(func(symbols *[]string, stmt string) {
			*symbols = existingSymbols
		}).Return(nil).Times(1)

	symbols, err := updatecoinlist.GetExistingTickers(mockDB)
	assert.NoError(t, err, "GetExistingTickers should not have returned an error")
	assert.Equal(t, len(existingSymbols), len(symbols), "GetExistingTickers did not returned the expected amount of tickers")
}

func TestUpdateListHappyPath(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Mock the database. We assume the database is empty
	existingSymbols := []string{}
	mockDB := mocksqldb.NewMockConnection(mockCtrl)
	mockDB.QEXPECT().Select(matcher.Interface(&[]string{}), matcher.String()).
		Do(func(symbols *[]string, stmt string) {
			*symbols = existingSymbols
		}).Return(nil).Times(1)

	// We assume coinmarketcap only has 7 tickers
	mockDB.QEXPECT().InsertSuccess(&tickers.Ticker{}).Times(7)

	// Mock CoinMarketCap. We assume there're only 2 pages:
	// page 1: 5 tickers, page 2: 2 tickers

	mockHTTP := mocks.NewMockHTTPGetter(mockCtrl)
	// Page 1: 5 tickers
	pageContent, err := json.Marshal(allTickers[0:5])
	if err != nil {
		t.Fatal(err)
	}
	mockHTTP.
		EXPECT().Get("https://api.coinmarketcap.com/v1/ticker/?start=0&limit=5").
		Return(&http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader(pageContent)),
		}, nil).
		Times(1)

	// Page 2: 2 tickers
	pageContent, err = json.Marshal(allTickers[5:7])
	if err != nil {
		t.Fatal(err)
	}
	mockHTTP.
		EXPECT().Get("https://api.coinmarketcap.com/v1/ticker/?start=5&limit=5").
		Return(&http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader(pageContent)),
		}, nil).
		Times(1)

	// Page 3 is empty and therefore returns a 404 error (that's how the API works)
	pageContent, err = json.Marshal(&updatecoinlist.CoinmarketcapAPIError{Error: "id not found"})
	if err != nil {
		t.Fatal(err)
	}
	mockHTTP.
		EXPECT().Get("https://api.coinmarketcap.com/v1/ticker/?start=10&limit=5").
		Return(&http.Response{
			StatusCode: 404,
			Body:       ioutil.NopCloser(bytes.NewReader(pageContent)),
		}, nil).
		Times(1)

	updatecoinlist.UpdateList(mockDB, mockHTTP, 5)
}
