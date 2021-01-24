package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"sort"
	"time"

	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"

	"github.com/fabritsius/envar"
	"github.com/fabritsius/investor/aggregator/messages"
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
	log.Printf("receive message body from client: %s", in.User[2:10])
	userPortfolio := getPortfolioStats(in.User)
	totalPortfolioValue := userPortfolio.Totals.Invested + userPortfolio.Totals.Yield
	return &messages.PortfolioReply{Data: fmt.Sprintf("%.2f", totalPortfolioValue)}, nil
}

func getPortfolioStats(token string) *PortfolioStats {
	client := sdk.NewRestClient(token)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	positions, err := client.PositionsPortfolio(ctx, sdk.DefaultAccount)
	if err != nil {
		log.Fatalln(err)
	}

	return getStatsFromPositions(client, positions)
}

func getStatsFromPositions(client *sdk.RestClient, positions []sdk.PositionBalance) *PortfolioStats {
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

// GetTotalValue calculates and return total value of all shares for a stock
func (s StockData) GetTotalValue() float64 {
	return s.Balance*s.Price + s.Yield
}

// ByTotalValue implements sort.Interface for []StockData based on getTotalValue()
type ByTotalValue []StockData

func (v ByTotalValue) Len() int      { return len(v) }
func (v ByTotalValue) Swap(i, j int) { v[i], v[j] = v[j], v[i] }
func (v ByTotalValue) Less(i, j int) bool {
	return v[i].GetTotalValue() < v[j].GetTotalValue()
}

// DollarConversionMap has conversion scalars for currencies
type DollarConversionMap = map[sdk.Currency]float64

// String method prints PortfolioStats in a table-like form
func (s PortfolioStats) String() string {
	sort.Sort(ByTotalValue(s.Stocks))
	stats, _ := json.MarshalIndent(s, "", " ")
	return fmt.Sprintf(string(stats))
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
