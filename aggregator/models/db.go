package models

import (
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

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	return &DB{
		session: session,
	}, nil
}

// Close end the DB session
func (db *DB) Close() {
	db.session.Close()
}
