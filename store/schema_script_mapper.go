package store

import (
	"errors"
	"fmt"
)

const (
	// DateTimeFormat represents the default sql datetime format.
	DateTimeFormat = "2006-01-02T15:04:05Z"
)

var (
	// ErrNoDataset is used if no data for row is provided
	ErrNoDataset = errors.New("dataset must be provided")

	// ErrNoID is used if no or wrong id was provided
	ErrNoID = errors.New("id must be greater than zero")

	// ErrNoScript is used if no script name was provided
	ErrNoScript = errors.New("script name must be provided")
)

// SchemaScriptMapper is responsible for mapping and storing SchemaScript struct in database.
type SchemaScriptMapper struct {
	db DatabaseConnector
}

// NewSchemaScriptMapper returns a new SchemaScriptMapper.
func NewSchemaScriptMapper(db DatabaseConnector) *SchemaScriptMapper {
	return &SchemaScriptMapper{db: db}
}

// Add adds a new row to the table.
func (ssm SchemaScriptMapper) Add(entry *SchemaScript) error {
	if entry == nil {
		return fmt.Errorf("SchemaScriptMapper, add: %w", ErrNoDataset)
	}

	q := `
		INSERT INTO schema_script (
			script_name,
  			executed_at,
  			execution_status,
  			error_msg,
			app_version
		) VALUES (?, ?, ?, ?, ?)
	`

	res, err := ssm.db.Exec(
		q,
		entry.ScriptName,
		entry.ExecutedAt.Format(DateTimeFormat),
		entry.Status,
		entry.ErrorMsg,
		entry.AppVersion,
	)

	if err != nil {
		return fmt.Errorf("SchemaScriptMapper, add failed: %w", err)
	}

	entry.ID, err = res.LastInsertId()
	if err != nil {
		return fmt.Errorf("SchemaScriptMapper, add returns no new id: %w", err)
	}

	return nil
}

// Remove deletes an entry from table based on scriptName.
func (ssm *SchemaScriptMapper) Remove(scriptName string) error {
	if scriptName == "" {
		return fmt.Errorf("SchemaScriptMapper, remove: %w", ErrNoScript)
	}

	q := `DELETE FROM schema_script WHERE script_name = ?;`
	if _, err := ssm.db.Exec(q, scriptName); err != nil {
		return err
	}

	return nil
}

// GetByID returns the SchemaScript entry found for provided id.
func (ssm SchemaScriptMapper) GetByID(id int64) (*SchemaScript, error) {
	if id < 1 {
		return nil, fmt.Errorf("SchemaScriptMapper, get by id: %w", ErrNoID)
	}

	sv := &SchemaScript{}
	q := `SELECT * from schema_script WHERE id = ?`

	if err := ssm.db.Get(sv, q, id); err != nil {
		return nil, fmt.Errorf("SchemaScriptMapper, get by id failed: %w", err)
	}

	return sv, nil
}

// GetAll returns all SchemaScript entries.
func (ssm SchemaScriptMapper) GetAll() (SchemaScriptCollection, error) {
	var versions []*SchemaScript

	q := `SELECT * FROM schema_script`
	if err := ssm.db.Select(&versions, q); err != nil {
		return nil, err
	}

	return versions, nil
}
