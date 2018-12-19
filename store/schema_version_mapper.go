package store

import (
	"fmt"
)

const dateTimeFormat = "2006-01-02T15:04:05Z"

// SchemaVersionMapper is responsible for mapping and storing SchemaVersion struct in database
type SchemaVersionMapper struct {
	db DatabaseConnector
}

// NewSchemaVersionMapper returns a new SchemaVersionMapper
func NewSchemaVersionMapper(db DatabaseConnector) SchemaVersionMapper {
	return SchemaVersionMapper{db: db}
}

// Add adds a new row to the table
func (svm SchemaVersionMapper) Add(entry *SchemaVersion) error {
	if entry == nil {
		return fmt.Errorf("SchemaVersion, add: dataset must be provided")
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
		return fmt.Errorf("SchemaVersion, add failed: %s", err)
	}

	entry.ID, err = res.LastInsertId()
	if err != nil {
		return fmt.Errorf("SchemaVersion, add returns no new id: %s", err)
	}

	return nil
}

// GetByID returns the SchemaVersion entry found for provided id
func (svm SchemaVersionMapper) GetByID(id int64) (*SchemaVersion, error) {
	if id < 1 {
		return nil, fmt.Errorf("SchemaVersion, get by id: id must be greater than zero")
	}

	sv := &SchemaVersion{}
	q := `SELECT * from schema_version WHERE id = ?`
	if err := svm.db.Select(sv, q, id); err != nil {
		return nil, fmt.Errorf("SchemaVersion, get by id failed: %s", err)
	}

	return sv, nil
}

// GetAll returns all SchemaVersion entries
func (svm SchemaVersionMapper) GetAll() ([]*SchemaVersion, error) {
	var versions []*SchemaVersion
	q := `SELECT * FROM schema_version`
	if err := svm.db.Select(versions, q); err != nil {
		return nil, err
	}

	return versions, nil
}
