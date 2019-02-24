package initdb_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/rebel-l/schema/store"

	"github.com/golang/mock/gomock"

	"github.com/rebel-l/schema/initdb"
	"github.com/rebel-l/schema/mocks/store_mock"
	"github.com/rebel-l/schema/tests/integration"
)

func TestInitDB_ApplyScript_Integration_Happy(t *testing.T) {
	if testing.Short() {
		t.Skip("skipped because of long running")
	}

	db, err := integration.InitDB("./../tests/data/storage/apply_script_integration.db")
	if err != nil {
		t.Fatalf("Failed to open database: %s", err)
	}
	defer integration.ShutdownDB(db, t)

	in := initdb.New(db)
	err = in.ApplyScript("./../tests/data/initdb/001.sql")
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	}

	var counter []uint32
	q := db.Rebind("SELECT count(id) FROM something;")
	err = db.Select(&counter, q)
	if err != nil {
		t.Fatalf("not able count rows in table: %s", err)
	}

	if len(counter) == 0 || counter[0] != 0 {
		t.Error("not able to select from table")
	}
}

func TestInitDB_ApplyScript_Integration_Unhappy(t *testing.T) {
	if testing.Short() {
		t.Skip("skipped because of long running")
	}

	testCases := []struct {
		name       string
		scriptName string
		expected   string
		dbErrorMsg string
	}{
		{
			name:       "not existing script",
			scriptName: "no_existing.sql",
			expected:   "open no_existing.sql:",
		},
		{
			name:       "database error",
			scriptName: "./../tests/data/initdb/001.sql",
			expected:   "something happened",
			dbErrorMsg: "something happened",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl, mockDB := getMockDB(t, testCase.dbErrorMsg)
			defer ctrl.Finish()

			in := initdb.New(mockDB)
			err := in.ApplyScript(testCase.scriptName)
			if err == nil {
				t.Error("Expected that error is returned")
			}

			if err != nil && !strings.Contains(err.Error(), testCase.expected) {
				t.Errorf("Expected error message '%s' but got '%s'", testCase.expected, err)
			}
		})
	}
}

func TestInitDB_RevertScript_Integration_Happy(t *testing.T) {
	if testing.Short() {
		t.Skip("skipped because of long running")
	}

	db, err := integration.InitDB("./../tests/data/storage/revert_script_integration.db")
	if err != nil {
		t.Fatalf("Failed to open database: %s", err)
	}
	defer integration.ShutdownDB(db, t)

	in := initdb.New(db)
	err = in.RevertScript("./../tests/data/initdb/001.sql")
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	}

	var counter []uint32
	q := db.Rebind("SELECT count(id) FROM something;")
	err = db.Select(&counter, q)
	if err == nil {
		t.Fatalf("table wasn't dropped: %s", err)
	}
}

func TestInitDB_RevertScript_Integration_Unhappy(t *testing.T) {
	if testing.Short() {
		t.Skip("skipped because of long running")
	}

	testCases := []struct {
		name       string
		scriptName string
		expected   string
		dbErrorMsg string
	}{
		{
			name:       "not existing script",
			scriptName: "no_existing.sql",
			expected:   "open no_existing.sql:",
		},
		{
			name:       "database error",
			scriptName: "./../tests/data/initdb/001.sql",
			expected:   "something happened",
			dbErrorMsg: "something happened",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl, mockDB := getMockDB(t, testCase.dbErrorMsg)
			defer ctrl.Finish()

			in := initdb.New(mockDB)
			err := in.RevertScript(testCase.scriptName)
			if err == nil {
				t.Error("Expected that error is returned")
			}

			if err != nil && !strings.Contains(err.Error(), testCase.expected) {
				t.Errorf("Expected error message '%s' but got '%s'", testCase.expected, err)
			}
		})
	}
}

func getMockDB(t *testing.T, errorMsg string) (*gomock.Controller, *store_mock.MockDatabaseConnector) {
	ctrl := gomock.NewController(t)
	mockDB := store_mock.NewMockDatabaseConnector(ctrl)
	if errorMsg != "" {
		mockDB.EXPECT().Exec(gomock.Any()).Return(nil, errors.New("something happened"))
	}
	return ctrl, mockDB
}

func TestInitDB_Init_Happy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	q := `CREATE TABLE IF NOT EXISTS schema_script (
  			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, -- TODO: AUTOINCREMENT is not available in every database
  			script_name TEXT NOT NULL,
  			executed_at DATETIME NOT NULL,
  			execution_status VARCHAR(100) NOT NULL,
  			app_version CHAR(30) NULL,
  			error_msg TEXT NULL
		);`

	mockDB := store_mock.NewMockDatabaseConnector(ctrl)
	mockDB.EXPECT().Exec(q).Return(nil, nil)

	in := initdb.New(mockDB)
	if err := in.Init(); err != nil {
		t.Errorf("Expected no error but got %s", err)
	}
}

func TestInitDB_Init_Unhappy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := store_mock.NewMockDatabaseConnector(ctrl)
	mockDB.EXPECT().Exec(gomock.Any()).Return(nil, errors.New("something happened"))

	in := initdb.New(mockDB)
	if err := in.Init(); err == nil {
		t.Error("Expected that errors are returned")
	}
}

func TestInitDB_Init_Integration_Happy(t *testing.T) {
	if testing.Short() {
		t.Skip("skipped because of long running")
	}

	db, err := integration.GetDB("./../tests/data/storage/init_integration.db")
	if err != nil {
		t.Fatalf("not able to open database connection: %s", err)
	}
	defer integration.ShutdownDB(db, t)

	in := initdb.New(db)
	err = in.Init()
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	}

	var counter []uint32
	q := db.Rebind("SELECT count(id) FROM schema_script;")
	err = db.Select(&counter, q)
	if err != nil {
		t.Fatalf("not able count rows in table: %s", err)
	}

	if len(counter) == 0 || counter[0] != 0 {
		t.Error("not able to select from table")
	}
}

func TestInitDB_ReInit_Happy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	q1 := `DROP TABLE IF EXISTS schema_script;`
	q2 := `CREATE TABLE IF NOT EXISTS schema_script (
  			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, -- TODO: AUTOINCREMENT is not available in every database
  			script_name TEXT NOT NULL,
  			executed_at DATETIME NOT NULL,
  			execution_status VARCHAR(100) NOT NULL,
  			app_version CHAR(30) NULL,
  			error_msg TEXT NULL
		);`

	mockDB := store_mock.NewMockDatabaseConnector(ctrl)
	mockDB.EXPECT().Exec(q1).Return(nil, nil)
	mockDB.EXPECT().Exec(q2).Return(nil, nil)

	in := initdb.New(mockDB)
	if err := in.ReInit(); err != nil {
		t.Errorf("Expected no error but got %s", err)
	}
}

func TestInitDB_ReInit_Unhappy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := store_mock.NewMockDatabaseConnector(ctrl)
	mockDB.EXPECT().Exec(gomock.Any()).Return(nil, errors.New("something happened"))

	in := initdb.New(mockDB)
	if err := in.ReInit(); err == nil {
		t.Error("Expected that errors are returned")
	}
}

func TestInitDB_ReInit_Integration_Happy(t *testing.T) {
	if testing.Short() {
		t.Skip("skipped because of long running")
	}

	db, err := integration.GetDB("./../tests/data/storage/reinit_integration.db")
	if err != nil {
		t.Fatalf("not able to open database connection: %s", err)
	}
	defer integration.ShutdownDB(db, t)

	// prepare
	in := initdb.New(db)
	if err = in.Init(); err != nil {
		t.Fatalf("Prepare: init failed %s", err)
	}

	m := store.NewSchemaScriptMapper(db)
	if err = m.Add(store.NewSchemaScriptSuccess("something.sql", "")); err != nil {
		t.Fatalf("Prepare: failed to add entry %s", err)
	}

	var counter []uint32
	q := db.Rebind("SELECT count(id) FROM schema_script;")
	if err = db.Select(&counter, q); err != nil {
		t.Fatalf("Prepare: not able count rows in table: %s", err)
	}

	if len(counter) == 0 || counter[0] != 1 {
		t.Fatalf("expected number of %d rows but got %d", 1, counter[0])
	}

	// now the test
	if err = in.ReInit(); err != nil {
		t.Fatalf("Failed to reinit: %s", err)
	}

	counter = make([]uint32, 0)
	if err = db.Select(&counter, q); err != nil {
		t.Fatalf("not able count rows in table: %s", err)
	}

	if len(counter) == 0 || counter[0] != 0 {
		t.Errorf("expected number of %d rows but got %d", 0, counter[0])
	}
}
