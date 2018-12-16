// Package store is responsible to persist the data of schema package
package store

import "time"

const (
	// StatusSuccess is the status name for 'success'
	StatusSuccess = "success"

	// StatusError is the status name for 'error'
	StatusError = "error"
)

// SchemaVersion represents the version information stored in the database
type SchemaVersion struct {
	ID         int64     `db:"id"`
	ScriptName string    `db:"script_name"`
	ExecutedAt time.Time `db:"executed_at"`
	Status     string    `db:"execution_status"`
	ErrorMsg   string    `db:"error_msg"`
}

// NewSchemaVersionSuccess returns a new SchemaVersion struct prepared for successful execution
func NewSchemaVersionSuccess(scriptName string) *SchemaVersion {
	return &SchemaVersion{
		ScriptName: scriptName,
		ExecutedAt: time.Now(),
		Status:     StatusSuccess,
	}
}

// NewSchemaVersionError returns a new SchemaVersion struct prepared for failed execution
func NewSchemaVersionError(scriptName string, errorMsg string) *SchemaVersion {
	return &SchemaVersion{
		ScriptName: scriptName,
		ExecutedAt: time.Now(),
		Status:     StatusError,
		ErrorMsg:   errorMsg,
	}
}
