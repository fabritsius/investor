package models

import (
	"context"

	"github.com/gocql/gocql"
)

// UserAccount represents one of user's stocks account
type UserAccount struct {
	UserID      gocql.UUID
	AccountType string
	Token       string
}

// UserAccountWithError contains UserAccount bundled with an error
type UserAccountWithError struct {
	Account *UserAccount
	Error   error
}

// GetAllUserAccounts returns all accounts for all users with a channel
func (db *DB) GetAllUserAccounts(ctx context.Context) <-chan *UserAccountWithError {
	query := "SELECT * FROM accounts_by_user"
	scanner := db.session.Query(query).WithContext(ctx).Iter().Scanner()
	result := make(chan *UserAccountWithError)
	go func() {
		defer close(result)
		for scanner.Next() {
			account := &UserAccount{}
			if err := scanner.Scan(&account.UserID, &account.AccountType, &account.Token); err != nil {
				result <- &UserAccountWithError{nil, err}
				continue
			}
			result <- &UserAccountWithError{account, nil}
		}
	}()
	return result
}
