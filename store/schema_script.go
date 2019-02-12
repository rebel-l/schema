package store

import "time"

const (
	// StatusSuccess is the status name for 'success'
	StatusSuccess = "success"

	// StatusError is the status name for 'error'
	StatusError = "error"
)

// SchemaScript represents the version information stored in the database
type SchemaScript struct {
	ID         int64     `db:"id"`
	ScriptName string    `db:"script_name"`
	ExecutedAt time.Time `db:"executed_at"`
	Status     string    `db:"execution_status"`
	ErrorMsg   string    `db:"error_msg"`
	AppVersion string    `db:"app_version"`
}

// NewSchemaScriptSuccess returns a new SchemaScript struct prepared for successful execution
func NewSchemaScriptSuccess(scriptName string, appVersion string) *SchemaScript {
	return &SchemaScript{
		ScriptName: scriptName,
		ExecutedAt: time.Now(),
		Status:     StatusSuccess,
		AppVersion: appVersion,
	}
}

// NewSchemaScriptError returns a new SchemaScript struct prepared for failed execution
func NewSchemaScriptError(scriptName string, appVersion string, errorMsg string) *SchemaScript {
	return &SchemaScript{
		ScriptName: scriptName,
		ExecutedAt: time.Now(),
		Status:     StatusError,
		ErrorMsg:   errorMsg,
		AppVersion: appVersion,
	}
}
