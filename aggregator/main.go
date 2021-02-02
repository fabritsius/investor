package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/fabritsius/envar"
	"github.com/fabritsius/investor/messages"
	"github.com/gocql/gocql"

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

	dbSession, err := getDbSession()
	if err != nil {
		log.Fatalf("did not connect to the database: %s", err)
	}
	defer dbSession.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	for account := range getUserAccounts(ctx, dbSession) {
		log.Println("Account:", account.userID, account.accountType, account.token)

		if account.accountType == "tinkoff" {
			portfolio, err := GetTinkoffPortfolio(conn, account.token)
			if err != nil {
				log.Fatalf("error when calling getTinkoffPortfolio: %s", err)
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

func getDbSession() (*gocql.Session, error) {
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "investor"
	return cluster.CreateSession()
}

func getUserAccounts(ctx context.Context, dbSession *gocql.Session) <-chan *UserAccount {
	query := "SELECT * FROM accounts_by_user"
	scanner := dbSession.Query(query).WithContext(ctx).Iter().Scanner()
	result := make(chan *UserAccount)
	go func() {
		defer close(result)
		for scanner.Next() {
			account := &UserAccount{}
			if err := scanner.Scan(&account.userID, &account.accountType, &account.token); err != nil {
				log.Fatalf("gocql scanner error: %s", err)
				continue
			}
			result <- account
		}
	}()
	return result
}

// UserAccount represents one of user's stocks account
type UserAccount struct {
	userID      gocql.UUID
	accountType string
	token       string
}
