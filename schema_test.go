package schema_test

import (
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"

	"github.com/rebel-l/schema"
	"github.com/rebel-l/schema/store"
	"github.com/rebel-l/schema/tests/integration"

	"github.com/sirupsen/logrus"
)

func TestSchema_Execute_CommandUpgrade_Happy(t *testing.T) {
	if testing.Short() {
		t.Skip("skipped because of long running")
	}

	log := logrus.New()
	db, err := integration.GetDB("./tests/data/storage/schema_execute_upgrade.db")
	if err != nil {
		t.Fatalf("failed to init database: %s", err)
	}
	defer integration.ShutdownDB(db, t)

	s := schema.New(log, db)
	err = s.Execute("./tests/data/schema", schema.CommandUpgrade, "")
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

	checkScriptTable(expected, data, t)
	checkTable("something", db, t)
	checkTable("something_new", db, t)
}

func checkScriptTable(expected store.SchemaScriptCollection, actual store.SchemaScriptCollection, t *testing.T) {
	if len(expected) != len(actual) {
		t.Fatalf("Expeted %d rows in table but got %d", len(expected), len(actual))
	}

	for i, v := range expected {
		w := actual[i]
		if v.ScriptName != w.ScriptName {
			t.Errorf("Expected script name %s but got %s", v.ScriptName, w.ScriptName)
		}

		if v.Status != w.Status {
			t.Errorf("Expected status %s but got %s", v.Status, w.Status)
		}

		if v.ErrorMsg != w.ErrorMsg {
			t.Errorf("Expected error message %s but got %s", v.ErrorMsg, w.ErrorMsg)
		}

		if v.AppVersion != w.AppVersion {
			t.Errorf("Expected app version %s but got %s", v.AppVersion, w.AppVersion)
		}
	}
}

func checkTable(tableName string, db *sqlx.DB, t *testing.T) {
	var counter []uint32
	q := db.Rebind(fmt.Sprintf("SELECT count(id) FROM %s;", tableName))
	err := db.Select(&counter, q)
	if err != nil {
		t.Fatalf("not able count rows in table: %s", err)
	}

	if len(counter) == 0 || counter[0] != 0 {
		t.Error("not able to select from table")
	}
}

/*
TODO:
	Integration tests:
		1. command upgrade happy with already executed scripts
		2. command revert happy
		3. command recreate ==> later
	Unit tests ...
*/
