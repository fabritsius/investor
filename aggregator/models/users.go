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

// EnsureUsers creates all user related tables if they are missing
func (db *DB) EnsureUsers(ctx context.Context) error {
	query := `CREATE TABLE IF NOT EXISTS accounts_by_user (
		user_id UUID,
		account text,
		key text,
		PRIMARY KEY (user_id, account));`
	return db.session.Query(query).WithContext(ctx).Exec()
}
