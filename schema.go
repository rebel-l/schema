// Package schema provides a library to organize and deploy your database schema
package schema

//go:generate mockgen -destination=mocks/schema_mock/schema_mock.go -package=schema_mock github.com/rebel-l/schema Applier,Scripter

import (
	"fmt"

	"github.com/rebel-l/go-utils/osutils"
	"github.com/rebel-l/schema/initdb"
	"github.com/rebel-l/schema/sqlfile"
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

// Scripter provides methods to manage the access to log of SQL script executions
type Scripter interface {
	Add(entry *store.SchemaScript) error
	GetAll() (store.SchemaScriptCollection, error)
}

// Applier provides methods to apply sql script to database
type Applier interface {
	ApplyScript(fileName string) error
	Init() error
}

// Schema provides commands to organize your database schema
type Schema struct {
	Logger   logrus.FieldLogger
	scripter Scripter
	applier  Applier
	db       store.DatabaseConnector
}

// New returns a Schema struct
func New(logger logrus.FieldLogger, db store.DatabaseConnector) Schema {
	return Schema{
		Logger:   logger,
		scripter: store.NewSchemaScriptMapper(db),
		applier:  initdb.New(db),
		db:       db,
	}
}

// Execute applies all sql scripts for a given folder
func (s *Schema) Execute(path string, command string, version string) error {
	// check path
	if osutils.FileOrPathExists(path) == false {
		return fmt.Errorf("path '%s' doesn'T exists", path)
	}

	// preparation
	switch command {
	case CommandCreate:
		if checkDatabaseExists(s.db) {
			return fmt.Errorf("create database not possible if in use, to force please use command %s", CommandRecreate)
		}
	case CommandRecreate:
		// TODO: drop all tables in schema ==> rollback all scripts
	}

	if command != CommandMigrate {
		if err := s.applier.Init(); err != nil {
			return err
		}
	}

	dbScripts, err := s.scripter.GetAll()
	if err != nil {
		return err
	}

	/**
	1. scan files
	2. iterate over files in directory
	2a. check if file is applied
	2b. if 2a) is false load each file apply to database
	2c. store executed script from 2b) to database as success or error
	*/
	files, err := sqlfile.Scan(path)
	if err != nil {
		return err
	}

	for _, f := range files {
		if dbScripts.ScriptExecuted(f) {
			continue
		}

		if err = s.applier.ApplyScript(f); err != nil {
			s.Logger.Errorf("failed to execute script %s: %s", f, err)
			if err2 := s.scripter.Add(store.NewSchemaScriptError(f, version, err.Error())); err2 != nil {
				err = fmt.Errorf("original error: %s, follow up error: %s", err, err2)
			}

			return err
		}

		if err = s.scripter.Add(store.NewSchemaScriptSuccess(f, version)); err != nil {
			return err
		}
	}

	return nil
}

func checkDatabaseExists(db store.DatabaseConnector) bool {
	var counter []uint32
	q := "SELECT count(id) FROM schema_script;"
	if err := db.Select(&counter, q); err != nil {
		return false
	}

	return true
}
