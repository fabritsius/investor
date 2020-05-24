package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
)

var token = flag.String("token", "", "your token")

func main() {
	rand.Seed(time.Now().UnixNano())
	flag.Parse()

	client := sdk.NewRestClient(*token)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	positions, err := client.PositionsPortfolio(ctx, sdk.DefaultAccount)
	if err != nil {
		log.Fatalln(err)
	}

	var portfolioValue PortfolioValue

	portfolioValue.ByCurr = getTotalPositionsValue(positions)
	dollarPrice := getDollarPrice(client)

	portfolioValue.CalcTotals(dollarPrice)
	fmt.Print(portfolioValue)
}

// PortfolioValue contains potfolio value by currency and converted totals
type PortfolioValue struct {
	ByCurr TotalAvgValueByCurrency
	Totals ConvertedTotalAvgValue
}

// CalcTotals calculates and fills in the Totals field
func (v *PortfolioValue) CalcTotals(dollarPrice float64) {
	v.Totals.USD = v.ByCurr.USD + (v.ByCurr.RUB / dollarPrice)
	v.Totals.RUB = v.ByCurr.RUB + (v.ByCurr.USD * dollarPrice)
}

func (v PortfolioValue) String() string {
	return fmt.Sprintf(`by currency:
	USD: %10.2f
	RUB: %10.2f
converted totals:
	USD: %10.2f
	RUB: %10.2f
`, v.ByCurr.USD, v.ByCurr.RUB, v.Totals.USD, v.Totals.RUB)
}

func getTotalPositionsValue(positions []sdk.PositionBalance) TotalAvgValueByCurrency {
	var positionTotals TotalAvgValueByCurrency
	for _, pos := range positions {
		posTotal := pos.Balance * pos.AveragePositionPrice.Value
		currency := pos.AveragePositionPrice.Currency
		switch currency {
		case sdk.RUB:
			positionTotals.RUB += posTotal
		case sdk.USD:
			positionTotals.USD += posTotal
		}
	}
	return positionTotals
}

func getDollarPrice(client *sdk.RestClient) float64 {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dollarFIGI := "BBG0013HGFT4"
	candles, err := client.Candles(ctx, time.Now().AddDate(0, 0, -7), time.Now(), sdk.CandleInterval1Day, dollarFIGI)
	if err != nil {
		log.Fatalln(err)
	}
	latestCandle := candles[len(candles)-1]
	dollarPrice := latestCandle.ClosePrice
	return dollarPrice
}

// TotalAvgValueByCurrency contains portfolio value for USD and RUB
type TotalAvgValueByCurrency usdrubs

// ConvertedTotalAvgValue contains total value converted for USD and RUB
type ConvertedTotalAvgValue usdrubs

type usdrubs struct {
	USD float64
	RUB float64
}
