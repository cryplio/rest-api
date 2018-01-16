package tickers

import (
	"github.com/Nivl/go-types/datetime"
)

// Ticker represents a coin/token
//go:generate api-cli generate model Ticker -t tickers
type Ticker struct {
	ID        string             `db:"id"`
	CreatedAt *datetime.DateTime `db:"created_at"`
	UpdatedAt *datetime.DateTime `db:"updated_at"`
	DeletedAt *datetime.DateTime `db:"deleted_at"`

	Name             string   `db:"name"`
	Symbol           string   `db:"symbol"`
	Unit             *string  `db:"unit"`
	Marketcap        *int64   `db:"marketcap"`
	Volume24h        *int64   `db:"volume_24h"`
	MaxSupply        *int64   `db:"max_supply"`
	CurrentSupply    *int64   `db:"current_supply"`
	LogoURL          *string  `db:"logo_url"`
	Website          *string  `db:"website"`
	PriceUSD         float64  `db:"price_usd"`
	PercentChange1h  *float64 `db:"percent_change_1h"`
	PercentChange24h *float64 `db:"percent_change_24h"`
	PercentChange7d  *float64 `db:"percent_change_7d"`
	CoinMarketCapID  string   `db:"coinmarketcap_id"`
}

// Tickers represents a list of Ticker
type Tickers []*Ticker
