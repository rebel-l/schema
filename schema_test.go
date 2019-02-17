package schema_test

import (
	"testing"

	"github.com/rebel-l/schema/store"

	"github.com/rebel-l/schema/tests/integration"
	"github.com/sirupsen/logrus"

	"github.com/rebel-l/schema"
)

func TestSchema_Execute_CommandMigrate_Happy(t *testing.T) {
	if testing.Short() {
		t.Skip("skipped because of long running")
	}

	log := logrus.New()
	db, err := integration.InitDB("./tests/data/storage/schema_execute_migrate.db")
	if err != nil {
		t.Fatalf("failed to init database: %s", err)
	}
	defer integration.ShutdownDB(db, t)

	s := schema.New(log, db)
	err = s.Execute("./tests/data/schema", schema.CommandMigrate, "")
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	}

	// check entries in table schema_script
	var data store.SchemaScriptCollection
	q := "SELECT * FROM schema_script;"
	err = db.Select(&data, q)
	if err != nil {
		t.Fatalf("not able count rows in table: %s", err)
	}

	if len(data) != 2 {
		t.Fatalf("Expeted 2 rows in table but got %d", len(data))
	}

	expected := store.SchemaScriptCollection{
		&store.SchemaScript{
			ScriptName: "./tests/data/schema/001.sql",
			Status:     store.StatusSuccess,
		},
		&store.SchemaScript{
			ScriptName: "./tests/data/schema/002.sql",
			Status:     store.StatusSuccess,
		},
	}

	for i, v := range expected {
		actual := data[i]
		if v.ScriptName != actual.ScriptName {
			t.Errorf("Expected script name %s but got %s", v.ScriptName, actual.ScriptName)
		}

		if v.Status != actual.Status {
			t.Errorf("Expected status %s but got %s", v.Status, actual.Status)
		}

		if v.ErrorMsg != actual.ErrorMsg {
			t.Errorf("Expected error message %s but got %s", v.ErrorMsg, actual.ErrorMsg)
		}

		if v.AppVersion != actual.AppVersion {
			t.Errorf("Expected app version %s but got %s", v.AppVersion, actual.AppVersion)
		}
	}
}

/*
TODO:
	Integration tests:
		1. command create happy
		3. command recreate ==> later
*/
