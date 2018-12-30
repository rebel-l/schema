package integration

import (
	"bufio"
	"os"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/rebel-l/go-utils/osutils"
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

	file, err := os.Open(schemaFile)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	var buffer string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) > 0 {
			buffer += "\n" + line
		}
	}

	_, err = db.Exec(buffer)
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
