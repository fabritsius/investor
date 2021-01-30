package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"sort"
	"time"

	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"

	"github.com/fabritsius/envar"
	"github.com/fabritsius/investor/messages"
	"google.golang.org/grpc"
)

type config struct {
	Port string `env:"PORT"`
}

func main() {
	cfg := config{}
	if err := envar.Fill(&cfg); err != nil {
		log.Fatalln(err)
	}

	rand.Seed(time.Now().UnixNano())
	flag.Parse()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", cfg.Port, err)
	}

	log.Printf("started tinkoff server on port %s", cfg.Port)

	grpcServer := grpc.NewServer()

	messages.RegisterPortfolioServer(grpcServer, &messageServer{})

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve grpc: %s", err)
	}
}

type messageServer struct{}

func (s *messageServer) GetPortfolio(ctx context.Context, in *messages.PortfolioRequest) (*messages.PortfolioReply, error) {
	userToken, ok := in.Options["token"]
	if !ok {
		return nil, errors.New("\"token\" is missing")
	}
	log.Printf("receive message body from client: %s", userToken[2:10])
	userPortfolio := getPortfolioStats(userToken)
	return &messages.PortfolioReply{Data: userPortfolio}, nil
}

func getPortfolioStats(token string) *messages.PortfolioStats {
	client := sdk.NewRestClient(token)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	positions, err := client.PositionsPortfolio(ctx, sdk.DefaultAccount)
	if err != nil {
		log.Fatalln(err)
	}

	return getStatsFromPositions(client, positions)
}

func getStatsFromPositions(client *sdk.RestClient, positions []sdk.PositionBalance) *messages.PortfolioStats {
	portfolioStats := &messages.PortfolioStats{
		Currency: messages.Currency_USD,
		Date:     time.Now().Unix(),
		Stocks:   []*messages.StockData{},
		Totals:   &messages.PortfolioTotals{},
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

		otherFields := map[string]string{
			"FIGI": pos.FIGI,
			"ISIN": pos.ISIN,
			"type": string(pos.InstrumentType),
		}
		stockData := &messages.StockData{
			Name:    pos.Name,
			Balance: pos.Balance,
			Price:   price,
			Yield:   yield,
			Other:   otherFields,
		}

		portfolioStats.Totals.Invested += stockData.Balance * stockData.Price
		portfolioStats.Totals.Yield += stockData.Yield

		portfolioStats.Stocks = append(portfolioStats.Stocks, stockData)
	}

	sort.Sort(sort.Reverse(ByTotalValue(portfolioStats.Stocks)))
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

// getStockTotalValue calculates and return total value of all shares for a stock
func getStockTotalValue(s *messages.StockData) float64 {
	return s.Balance*s.Price + s.Yield
}

// ByTotalValue implements sort.Interface for []StockData based on getStockTotalValue function
type ByTotalValue []*messages.StockData

func (v ByTotalValue) Len() int      { return len(v) }
func (v ByTotalValue) Swap(i, j int) { v[i], v[j] = v[j], v[i] }
func (v ByTotalValue) Less(i, j int) bool {
	return getStockTotalValue(v[i]) < getStockTotalValue(v[j])
}

// DollarConversionMap has conversion scalars for currencies
type DollarConversionMap = map[sdk.Currency]float64

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
