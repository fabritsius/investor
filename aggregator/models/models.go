package models

import (
	"errors"
	"reflect"
	"time"

	"github.com/gocql/gocql"
)

// HasSession can store and return gocql session
type HasSession interface {
	SetSession(*gocql.Session) error
	GetSession() *gocql.Session
}

// OpenSession creates a news db session and sets to HasSession interface
func OpenSession(db HasSession, hosts []string) error {
	if reflect.ValueOf(db).IsNil() {
		return errors.New("please pass non-nil HasSession object to OpenSession")
	}

	session, err := createSession(hosts)
	if err != nil {
		return err
	}

	return db.SetSession(session)
}

// CloseSession end the DB session
func CloseSession(db HasSession) {
	db.GetSession().Close()
}

// createSession returns a new session with couple predefined parameters
func createSession(hosts []string) (*gocql.Session, error) {
	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = "investor"
	cluster.Timeout = 5 * time.Second
	return cluster.CreateSession()
}
