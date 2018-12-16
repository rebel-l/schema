// Package schema provides a library to organize and deploy your database schema
package schema

import (
	"github.com/rebel-l/schema/store"

	"github.com/sirupsen/logrus"
)

const (
	// CommandCycle is the command to recreate the schema
	CommandCycle = "cycle"

	// CommandMigrate is the command to apply the latest schema changes
	CommandMigrate = "migrate"
)

// Schema provides commands to organize your database schema
type Schema struct {
	PathOfSchemaFiles string
	Logger            logrus.FieldLogger
	Command           string
	mapper            store.SchemaVersionMapper
}

// New returns a Schema struct
func New(logger logrus.FieldLogger, db store.DatabaseConnector) Schema {
	return Schema{
		Logger: logger,
		mapper: store.NewSchemaVersionMapper(db),
	}
}

// WithFlags initialises the CLI flags
func (s Schema) WithFlags() {
	// TODO
}
