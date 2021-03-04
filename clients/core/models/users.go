package models

import (
	"context"

	"github.com/fabritsius/investor/aggregator/models"
	"github.com/gocql/gocql"
)

// GetAccountsForUser returns all accounts for a user
func (db *DB) GetAccountsForUser(ctx context.Context, userID gocql.UUID) ([]string, error) {
	query := "SELECT (user_id, account) FROM accounts_by_user WHERE user_id = ?;"
	scanner := db.session.Query(query, userID).WithContext(ctx).Iter().Scanner()
	result := []string{}
	for scanner.Next() {
		account := models.UserAccount{}
		if err := scanner.Scan(&account.UserID, &account.AccountType); err != nil {
			return result, err
		}
		result = append(result, account.AccountType)
	}
	return result, nil
}
