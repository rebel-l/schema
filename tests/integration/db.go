package integration

import (
	"os"
	"testing"

	"github.com/jmoiron/sqlx"

	"github.com/rebel-l/go-utils/osutils"
	"github.com/rebel-l/schema/sqlfile"
)

const (
	schemaFile = "./../../scripts/sql/001_create_table_schema_version.sql"
)

func initDB(dbFile string) (*sqlx.DB, error) {
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
	if err != nil {
		return nil, err
	}

	content, err := sqlfile.Read(schemaFile)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(content)
	if err != nil {
		return nil, err
	}

	return db, err
}

func shutdownDB(db *sqlx.DB, t *testing.T) {
	if err := db.Close(); err != nil {
		t.Fatalf("failed to close database: %s", err)
	}
}
