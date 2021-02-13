package models

import (
	"context"
	"fmt"
	"time"

	"github.com/gocql/gocql"
)

// DB stores the session and model methods are defined on it
type DB struct {
	session *gocql.Session
}

// Connect creates a new session and returns a DB object
func Connect(hosts ...string) (*DB, error) {
	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = "investor"
	cluster.Timeout = 5 * time.Second

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	return &DB{
		session: session,
	}, nil
}

// Init creates missing tables
func (db *DB) Init() error {
	if err := db.EnsureUsers(context.Background()); err != nil {
		return fmt.Errorf("failed to unsure users tables: %s", err)
	}

	return nil
}

// Close end the DB session
func (db *DB) Close() {
	db.session.Close()
}
