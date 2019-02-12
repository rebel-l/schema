// Package schema provides a library to organize and deploy your database schema
package schema

import (
	"github.com/rebel-l/schema/store"

	"github.com/sirupsen/logrus"
)

const (
	// CommandCreate is the command to create the schema
	CommandCreate = "create"

	// CommandMigrate is the command to apply the latest schema changes
	CommandMigrate = "migrate"

	// CommandRecreate is the command to recreate the schema
	CommandRecreate = "recreate"
)

// Versioner provides methods to manage the access to log of SQL script executions
type Versioner interface {
	Add(entry *store.SchemaScript) error
	GetByID(id int64) (*store.SchemaScript, error)
}

// Schema provides commands to organize your database schema
type Schema struct {
	PathOfSchemaFiles string
	Logger            logrus.FieldLogger
	Command           string
	versioner         Versioner
}

// New returns a Schema struct
func New(logger logrus.FieldLogger, db store.DatabaseConnector) Schema {
	return Schema{
		Logger:    logger,
		versioner: store.NewSchemaScriptMapper(db),
	}
}

// WithFlags initialises the CLI flags
func (s Schema) WithFlags() {
	// TODO: should be not in lib, but in cli tool
}
