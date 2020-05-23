package store

//go:generate mockgen -destination=../mocks/store_mock/database_connector_mock.go -package=store_mock github.com/rebel-l/schema/store DatabaseConnector

import (
	"io"

	"github.com/jmoiron/sqlx"
)

// DatabaseConnector provides methods to interact with a database.
type DatabaseConnector interface {
	sqlx.Execer
	Select(dest interface{}, query string, args ...interface{}) error
	Get(dest interface{}, query string, args ...interface{}) error
	io.Closer
	Rebind(string) string
}
