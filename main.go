package main

import (
	"context"
	"flag"
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Get non-currency assets")
	positions, err := client.PositionsPortfolio(ctx, sdk.DefaultAccount)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("%+v\n\n", positions)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Get currency assets")
	positionCurrencies, err := client.CurrenciesPortfolio(ctx, sdk.DefaultAccount)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("%+v\n", positionCurrencies)
}
