package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/fabritsius/envar"
	"github.com/fabritsius/investor/aggregator/models"
	"github.com/fabritsius/investor/messages"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type config struct {
	Port string `env:"TINKOFF_PORT"`
}

func main() {
	cfg := config{}
	if err := envar.Fill(&cfg); err != nil {
		log.Fatalln(err)
	}

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(fmt.Sprintf(":%s", cfg.Port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect to grpc: %s", err)
	}
	defer conn.Close()

	var db *models.DB
	if db, err = models.Connect("127.0.0.1"); err != nil {
		log.Fatalf("did not connect to the database: %s", err)
	}
	if err = db.Init(); err != nil {
		log.Fatalf("database init error: %s", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	for item := range db.GetAllUserAccounts(ctx) {
		if item.Error != nil {
			log.Printf("account error: %s", err)
			continue
		}

		account := item.Account
		log.Println("Account:", account.UserID, account.AccountType, account.Token)

		var portfolio *messages.PortfolioReply
		switch account.AccountType {
		case "tinkoff":
			if portfolio, err = GetTinkoffPortfolio(conn, account.Token); err != nil {
				log.Printf("error when calling getTinkoffPortfolio: %s", err)
			}
		}

		if portfolio != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			if err := db.UpdateDailyRecord(ctx, &models.PortfolioValue{
				UserID:      account.UserID,
				AccountType: account.AccountType,
				Date:        time.Now().Format("2006-01-02"),
				Invested:    portfolio.Data.Totals.Invested,
				Yield:       portfolio.Data.Totals.Yield,
			}); err != nil {
				log.Printf("failed to update portfolio stats: %s", err)
			}

			stats, _ := json.MarshalIndent(portfolio, "", " ")
			log.Printf("Portfolio: %s", string(stats))
		}
	}

}

// GetTinkoffPortfolio takes grpc connection and Tinkoff API Token
// and returns user's Tinkoff portfolio with stocks and calculated totals
func GetTinkoffPortfolio(conn *grpc.ClientConn, tinkoffToken string) (*messages.PortfolioReply, error) {
	client := messages.NewPortfolioClient(conn)

	response, err := client.GetPortfolio(context.Background(), &messages.PortfolioRequest{
		Options: map[string]string{"token": tinkoffToken},
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}
