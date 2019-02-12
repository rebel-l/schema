package store

import (
	"fmt"
)

const dateTimeFormat = "2006-01-02T15:04:05Z"

// SchemaScriptMapper is responsible for mapping and storing SchemaScript struct in database
type SchemaScriptMapper struct {
	db DatabaseConnector
}

// NewSchemaScriptMapper returns a new SchemaScriptMapper
func NewSchemaScriptMapper(db DatabaseConnector) SchemaScriptMapper {
	return SchemaScriptMapper{db: db}
}

// Add adds a new row to the table
func (svm SchemaScriptMapper) Add(entry *SchemaScript) error {
	if entry == nil {
		return fmt.Errorf("SchemaScriptMapper, add: dataset must be provided")
	}

	q := `
		INSERT INTO schema_script (
			script_name,
  			executed_at,
  			execution_status,
  			error_msg
		) VALUES (?, ?, ?, ?)
	`

	res, err := svm.db.Exec(
		q,
		entry.ScriptName,
		entry.ExecutedAt.Format(dateTimeFormat),
		entry.Status,
		entry.ErrorMsg,
	)

	if err != nil {
		return fmt.Errorf("SchemaScriptMapper, add failed: %s", err)
	}

	entry.ID, err = res.LastInsertId()
	if err != nil {
		return fmt.Errorf("SchemaScriptMapper, add returns no new id: %s", err)
	}

	return nil
}

// GetByID returns the SchemaScript entry found for provided id
func (svm SchemaScriptMapper) GetByID(id int64) (*SchemaScript, error) {
	if id < 1 {
		return nil, fmt.Errorf("SchemaScriptMapper, get by id: id must be greater than zero")
	}

	sv := &SchemaScript{}
	q := `SELECT * from schema_script WHERE id = ?`
	if err := svm.db.Get(sv, q, id); err != nil {
		return nil, fmt.Errorf("SchemaScriptMapper, get by id failed: %s", err)
	}

	return sv, nil
}

// GetAll returns all SchemaScript entries
func (svm SchemaScriptMapper) GetAll() ([]*SchemaScript, error) {
	var versions []*SchemaScript
	q := `SELECT * FROM schema_script`
	if err := svm.db.Select(&versions, q); err != nil {
		return nil, err
	}

	return versions, nil
}
