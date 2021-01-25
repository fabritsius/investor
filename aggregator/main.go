package main

import (
	"fmt"
	"log"

	"github.com/fabritsius/envar"
	"github.com/fabritsius/investor/aggregator/messages"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type config struct {
	Port         string `env:"TINKOFF_PORT"`
	TinkoffToken string `env:"TINKOFF_API_TOKEN"`
}

func main() {
	cfg := config{}
	if err := envar.Fill(&cfg); err != nil {
		log.Fatalln(err)
	}

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(fmt.Sprintf(":%s", cfg.Port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	client := messages.NewPortfolioClient(conn)

	response, err := client.GetPortfolio(context.Background(), &messages.PortfolioRequest{
		Options: map[string]string{"token": cfg.TinkoffToken},
	})
	if err != nil {
		log.Fatalf("error when calling GetPortfolio: %s", err)
	}
	log.Printf("Total Portfolio Value: $%s", response.Data)

}
