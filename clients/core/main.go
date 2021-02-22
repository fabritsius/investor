package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fabritsius/envar"
	"github.com/fabritsius/investor/clients/core/models"
	"github.com/fabritsius/investor/clients/core/utils"
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

	daysBack := 15
	userAccounts, err := db.GetAccountsForUser(ctx, userID)
	if err != nil {
		log.Fatalf("failed to load accounts by userID: %s", err)
	}

	userRecords := db.GetUserDailyRecordsForPeriod(ctx, userID, userAccounts, &models.Period{
		Start: time.Now().Add(-time.Duration(daysBack) * 24 * time.Hour),
		End:   time.Now(),
	})

	dailyTotals := make(map[string]*DailyTotals)
	dates := make([]string, daysBack)
	for item := range userRecords {
		if item.Error != nil {
			log.Printf("account error: %s", item.Error)
			continue
		}

		record := item.Record

		var dailyAvgRecord *DailyTotals
		var ok bool
		if dailyAvgRecord, ok = dailyTotals[record.Date]; !ok {
			dailyAvgRecord = &DailyTotals{}
			dates = append(dates, record.Date)
		}

		switch record.AccountType {
		case "tinkoff":
			dailyAvgRecord.tinkoff += record.Avg
		case "ethereum":
			dailyAvgRecord.ethereum += record.Avg
		}

		dailyAvgRecord.total += record.Avg
		dailyTotals[record.Date] = dailyAvgRecord
	}

	var prevRecord *DailyTotals
	for _, date := range dates {
		if dailyRecord, ok := dailyTotals[date]; ok {
			if prevRecord == nil {
				prevRecord = dailyRecord
				continue
			}

			tinkoffPercDiff := utils.GetPercDiff(prevRecord.tinkoff, dailyRecord.tinkoff)
			ethereumPercDiff := utils.GetPercDiff(prevRecord.ethereum, dailyRecord.ethereum)
			totalPercDiff := utils.GetPercDiff(prevRecord.total, dailyRecord.total)

			colorFunc := utils.GreenString
			if strings.HasPrefix(totalPercDiff, "-") {
				colorFunc = utils.RedString
			}

			fmt.Printf(colorFunc("%s | $%.1f (%s%%) | $%.1f (%s%%) | $%.1f (%s%%)\n"), date, dailyRecord.tinkoff, tinkoffPercDiff, dailyRecord.ethereum, ethereumPercDiff, dailyRecord.total, totalPercDiff)

			prevRecord = dailyRecord
		}
	}
}

// DailyTotals contains combined and separate portfolio totals
type DailyTotals struct {
	tinkoff  float64
	ethereum float64
	total    float64
}
