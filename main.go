package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	"github.com/fabritsius/envar"
)

type config struct {
	TinkoffToken string `env:"TINKOFF_API_TOKEN"`
}

func main() {
	cfg := config{}
	if err := envar.Fill(&cfg); err != nil {
		log.Fatalln(err)
	}

	rand.Seed(time.Now().UnixNano())
	flag.Parse()

	client := sdk.NewRestClient(cfg.TinkoffToken)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	positions, err := client.PositionsPortfolio(ctx, sdk.DefaultAccount)
	if err != nil {
		log.Fatalln(err)
	}

	portfolioStats := getPortfolioStats(client, positions)
	fmt.Println(*portfolioStats)
}

func getPortfolioStats(client *sdk.RestClient, positions []sdk.PositionBalance) *PortfolioStats {
	portfolioStats := &PortfolioStats{
		Currency: "USD",
		Date:     time.Now(),
		Stocks:   []StockData{},
		Totals:   PortfolioTotals{},
	}
	dollarConversionMap := make(DollarConversionMap)
	dollarRubPrice, err := getDollarPriceInRubles(client)
	if err == nil {
		dollarConversionMap["RUB"] = dollarRubPrice
	}

	for _, pos := range positions {
		currency := pos.AveragePositionPrice.Currency

		price, err := convertToDollar(currency, pos.AveragePositionPrice.Value, dollarConversionMap)
		if err != nil {
			log.Fatalln(err)
		}
		yield, err := convertToDollar(currency, pos.ExpectedYield.Value, dollarConversionMap)
		if err != nil {
			log.Fatalln(err)
		}

		stockIDs := map[string]string{
			"FIGI": pos.FIGI,
			"ISIN": pos.ISIN,
		}
		stockData := StockData{
			IDs:     stockIDs,
			Name:    pos.Name,
			Balance: pos.Balance,
			Price:   price,
			Yield:   yield,
		}

		portfolioStats.Totals.Invested += stockData.Balance * stockData.Price
		portfolioStats.Totals.Yield += stockData.Yield
		portfolioStats.Stocks = append(portfolioStats.Stocks, stockData)
	}
	return portfolioStats
}

func convertToDollar(currency sdk.Currency, value float64, dollarConversionMap DollarConversionMap) (float64, error) {
	switch currency {
	case "USD":
		return value, nil
	case "RUB":
		dollarPrice := dollarConversionMap["RUB"]
		return value / dollarPrice, nil
	default:
		return 0, fmt.Errorf("cannot convert: unrecognised currency %q", currency)
	}
}

// PortfolioStats contains main portfolio stats for the moment
type PortfolioStats struct {
	Currency sdk.Currency
	Date     time.Time
	Stocks   []StockData
	Totals   PortfolioTotals
}

// PortfolioTotals contains calculated totals for portfolio
type PortfolioTotals struct {
	Invested float64
	Yield    float64
}

// StockData contains main data about the stock
type StockData struct {
	IDs     map[string]string
	Name    string
	Balance float64
	Price   float64
	Yield   float64
}

// DollarConversionMap has conversion scalars for currencies
type DollarConversionMap = map[sdk.Currency]float64

// String method prints PortfolioStats in a table-like form
func (s PortfolioStats) String() string {
	stocks, _ := json.MarshalIndent(s.Stocks, "", " ")
	return fmt.Sprintf("%s\n\nTotals in %3s >> invested: %10.2f | yield: %10.2f | grand total: %10.2f", string(stocks), s.Currency, s.Totals.Invested, s.Totals.Yield, s.Totals.Invested+s.Totals.Yield)
}

func getDollarPriceInRubles(client *sdk.RestClient) (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dollarFIGI := "BBG0013HGFT4"
	candles, err := client.Candles(ctx, time.Now().AddDate(0, 0, -7), time.Now(), sdk.CandleInterval1Day, dollarFIGI)
	if err != nil {
		return 0, err
	}
	latestCandle := candles[len(candles)-1]
	dollarPrice := latestCandle.ClosePrice
	return dollarPrice, nil
}
