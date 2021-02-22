package models

import (
	"context"
	"fmt"

	"github.com/fabritsius/investor/aggregator/models"
	"github.com/gocql/gocql"
)

// DB stores the session and model methods are defined on it
type DB struct {
	session *gocql.Session
}

// Connect opens a new DB session
func Connect(hosts ...string) (*DB, error) {
	db := &DB{}
	if err := models.OpenSession(db, hosts); err != nil {
		return nil, err
	}

	return db, nil
}

// Disconnect closes the DB session
func (db *DB) Disconnect() {
	models.CloseSession(db)
}

// Init creates missing tables and fills in default values
func (db *DB) Init() error {
	if err := models.EnsureStats(context.Background(), db); err != nil {
		return fmt.Errorf("failed to unsure stats tables: %s", err)
	}

	return nil
}

// GetSession return session to satisfy HasSession interface
func (db *DB) GetSession() *gocql.Session {
	return db.session
}

// SetSession sets the session to satisfy HasSession interface
func (db *DB) SetSession(session *gocql.Session) error {
	db.session = session
	return nil
}
