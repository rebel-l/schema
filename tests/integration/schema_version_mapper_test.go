package integration

import (
	"testing"
	"time"

	"github.com/rebel-l/schema/store"

	_ "github.com/mattn/go-sqlite3"
)

func equalDateTime(expected time.Time, actual time.Time) bool {
	e := expected.Format(time.ANSIC)
	a := actual.Format(time.ANSIC)

	return e == a
}

func TestSchemaVersionMapper_Add_Integration(t *testing.T) {
	if testing.Short() {
		t.Skipf("skipped because of long running")
	}

	db, err := initDB("./../data/storage/add_integration_tests.db")
	if err != nil {
		t.Fatalf("not able to open database connection: %s", err)
	}
	defer shutdownDB(db, t)

	// now the test
	expected := store.NewSchemaVersionSuccess("some_script.sql")

	vm := store.NewSchemaVersionMapper(db)
	err = vm.Add(expected)
	if err != nil {
		t.Fatalf("No error expected on adding entry to database: %s", err)
	}

	if expected.ID != 1 {
		t.Errorf("Expected that id is set with 1 but got %d", expected.ID)
	}
}

func TestSchemaVersionMapper_GetByID_Integration(t *testing.T) {
	if testing.Short() {
		t.Skipf("skipped because of long running")
	}

	testcases := []struct {
		name     string
		dbFile   string
		expected *store.SchemaVersion
	}{
		{
			name:     "success entry",
			dbFile:   "./../data/storage/get_success_integration_tests.db",
			expected: store.NewSchemaVersionSuccess("success.sql"),
		},
		{
			name:     "error entry",
			dbFile:   "./../data/storage/get_error_integration_tests.db",
			expected: store.NewSchemaVersionError("error.sql", "an error message"),
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			db, err := initDB(testcase.dbFile)
			if err != nil {
				t.Fatalf("not able to open database connection: %s", err)
			}
			defer shutdownDB(db, t)

			expected := testcase.expected
			vm := store.NewSchemaVersionMapper(db)
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
				t.Errorf("Expected script name '%s' but git '%s'", expected.ScriptName, actual.ScriptName)
			}
		})
	}
}

func TestSchemaVersionMapper_GetAll_Integration(t *testing.T) {
	if testing.Short() {
		t.Skipf("skipped because of long running")
	}

	db, err := initDB("./../data/storage/getall_integration_tests.db")
	if err != nil {
		t.Fatalf("not able to open database connection: %s", err)
	}
	defer shutdownDB(db, t)

	expected := []*store.SchemaVersion{
		store.NewSchemaVersionSuccess("success.sql"),
		store.NewSchemaVersionError("error.sql", "a message"),
	}

	vm := store.NewSchemaVersionMapper(db)
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
			t.Errorf("Expected script name '%s' but git '%s'", e.ScriptName, a.ScriptName)
		}
	}
}
