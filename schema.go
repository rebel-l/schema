// Package schema provides a library to organize and deploy your database schema
package schema

//go:generate mockgen -destination=mocks/schema_mock/schema_mock.go -package=schema_mock github.com/rebel-l/schema Applier,Scripter

import (
	"errors"
	"fmt"

	"github.com/rebel-l/schema/bar"
	"github.com/rebel-l/schema/initdb"
	"github.com/rebel-l/schema/sqlfile"
	"github.com/rebel-l/schema/store"

	"gopkg.in/cheggaaa/pb.v1"
)

const (
	// CommandRecreate is the command to recreate the schema
	CommandRecreate = "recreate"
)

// Scripter provides methods to manage the access to log of SQL script executions
type Scripter interface {
	Add(entry *store.SchemaScript) error
	GetAll() (store.SchemaScriptCollection, error)
	Remove(scriptName string) error
}

// Applier provides methods to apply sql script to database
type Applier interface {
	ApplyScript(fileName string) error
	RevertScript(fileName string) error
	Init() error
	ReInit() error
}

// Progressor provides methods to steer a progress bar
type Progressor interface {
	Increment() int
	FinishPrint(msg string)
}

// Schema provides commands to organize your database schema
type Schema struct {
	Scripter    Scripter
	Applier     Applier
	progressBar bool
	db          store.DatabaseConnector
}

// New returns a Schema struct
func New(db store.DatabaseConnector) Schema {
	return Schema{
		Scripter: store.NewSchemaScriptMapper(db),
		Applier:  initdb.New(db),
		db:       db,
	}
}

// WithProgressBar activate the progress bar
func (s *Schema) WithProgressBar() {
	s.progressBar = true
}

// Execute applies all sql scripts for a given folder
func (s *Schema) Execute(path string, command string, version string) error {
	var err error
	switch command {
	case CommandRecreate:
		err = s.Recreate(path, version)
	default:
		err = fmt.Errorf("command '%s' not found", command)
	}

	return err
}

// Upgrade applies new scripts to the database or if executed the first time applies all.
// A path to the sql scripts needs to be provided. It applies only files with ending ".sql", sub folders are ignored.
// The version of your application can be provided too, use empty string to ignore it.
func (s *Schema) Upgrade(path string, version string) error {
	if !checkDatabaseExists(s.db) {
		if err := s.Applier.Init(); err != nil {
			return err
		}
	}

	executedScripts, err := s.Scripter.GetAll()
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

	progressBar := s.startProgressBar(len(files))
	for _, f := range files {
		progressBar.Increment()
		if executedScripts.ScriptExecuted(f) {
			continue
		}

		if err = s.Applier.ApplyScript(f); err != nil {
			msg := fmt.Sprintf("failed to execute script %s: %s", f, err)
			if err := s.Scripter.Add(store.NewSchemaScriptError(f, version, err.Error())); err != nil {
				msg = fmt.Sprintf("original error: %s, following error: %s", msg, err)
			}

			return errors.New(msg)
		}

		if err = s.Scripter.Add(store.NewSchemaScriptSuccess(f, version)); err != nil {
			return err
		}
	}
	progressBar.FinishPrint("Schema Upgrade finished!")
	return nil
}

// RevertLast reverts the last applied script. If it is repeatedly called, it reverts every time one script: means if
// you run it twice it reverts the last two scripts and so on.
// A path to the sql scripts needs to be provided. It applies only files with ending ".sql", sub folders are ignored.
func (s *Schema) RevertLast(path string) error {
	return s.revert(path, 1)
}

// RevertAll reverts the all applied scripts.
// A path to the sql scripts needs to be provided. It applies only files with ending ".sql", sub folders are ignored.
func (s *Schema) RevertAll(path string) error {
	return s.revert(path, -1)
}

func (s *Schema) revert(path string, numOfScripts int) error {
	executedScripts, err := s.Scripter.GetAll()
	if err != nil {
		return err
	}

	/**
	1. scan files reverse
	2. iterate over files in directory
	2a. check if file is applied
	2b. if 2a) is true load each file revert from database
	2c. remove executed script from 2b) from store as success or error
	3. return after numOfScripts was reverted, -1 means all
	*/
	files, err := sqlfile.ScanReverse(path)
	if err != nil {
		return err
	}

	counter := 0
	progressBar := s.startProgressBar(numOfScripts)
	if numOfScripts < 1 {
		progressBar = s.startProgressBar(len(files))
	}
	for _, f := range files {
		progressBar.Increment()
		if !executedScripts.ScriptExecuted(f) {
			continue
		}

		if err = s.Applier.RevertScript(f); err != nil {
			return err
		}

		if err = s.Scripter.Remove(f); err != nil {
			return err
		}

		counter++
		if numOfScripts > 0 && counter >= numOfScripts {
			break
		}
	}
	progressBar.FinishPrint("Schema revert finished!")

	return nil
}

// Recreate reverts all applied scripts and apply them again. Internally it usues RevertAll() and Upgrade().
func (s *Schema) Recreate(path string, version string) error {
	var err error
	if err = s.RevertAll(path); err != nil {
		return err
	}

	if err = s.Applier.ReInit(); err != nil {
		return err
	}
	return s.Upgrade(path, version)
}

func checkDatabaseExists(db store.DatabaseConnector) bool {
	var counter []uint32
	q := "SELECT count(id) FROM schema_script;"
	if err := db.Select(&counter, q); err != nil {
		return false
	}

	return true
}

func (s *Schema) startProgressBar(count int) Progressor {
	if s.progressBar {
		return pb.StartNew(count)
	}

	return &bar.BlackHole{}
}
