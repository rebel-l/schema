package store

//go:generate mockgen -destination=../mocks/database_connector_mock.go -package=mocks github.com/rebel-l/schema/store DatabaseConnector

import (
	"github.com/jmoiron/sqlx"
)

// DatabaseConnector provides methods to interact with a database
type DatabaseConnector interface {
	sqlx.Execer
	Select(dest interface{}, query string, args ...interface{}) error
}
