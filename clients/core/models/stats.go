package models

import (
	"context"
	"time"

	"github.com/fabritsius/investor/aggregator/models"
	"github.com/gocql/gocql"
)

// DailyPortfolioStatsWithError contains DailyPortfolioStats record bundled with an error
type DailyPortfolioStatsWithError struct {
	Record *models.DailyPortfolioStats
	Error  error
}

// GetUserDailyRecordsForPeriod returns
func (db *DB) GetUserDailyRecordsForPeriod(ctx context.Context, userID gocql.UUID, accounts []string, period *Period) <-chan *DailyPortfolioStatsWithError {
	getQuery := `SELECT (account, date, avg, last, max, min, n)
		FROM daily_portfolio_stats_by_user
		WHERE user_id = ? AND account IN ? AND date > ? AND date <= ?;`
	q := db.session.Query(getQuery, userID, accounts, period.Start, period.End)
	scanner := q.WithContext(ctx).Iter().Scanner()
	result := make(chan *DailyPortfolioStatsWithError)
	go func() {
		defer close(result)
		for scanner.Next() {
			r := &models.DailyPortfolioStats{
				UserID: userID,
			}
			if err := scanner.Scan(&r.AccountType, &r.Date, &r.Avg, &r.Last, &r.Max, &r.Min, &r.N); err != nil {
				result <- &DailyPortfolioStatsWithError{nil, err}
				continue
			}
			result <- &DailyPortfolioStatsWithError{r, nil}
		}
	}()
	return result
}

// Period contains a start and an end
type Period struct {
	Start time.Time
	End   time.Time
}
