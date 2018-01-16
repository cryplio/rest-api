package updatecoinlist

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cryplio/rest-api/src/interfaces"
)

// CoinmarketcapAPIError represents an error thrown by the coinmarketcap API
type CoinmarketcapAPIError struct {
	Error string `json:"error"`
}

// CoinmarketcapTicker represents a ticker on coinmarketcap
type CoinmarketcapTicker struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	Symbol           string  `json:"symbol"`
	PriceUSD         string  `json:"price_usd"`
	Volume           *string `json:"24h_volume_usd"`
	MarketCap        *string `json:"market_cap_usd"`
	AvailableSupply  *string `json:"available_supply"`
	MaxSupply        *string `json:"max_supply"`
	PercentChange1h  *string `json:"percent_change_1h"`
	PercentChange24h *string `json:"percent_change_24h"`
	PercentChange7d  *string `json:"percent_change_7d"`
}

// GetCoinmarketcapTickers returns the tickers of CoinmarketCap
func GetCoinmarketcapTickers(getter interfaces.HTTPGetter, page, itemsPerPage int) ([]*CoinmarketcapTicker, error) {
	endpoint := fmt.Sprintf("https://api.coinmarketcap.com/v1/ticker/?start=%d&limit=%d", (page-1)*itemsPerPage, itemsPerPage)
	resp, err := getter.Get(endpoint)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case 200:
		var pld []*CoinmarketcapTicker
		if err := json.NewDecoder(resp.Body).Decode(&pld); err != nil {
			return nil, err
		}
		return pld, nil
	case 404:
		// nothing else to get
		return nil, nil
	default:
		var pld *CoinmarketcapAPIError
		if err := json.NewDecoder(resp.Body).Decode(&pld); err != nil {
			return nil, err
		}
		return nil, errors.New(pld.Error)
	}
}
