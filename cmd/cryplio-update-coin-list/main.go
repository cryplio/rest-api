package updatecoinlist

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	sqldb "github.com/Nivl/go-sqldb"
	"github.com/cryplio/rest-api/src/interfaces"
	"github.com/cryplio/rest-api/src/modules/api"
	"github.com/cryplio/rest-api/src/modules/tickers"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "update-coin-list"
	app.Usage = "Update the list of coins"
	app.Version = "1.0.0"

	app.Action = func(c *cli.Context) error {
		httpClient := &http.Client{
			Timeout: time.Duration(10 * time.Second),
		}
		_, deps, err := api.DefaultSetup()
		if err != nil {
			fmt.Printf("[error] could not connect to the database: %s\n", err.Error())
			return nil
		}
		UpdateList(deps.DB(), httpClient, 1000)
		return nil
	}
	app.Run(os.Args)
}

// UpdateList fetch the remote tickers and adds the missing one in the database
func UpdateList(con sqldb.Connection, httpClient interfaces.HTTPGetter, pageSize int) {
	existingTickers, err := GetExistingTickers(con)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// We paginate over all coinMarketCap tickers
	for currentPage := 1; ; currentPage++ {
		cmcTickers, err := GetCoinmarketcapTickers(httpClient, currentPage, pageSize)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// if there're no cmcTickers, then we're done
		if len(cmcTickers) == 0 {
			break
		}

		// we browse all the tickers on the current page
		for _, ticker := range cmcTickers {
			_, alreadyInDB := existingTickers[ticker.Symbol]
			if !alreadyInDB {
				t := &tickers.Ticker{
					ID:              ticker.Symbol,
					Name:            ticker.Name,
					CoinMarketCapID: ticker.ID,
				}

				t.PriceUSD, err = strconv.ParseFloat(ticker.PriceUSD, 64)
				if err != nil {
					fmt.Printf("[error] Failed parsing the USD price for %s: %s\n", ticker.Symbol, err.Error())
					continue
				}

				if ticker.MarketCap != nil {
					value, err := extractInt64(*ticker.MarketCap)
					if err != nil {
						fmt.Printf("[error] Failed parsing the marketcap for %s: %s\n", ticker.Symbol, err.Error())
						continue
					}
					t.Marketcap = &value
				}
				if ticker.Volume != nil {
					value, err := extractInt64(*ticker.Volume)
					if err != nil {
						fmt.Printf("[error] Failed parsing the 24h volume for %s: %s\n", ticker.Symbol, err.Error())
						continue
					}
					t.Volume24h = &value
				}
				if ticker.MaxSupply != nil {
					value, err := extractInt64(*ticker.MaxSupply)
					if err != nil {
						fmt.Printf("[error] Failed parsing the max supply for %s: %s\n", ticker.Symbol, err.Error())
						continue
					}
					t.MaxSupply = &value
				}
				if ticker.AvailableSupply != nil {
					value, err := extractInt64(*ticker.AvailableSupply)
					if err != nil {
						fmt.Printf("[error] Failed parsing the current supply for %s: %s\n", ticker.Symbol, err.Error())
						continue
					}
					t.CurrentSupply = &value
				}
				if ticker.PercentChange1h != nil {
					value, err := strconv.ParseFloat(*ticker.PercentChange1h, 64)
					if err != nil {
						fmt.Printf("[error] Failed parsing the 1h percentage change for %s: %s\n", ticker.Symbol, err.Error())
						continue
					}
					t.PercentChange1h = &value
				}
				if ticker.PercentChange24h != nil {
					value, err := strconv.ParseFloat(*ticker.PercentChange24h, 64)
					if err != nil {
						fmt.Printf("[error] Failed parsing the 24h percentage change for %s: %s\n", ticker.Symbol, err.Error())
						continue
					}
					t.PercentChange24h = &value
				}
				if ticker.PercentChange7d != nil {
					value, err := strconv.ParseFloat(*ticker.PercentChange7d, 64)
					if err != nil {
						fmt.Printf("[error] Failed parsing the 7d percentage change for %s: %s\n", ticker.Symbol, err.Error())
						continue
					}
					t.PercentChange7d = &value
				}

				if err := t.Create(con); err != nil {
					fmt.Println(err.Error())
					return
				}
			}
		}
	}
}

// GetExistingTickers returns a map of symbols that are already in the database
func GetExistingTickers(con sqldb.Connection) (map[string]bool, error) {
	var symbols []string
	stmt := "SELECT symbol FROM tickers"
	err := con.Select(&symbols, stmt)
	if err != nil {
		return nil, err
	}

	mappedSymbols := map[string]bool{}
	for _, symbol := range symbols {
		mappedSymbols[symbol] = true
	}
	return mappedSymbols, nil
}

func extractInt64(originalValue string) (int64, error) {
	value, err := strconv.ParseFloat(originalValue, 64)
	if err != nil {
		return 0, err
	}
	return int64(value), nil
}
