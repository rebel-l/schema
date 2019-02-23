package schema_test

import (
	"fmt"
	"io"
	"os"
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
	err = s.Execute("./tests/data/schema/upgrade", schema.CommandUpgrade, "")
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	}

	data, err := s.Scripter.GetAll()
	if err != nil {
		t.Fatalf("not able get rows from table: %s", err)
	}

	expected := store.SchemaScriptCollection{
		&store.SchemaScript{
			ScriptName: "./tests/data/schema/upgrade/001.sql",
			Status:     store.StatusSuccess,
		},
		&store.SchemaScript{
			ScriptName: "./tests/data/schema/upgrade/002.sql",
			Status:     store.StatusSuccess,
		},
	}

	checkScriptTable(expected, data, t)
	checkTable("something", db, t)
	checkTable("something_new", db, t)
}

func TestSchema_Execute_CommandUpgrade_Happy_TwoSteps(t *testing.T) {
	if testing.Short() {
		t.Skip("skipped because of long running")
	}

	log := logrus.New()
	db, err := integration.GetDB("./tests/data/storage/schema_execute_upgrade_twosteps.db")
	if err != nil {
		t.Fatalf("failed to init database: %s", err)
	}
	defer integration.ShutdownDB(db, t)

	s := schema.New(log, db)

	/**
	STEP 1
	*/
	if err = copyFile("./tests/data/schema/upgrade/step1/001.sql", "./tests/data/schema/upgrade/two_steps/001.sql"); err != nil {
		t.Fatalf("failed to copy file: %s", err)
	}

	err = s.Execute("./tests/data/schema/upgrade/two_steps", schema.CommandUpgrade, "")
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	}

	data, err := s.Scripter.GetAll()
	if err != nil {
		t.Fatalf("not able get rows from table: %s", err)
	}

	expected := store.SchemaScriptCollection{
		&store.SchemaScript{
			ScriptName: "./tests/data/schema/upgrade/two_steps/001.sql",
			Status:     store.StatusSuccess,
		},
	}

	checkScriptTable(expected, data, t)
	checkTable("something", db, t)

	/**
	STEP 2
	*/
	if err = copyFile("./tests/data/schema/upgrade/step2/002.sql", "./tests/data/schema/upgrade/two_steps/002.sql"); err != nil {
		t.Fatalf("failed to copy file: %s", err)
	}

	err = s.Execute("./tests/data/schema/upgrade/two_steps", schema.CommandUpgrade, "")
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	}

	data, err = s.Scripter.GetAll()
	if err != nil {
		t.Fatalf("not able get rows from table: %s", err)
	}

	expected = append(expected, &store.SchemaScript{
		ScriptName: "./tests/data/schema/upgrade/two_steps/002.sql",
		Status:     store.StatusSuccess,
	})
	checkScriptTable(expected, data, t)
	checkTable("something", db, t)
	checkTable("something_new", db, t)

	// cleanup
	if err = os.Remove("./tests/data/schema/upgrade/two_steps/001.sql"); err != nil {
		t.Fatalf("Cleanup files failed: %s", err)
	}

	if err = os.Remove("./tests/data/schema/upgrade/two_steps/002.sql"); err != nil {
		t.Fatalf("Cleanup files failed: %s", err)
	}
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

// TODO: move this to go-utils
func copyFile(src, dest string) error {
	from, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		_ = from.Close()
	}()

	to, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer func() {
		_ = to.Close()
	}()

	_, err = io.Copy(to, from)
	return err
}

/*
TODO:
	Integration tests:
		2. command revert happy
		3. command recreate ==> later
*/
