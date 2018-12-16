package store

import (
	"github.com/jmoiron/sqlx"
)

// DatabaseConnector provides methods to interact with a database
type DatabaseConnector interface {
	sqlx.Execer
	Select(dest interface{}, query string, args ...interface{}) error
}
