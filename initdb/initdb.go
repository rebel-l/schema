package initdb

import (
	"github.com/rebel-l/schema/sqlfile"
	"github.com/rebel-l/schema/store"
)

// InitDB provides functionality to initialize the database
type InitDB struct {
	db store.DatabaseConnector
}

// New returns an InitDB struct
func New(db store.DatabaseConnector) *InitDB {
	return &InitDB{
		db: db,
	}
}

// ApplyScript appliers a script to the database
func (i *InitDB) ApplyScript(fileName string) error {
	sqlScript, err := sqlfile.Read(fileName)
	if err != nil {
		return err
	}

	if _, err = i.db.Exec(sqlScript); err != nil {

	}
	return nil
}

// Init initializes the schema database
func (i *InitDB) Init() error {
	scripts := []string{
		`CREATE TABLE IF NOT EXISTS schema_script (
  			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, -- TODO: AUTOINCREMENT is not available in every database
  			script_name TEXT NOT NULL,
  			executed_at DATETIME NOT NULL,
  			execution_status VARCHAR(100) NOT NULL,
  			app_version CHAR(30) NULL,
  			error_msg TEXT NULL
		);`,
	}

	for _, q := range scripts {
		if _, err := i.db.Exec(q); err != nil {
			return err
		}
	}

	return nil
}
