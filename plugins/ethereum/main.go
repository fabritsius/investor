package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/fabritsius/envar"
	"github.com/fabritsius/investor/messages"
	"google.golang.org/grpc"
)

type config struct {
	Port         string `env:"ETHEREUM_PORT"`
	EtherscanKey string `env:"ETHERSCAN_API_KEY"`
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

	log.Printf("started ethereum server on port %s", cfg.Port)

	grpcServer := grpc.NewServer()

	messages.RegisterPortfolioServer(grpcServer, &messageServer{
		config: cfg,
	})

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve grpc: %s", err)
	}
}

type messageServer struct {
	config config
}

func (s *messageServer) GetPortfolio(ctx context.Context, in *messages.PortfolioRequest) (*messages.PortfolioReply, error) {
	etherAddress, ok := in.Options["token"]
	if !ok {
		return nil, errors.New("\"token\" is missing")
	}
	log.Printf("receive message body from client: %s", etherAddress[2:10])
	userPortfolio, err := doEtherscan(s.config, etherAddress)
	if err != nil {
		return nil, err
	}

	return &messages.PortfolioReply{Data: userPortfolio}, nil
}

// AccountBalanceResponse contains a balance of an Ethereum Account (returned by Etherscan)
type AccountBalanceResponse struct {
	Status  string
	Message string
	Result  string
}

func doEtherscan(cfg config, address string) (*messages.PortfolioStats, error) {
	requestURL := fmt.Sprintf("https://api.etherscan.io/api?module=account&action=balance&address=%s&tag=latest&apikey=%s", address, cfg.EtherscanKey)
	fmt.Println(requestURL)
	resp, err := http.Get(requestURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	account := &AccountBalanceResponse{}
	if err := json.NewDecoder(resp.Body).Decode(account); err != nil {
		return nil, fmt.Errorf("cannot decode: %s", err)
	}

	etherPrice, err := getEtherPrice(cfg)
	if err != nil {
		return nil, fmt.Errorf("cannot get dollar price of ether: %s", err)
	}

	balance, err := strconv.Atoi(account.Result)
	if err != nil {
		return nil, fmt.Errorf("cannot convert balance: %s", err)
	}
	accountTotalValue := float64(balance) * etherPrice / 1000000000000000000
	portfolioStats := &messages.PortfolioStats{
		Currency: messages.Currency_USD,
		Date:     time.Now().Unix(),
		Stocks:   []*messages.StockData{},
		Totals: &messages.PortfolioTotals{
			Invested: accountTotalValue,
		},
	}

	return portfolioStats, nil
}

// EtherPriceResponse contains ETH prices (returned by Etherscan)
type EtherPriceResponse struct {
	Status  string
	Message string
	Result  EtherPrices
}

// EtherPrices contains a price ETH in dollars and bitcoins
type EtherPrices struct {
	Ethbtc          string
	EthbtcTimestamp string `json:"ethbtc_timestamp"`
	Ethusd          string
	EthusdTimestamp string `json:"ethusd_timestamp"`
}

func getEtherPrice(cfg config) (float64, error) {
	requestURL := fmt.Sprintf("https://api.etherscan.io/api?module=stats&action=ethprice&apikey=%s", cfg.EtherscanKey)
	resp, err := http.Get(requestURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	price := &EtherPriceResponse{}
	if err := json.NewDecoder(resp.Body).Decode(price); err != nil {
		return 0, fmt.Errorf("cannot decode: %s", err)
	}

	if price.Message != "OK" {
		return 0, fmt.Errorf("not OK response from Etherscan: %v", price)
	}

	priceValue, err := strconv.ParseFloat(price.Result.Ethusd, 64)
	if err != nil {
		return 0, fmt.Errorf("cannot convert price in dollars: %s", err)
	}

	return priceValue, nil
}
