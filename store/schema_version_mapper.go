package store

import (
	"fmt"
)

const dateTimeFormat = "2006-01-02T15:04:05Z"

// SchemaVersionMapper is responsible for mapping and storing SchemaScript struct in database
type SchemaVersionMapper struct {
	db DatabaseConnector
}

// NewSchemaVersionMapper returns a new SchemaVersionMapper
func NewSchemaVersionMapper(db DatabaseConnector) SchemaVersionMapper {
	return SchemaVersionMapper{db: db}
}

// Add adds a new row to the table
func (svm SchemaVersionMapper) Add(entry *SchemaScript) error {
	if entry == nil {
		return fmt.Errorf("SchemaScript, add: dataset must be provided")
	}

	q := `
		INSERT INTO schema_version (
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
		return fmt.Errorf("SchemaScript, add failed: %s", err)
	}

	entry.ID, err = res.LastInsertId()
	if err != nil {
		return fmt.Errorf("SchemaScript, add returns no new id: %s", err)
	}

	return nil
}

// GetByID returns the SchemaScript entry found for provided id
func (svm SchemaVersionMapper) GetByID(id int64) (*SchemaScript, error) {
	if id < 1 {
		return nil, fmt.Errorf("SchemaScript, get by id: id must be greater than zero")
	}

	sv := &SchemaScript{}
	q := `SELECT * from schema_version WHERE id = ?`
	if err := svm.db.Get(sv, q, id); err != nil {
		return nil, fmt.Errorf("SchemaScript, get by id failed: %s", err)
	}

	return sv, nil
}

// GetAll returns all SchemaScript entries
func (svm SchemaVersionMapper) GetAll() ([]*SchemaScript, error) {
	var versions []*SchemaScript
	q := `SELECT * FROM schema_version`
	if err := svm.db.Select(&versions, q); err != nil {
		return nil, err
	}

	return versions, nil
}
