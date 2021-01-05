package main

import (
	"context"
	"errors"
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

	portfolioStatsByCurrency := getTotalPositionsValue(positions)

	for _, portfolioStats := range portfolioStatsByCurrency {
		fmt.Println(portfolioStats)
	}
}

func getTotalPositionsValue(positions []sdk.PositionBalance) map[sdk.Currency]*PortfolioStats {
	portfolioStats := make(map[sdk.Currency]*PortfolioStats)
	for _, pos := range positions {
		currency := pos.AveragePositionPrice.Currency

		positionStats := &PortfolioStats{
			Date:     time.Now(),
			Invested: pos.AveragePositionPrice.Value * pos.Balance,
			Yield:    pos.ExpectedYield.Value,
			Stocks:   make(map[string]float64),
			Currency: currency,
		}
		positionStats.Stocks[pos.FIGI] = pos.Balance

		if prevStats, ok := portfolioStats[currency]; ok {
			if err := (*prevStats).add(positionStats); err != nil {
				log.Fatalln(err)
			}
		} else {
			portfolioStats[currency] = positionStats
		}
	}
	return portfolioStats
}

// PortfolioStats contains main portfolio stats for the moment
type PortfolioStats struct {
	Date     time.Time
	Invested float64
	Yield    float64
	Stocks   map[string]float64
	Currency sdk.Currency
}

func (s *PortfolioStats) add(new *PortfolioStats) error {
	if s.Currency != new.Currency {
		return errors.New("Can't add. Currencies do no match")
	}

	s.Invested += new.Invested
	s.Yield += new.Yield
	for bk, bv := range new.Stocks {
		s.Stocks[bk] = bv
	}

	return nil
}

// String method prints PortfolioStats in a table-like form
func (s PortfolioStats) String() string {
	return fmt.Sprintf("%3s: invested: %10.2f | yield: %10.2f | total: %10.2f", s.Currency, s.Invested, s.Yield, s.Invested+s.Yield)
}

// func getDollarPrice(client *sdk.RestClient) float64 {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	dollarFIGI := "BBG0013HGFT4"
// 	candles, err := client.Candles(ctx, time.Now().AddDate(0, 0, -7), time.Now(), sdk.CandleInterval1Day, dollarFIGI)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	latestCandle := candles[len(candles)-1]
// 	dollarPrice := latestCandle.ClosePrice
// 	return dollarPrice
// }
