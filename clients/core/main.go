package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/fabritsius/envar"
	"github.com/fabritsius/investor/clients/core/models"
	"github.com/gocql/gocql"
)

type config struct {
	UserID string `env:"DEFAULT_USER_ID"`
}

func main() {
	cfg := config{}
	if err := envar.Fill(&cfg); err != nil {
		log.Fatalln(err)
	}

	var db *models.DB
	var err error
	if db, err = models.Connect("127.0.0.1"); err != nil {
		log.Fatalf("did not connect to the database: %s", err)
	}
	defer db.Disconnect()

	if err := db.Init(); err != nil {
		log.Fatalf("database init error: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	userID, err := gocql.ParseUUID(cfg.UserID)
	if err != nil {
		log.Fatalf("failed to parse mandatory UserID: %s", err)
	}

	daysBack := 7
	userAccounts := []string{"tinkoff", "ethereum"}
	userRecords := db.GetUserDailyRecordsForPeriod(ctx, userID, userAccounts, &models.Period{
		Start: time.Now().Add(-time.Duration(daysBack) * 24 * time.Hour),
		End:   time.Now(),
	})

	for item := range userRecords {
		if item.Error != nil {
			log.Printf("account error: %s", item.Error)
			continue
		}

		record := item.Record
		fmt.Println(*record)
	}

	fmt.Println("\nClient done.")
}
