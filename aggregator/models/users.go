package models

import (
	"context"
	"log"
	"os"

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
	if err := db.session.Query(query).WithContext(ctx).Exec(); err != nil {
		return err
	}

	if defaultTinkoffToken, ok := os.LookupEnv("DEFAULT_TINKOFF_TOKEN"); ok {
		defaultID, err := gocql.UUIDFromBytes([]byte("default-tinkoff-token")[:16])
		if err != nil {
			log.Printf("failed to generate UUID so zero-ID is used: %s", err)
		}
		query = `UPDATE accounts_by_user SET key = ? WHERE user_id = ? AND account = ?;`
		err = db.session.Query(query, defaultTinkoffToken, defaultID, "tinkoff").WithContext(ctx).Exec()
		if err != nil {
			return err
		}
	}

	return nil
}
