// Package integration provides functions to setup or tear down integration tests
package integration

import (
	"os"
	"testing"

	"github.com/rebel-l/go-utils/osutils"
	"github.com/rebel-l/schema/initdb"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // sqlite3 driver is needed
)

// GetDB provides a database connection for integration tests
func GetDB(dbFile string) (*sqlx.DB, error) {
	if osutils.FileOrPathExists(dbFile) {
		err := os.Remove(dbFile)
		if err != nil {
			return nil, err
		}
	}

	_, err := os.Create(dbFile)
	if err != nil {
		return nil, err
	}

	db, err := sqlx.Open("sqlite3", dbFile)
	return db, err
}

// InitDB provides a database connection for integration tests and initialises the database
func InitDB(dbFile string) (*sqlx.DB, error) {
	db, err := GetDB(dbFile)
	if err != nil {
		return nil, err
	}

	in := initdb.New(db)
	err = in.Init()
	return db, err
}

// ShutdownDB closes database connection for integration tests
func ShutdownDB(db *sqlx.DB, t *testing.T) {
	if err := db.Close(); err != nil {
		t.Fatalf("failed to close database: %s", err)
	}
}
