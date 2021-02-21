package models

import (
	"context"
	"math"

	"github.com/gocql/gocql"
)

// DailyPortfolioStats contains main portfolio stats for a single day
type DailyPortfolioStats struct {
	UserID      gocql.UUID
	AccountType string
	Date        string
	Avg         float64
	Last        float64
	Max         float64
	Min         float64
	N           int
}

// PortfolioValue contains total portfolio value
// with separate Invested and Yield fields
type PortfolioValue struct {
	UserID      gocql.UUID
	AccountType string
	Date        string
	Invested    float64
	Yield       float64
}

// UpdateDailyRecord updates daily AVG, MAX, MIN records for portfolios
func (db *DB) UpdateDailyRecord(ctx context.Context, portfolio *PortfolioValue) error {
	var avg, last, max, min float64
	var n int
	min = math.MaxFloat64

	getQuery := `SELECT (avg, last, max, min, n)
		FROM daily_portfolio_stats_by_user
		WHERE user_id = ? AND account = ? AND date = ?;`
	q := db.session.Query(getQuery, portfolio.UserID, portfolio.AccountType, portfolio.Date)
	q.WithContext(ctx).Scan(&avg, &last, &max, &min, &n)

	current := portfolio.Invested + portfolio.Yield
	if current == last {
		return nil
	}

	avg = recalcAverage(avg, current, n)
	max = math.Max(max, current)
	min = math.Min(min, current)
	n++

	updateQuery := `UPDATE daily_portfolio_stats_by_user
		SET avg = ?, last = ?, max = ?, min = ?, n = ?
		WHERE user_id = ? AND account = ? AND date = ?;`
	q = db.session.Query(updateQuery, avg, current, max, min, n, portfolio.UserID, portfolio.AccountType, portfolio.Date)
	if err := q.WithContext(ctx).Exec(); err != nil {
		return err
	}

	return nil
}

// EnsureStats creates all stats related tables if they are missing
func EnsureStats(ctx context.Context, db HasSession) error {
	query := `CREATE TABLE IF NOT EXISTS daily_portfolio_stats_by_user (
		user_id uuid,
		account text,
		date date,
		avg double,
		last double,
		max double,
		min double,
		n int,
		PRIMARY KEY (user_id, account, date))
		WITH CLUSTERING ORDER BY (account ASC, date ASC);`
	if err := db.GetSession().Query(query).WithContext(ctx).Exec(); err != nil {
		return err
	}

	return nil
}

func recalcAverage(old, new float64, n int) float64 {
	return (old*float64(n) + new) / (float64(n) + 1)
}
