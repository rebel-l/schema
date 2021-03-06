package store_test

import (
	"testing"
	"time"

	"github.com/rebel-l/schema/store"
	"github.com/rebel-l/schema/utils/testdb"
)

func equalDateTime(expected time.Time, actual time.Time) bool {
	e := expected.Format(time.ANSIC)
	a := actual.Format(time.ANSIC)

	return e == a
}

func TestSchemaScriptMapper_Add_Integration(t *testing.T) {
	if testing.Short() {
		t.Skipf("skipped because of long running")
	}

	t.Parallel()

	db, err := testdb.InitDB("./testdata/tmp/add_integration_tests.db")
	if err != nil {
		t.Fatalf("not able to open database connection: %s", err)
	}

	defer testdb.ShutdownDB(db, t)

	// now the test
	expected := store.NewSchemaScriptSuccess("some_script.sql", "0.5.2")

	vm := store.NewSchemaScriptMapper(db)

	err = vm.Add(expected)
	if err != nil {
		t.Fatalf("No error expected on adding entry to database: %s", err)
	}

	if expected.ID != 1 {
		t.Errorf("Expected that id is set with 1 but got %d", expected.ID)
	}
}

func TestSchemaScriptMapper_Remove_Integration(t *testing.T) {
	if testing.Short() {
		t.Skipf("skipped because of long running")
	}

	t.Parallel()

	db, err := testdb.InitDB("./testdata/tmp/remove_integration_tests.db")
	if err != nil {
		t.Fatalf("not able to open database connection: %s", err)
	}

	defer testdb.ShutdownDB(db, t)

	// init data
	script := store.NewSchemaScriptSuccess("my.sql", "")

	sm := store.NewSchemaScriptMapper(db)
	if err = sm.Add(script); err != nil {
		t.Fatalf("Failed to prepare data: %s", err)
	}

	before, err := sm.GetAll()
	if err != nil {
		t.Fatalf("Failed to get all data: %s", err)
	}

	// now the test
	if err = sm.Remove(script.ScriptName); err != nil {
		t.Fatalf("Failed to remove entry: %s", err)
	}

	after, err := sm.GetAll()
	if err != nil {
		t.Fatalf("Failed to get all data again: %s", err)
	}

	if after.Len() >= before.Len() {
		t.Errorf("Expected number of rows before (%d) is greater than after (%d)", before.Len(), after.Len())
	}
}

func TestSchemaScriptMapper_GetByID_Integration(t *testing.T) { // nolint: gocognit
	if testing.Short() {
		t.Skipf("skipped because of long running")
	}

	t.Parallel()

	testcases := []struct {
		name     string
		dbFile   string
		expected *store.SchemaScript
	}{
		{
			name:     "success entry",
			dbFile:   "./testdata/tmp/get_success_integration_tests.db",
			expected: store.NewSchemaScriptSuccess("success.sql", ""),
		},
		{
			name:     "success entry with app version",
			dbFile:   "./testdata/tmp/get_success_with_app_version_integration_tests.db",
			expected: store.NewSchemaScriptSuccess("success.sql", "0.8.11"),
		},
		{
			name:     "error entry",
			dbFile:   "./testdata/tmp/get_error_integration_tests.db",
			expected: store.NewSchemaScriptError("error.sql", "", "an error message"),
		},
		{
			name:     "error entry with app version",
			dbFile:   "./testdata/tmp/get_error_with_app_version_integration_tests.db",
			expected: store.NewSchemaScriptError("error.sql", "master-20190212-2354", "an error message"),
		},
	}

	for _, testcase := range testcases {
		dbFile := testcase.dbFile
		expected := testcase.expected
		t.Run(testcase.name, func(t *testing.T) {
			db, err := testdb.InitDB(dbFile)
			if err != nil {
				t.Fatalf("not able to open database connection: %s", err)
			}
			defer testdb.ShutdownDB(db, t)

			vm := store.NewSchemaScriptMapper(db)
			err = vm.Add(expected)
			if err != nil {
				t.Fatalf("No error expected on adding entry to database: %s", err)
			}

			if expected.ID != 1 {
				t.Errorf("Expected that id is set with 1 but got %d", expected.ID)
			}

			actual, err := vm.GetByID(expected.ID)
			if err != nil {
				t.Fatalf("No error expected on loading entry from database: %s", err)
			}

			if actual == nil {
				t.Fatalf("Loaded entry for id %d should not be nil", expected.ID)
			}

			if expected.ID != actual.ID {
				t.Errorf("IDs should be identical: expected %d but got %d", expected.ID, actual.ID)
			}

			if expected.ErrorMsg != actual.ErrorMsg {
				t.Errorf("Expected error message '%s' but got '%s'", expected.ErrorMsg, actual.ErrorMsg)
			}

			if expected.Status != actual.Status {
				t.Errorf("Expected status '%s' but got '%s'", expected.Status, actual.Status)
			}

			if !equalDateTime(expected.ExecutedAt, actual.ExecutedAt) {
				t.Errorf("Expected executed at '%s' but got '%s'", expected.ExecutedAt.String(), actual.ExecutedAt.String())
			}

			if expected.ScriptName != actual.ScriptName {
				t.Errorf("Expected script name '%s' but got '%s'", expected.ScriptName, actual.ScriptName)
			}

			if expected.AppVersion != actual.AppVersion {
				t.Errorf("Expected app version '%s' but got '%s'", expected.AppVersion, actual.AppVersion)
			}
		})
	}
}

func TestSchemaScriptMapper_GetAll_Integration(t *testing.T) {
	if testing.Short() {
		t.Skipf("skipped because of long running")
	}

	t.Parallel()

	db, err := testdb.InitDB("./testdata/tmp/getall_integration_tests.db")
	if err != nil {
		t.Fatalf("not able to open database connection: %s", err)
	}

	defer testdb.ShutdownDB(db, t)

	expected := []*store.SchemaScript{
		store.NewSchemaScriptSuccess("success.sql", "0.7.3"),
		store.NewSchemaScriptError("error.sql", "", "a message"),
	}

	vm := store.NewSchemaScriptMapper(db)

	for _, v := range expected {
		err = vm.Add(v)
		if err != nil {
			t.Fatalf("No error expected on adding entry to database: %s", err)
		}
	}

	actual, err := vm.GetAll()
	if err != nil {
		t.Fatalf("No error expected on loading entries from database: %s", err)
	}

	if len(expected) != len(actual) {
		t.Fatalf("Expected number of entries %d but got %d", len(expected), len(actual))
	}

	for k, e := range expected {
		a := actual[k]

		if e.ID != a.ID {
			t.Errorf("IDs should be identical: expected %d but got %d", e.ID, a.ID)
		}

		if e.ErrorMsg != a.ErrorMsg {
			t.Errorf("Expected error message '%s' but got '%s'", e.ErrorMsg, a.ErrorMsg)
		}

		if e.Status != a.Status {
			t.Errorf("Expected status '%s' but got '%s'", e.Status, a.Status)
		}

		if !equalDateTime(e.ExecutedAt, a.ExecutedAt) {
			t.Errorf("Expected executed at '%s' but got '%s'", e.ExecutedAt.String(), a.ExecutedAt.String())
		}

		if e.ScriptName != a.ScriptName {
			t.Errorf("Expected script name '%s' but got '%s'", e.ScriptName, a.ScriptName)
		}

		if e.AppVersion != a.AppVersion {
			t.Errorf("Expected app version '%s' but got '%s'", e.AppVersion, a.AppVersion)
		}
	}
}
