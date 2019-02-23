package schema_test

import (
	"errors"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/rebel-l/schema"
	"github.com/rebel-l/schema/mocks/logrus_mock"
	"github.com/rebel-l/schema/mocks/schema_mock"
	"github.com/rebel-l/schema/mocks/store_mock"
	"github.com/rebel-l/schema/store"
	"github.com/rebel-l/schema/tests/integration"

	"github.com/golang/mock/gomock"

	"github.com/sirupsen/logrus"
)

func TestSchema_Execute_CommandUpgrade_Happy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := getMockDB(ctrl, false)
	mockDB.EXPECT().Select(gomock.Any(), gomock.Any()).Times(1).Return(errors.New("failed"))

	mockApplier := schema_mock.NewMockApplier(ctrl)
	mockApplier.EXPECT().Init().Times(1).Return(nil)
	mockApplier.EXPECT().ApplyScript(gomock.Eq("./tests/data/schema/unit/001.sql")).Times(1).Return(nil)
	mockApplier.EXPECT().ApplyScript(gomock.Eq("./tests/data/schema/unit/002.sql")).Times(1).Return(nil)

	mockScripter := schema_mock.NewMockScripter(ctrl)
	mockScripter.EXPECT().GetAll().Times(1).Return(store.SchemaScriptCollection{}, nil)
	mockScripter.EXPECT().Add(gomock.Any()).Times(2).Return(nil)

	s := schema.New(getMockLogger(ctrl, true), mockDB)
	s.Applier = mockApplier
	s.Scripter = mockScripter

	if err := s.Execute("./tests/data/schema/unit", schema.CommandUpgrade, ""); err != nil {
		t.Errorf("Expected no errors but got %s", err)
	}
}

func TestSchema_Execute_CommandUpgrade_Unhappy_GetAllError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := getMockDB(ctrl, false)
	mockDB.EXPECT().Select(gomock.Any(), gomock.Any()).Times(1).Return(nil)

	mockApplier := schema_mock.NewMockApplier(ctrl)
	mockApplier.EXPECT().Init().Times(1).Times(0)

	mockScripter := schema_mock.NewMockScripter(ctrl)
	mockScripter.EXPECT().GetAll().Times(1).Return(store.SchemaScriptCollection{}, errors.New("failed"))

	s := schema.New(getMockLogger(ctrl, true), mockDB)
	s.Applier = mockApplier
	s.Scripter = mockScripter

	if err := s.Execute("./tests/data/schema/unit", schema.CommandUpgrade, ""); err == nil {
		t.Error("Expected error is returned on failed database operation")
	}
}

func TestSchema_Execute_CommandUpgrade_Unhappy_ApplyError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expected := "failed apply"

	mockLogger := logrus_mock.NewMockFieldLogger(ctrl)
	mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any())

	mockDB := getMockDB(ctrl, false)
	mockDB.EXPECT().Select(gomock.Any(), gomock.Any()).Times(1).Return(nil)

	mockScripter := schema_mock.NewMockScripter(ctrl)
	res := store.SchemaScriptCollection{}
	mockScripter.EXPECT().GetAll().Times(1).Return(res, nil)
	mockScripter.EXPECT().Add(gomock.Any()).Return(nil)

	mockApplier := schema_mock.NewMockApplier(ctrl)

	mockApplier.EXPECT().ApplyScript("./tests/data/schema/unit/001.sql").Return(errors.New(expected))

	s := schema.New(mockLogger, mockDB)
	s.Applier = mockApplier
	s.Scripter = mockScripter

	err := s.Execute("./tests/data/schema/unit", schema.CommandUpgrade, "")
	if err == nil {
		t.Error("Expected error is returned on failed apply")
	}

	if err.Error() != expected {
		t.Errorf("Expected error message '%s' but got '%s'", expected, err.Error())
	}
}

func TestSchema_Execute_CommandUpgrade_Unhappy_AddError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expected := "failed add"

	mockLogger := logrus_mock.NewMockFieldLogger(ctrl)
	mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	mockDB := getMockDB(ctrl, false)
	mockDB.EXPECT().Select(gomock.Any(), gomock.Any()).Times(1).Return(nil)

	mockScripter := schema_mock.NewMockScripter(ctrl)
	res := store.SchemaScriptCollection{}
	mockScripter.EXPECT().GetAll().Times(1).Return(res, nil)
	mockScripter.EXPECT().Add(gomock.Any()).Return(errors.New(expected))

	mockApplier := schema_mock.NewMockApplier(ctrl)
	mockApplier.EXPECT().ApplyScript("./tests/data/schema/unit/001.sql").Return(nil)

	s := schema.New(mockLogger, mockDB)
	s.Applier = mockApplier
	s.Scripter = mockScripter

	err := s.Execute("./tests/data/schema/unit", schema.CommandUpgrade, "")
	if err == nil {
		t.Error("Expected error is returned on failed add")
	}

	if err.Error() != expected {
		t.Errorf("Expected error message '%s' but got '%s'", expected, err.Error())
	}
}

func TestSchema_Execute_CommandUpgrade_Unhappy_ApplyErrorAddError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expected1 := "failed apply"
	expected2 := "failed add"
	final := "original error: failed apply, follow up error: failed add"

	mockLogger := logrus_mock.NewMockFieldLogger(ctrl)
	mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any())

	mockDB := getMockDB(ctrl, false)
	mockDB.EXPECT().Select(gomock.Any(), gomock.Any()).Times(1).Return(nil)

	mockScripter := schema_mock.NewMockScripter(ctrl)
	res := store.SchemaScriptCollection{}
	mockScripter.EXPECT().GetAll().Times(1).Return(res, nil)
	mockScripter.EXPECT().Add(gomock.Any()).Return(errors.New(expected2))

	mockApplier := schema_mock.NewMockApplier(ctrl)
	mockApplier.EXPECT().ApplyScript("./tests/data/schema/unit/001.sql").Return(errors.New(expected1))

	s := schema.New(mockLogger, mockDB)
	s.Applier = mockApplier
	s.Scripter = mockScripter

	err := s.Execute("./tests/data/schema/unit", schema.CommandUpgrade, "")
	if err == nil {
		t.Error("Expected error is returned on failed apply")
	}

	if err.Error() != final {
		t.Errorf("Expected error message '%s' but got '%s'", final, err.Error())
	}
}

func TestSchema_Execute_CommandRevert_Unhappy_RevertError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockApplier := schema_mock.NewMockApplier(ctrl)
	mockApplier.EXPECT().RevertScript("./tests/data/schema/unit/002.sql").Return(errors.New("failed"))

	mockScripter := schema_mock.NewMockScripter(ctrl)
	res := store.SchemaScriptCollection{&store.SchemaScript{
		ScriptName: "./tests/data/schema/unit/002.sql",
		Status:     store.StatusSuccess,
	}}
	mockScripter.EXPECT().GetAll().Times(1).Return(res, nil)

	s := schema.New(getMockLogger(ctrl, true), getMockDB(ctrl, true))
	s.Applier = mockApplier
	s.Scripter = mockScripter

	if err := s.Execute("./tests/data/schema/unit", schema.CommandRevert, ""); err == nil {
		t.Error("Expected error is returned on failed revert")
	}
}

func TestSchema_Execute_CommandRevert_Unhappy_RemoveError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockApplier := schema_mock.NewMockApplier(ctrl)
	mockApplier.EXPECT().RevertScript("./tests/data/schema/unit/002.sql").Return(nil)

	mockScripter := schema_mock.NewMockScripter(ctrl)
	res := store.SchemaScriptCollection{&store.SchemaScript{
		ScriptName: "./tests/data/schema/unit/002.sql",
		Status:     store.StatusSuccess,
	}}
	mockScripter.EXPECT().GetAll().Times(1).Return(res, nil)
	mockScripter.EXPECT().Remove("./tests/data/schema/unit/002.sql").Return(errors.New("failed"))

	s := schema.New(getMockLogger(ctrl, true), getMockDB(ctrl, true))
	s.Applier = mockApplier
	s.Scripter = mockScripter

	if err := s.Execute("./tests/data/schema/unit", schema.CommandRevert, ""); err == nil {
		t.Error("Expected error is returned on failed remove")
	}
}

func TestSchema_Execute_Unhappy_NotExistingPath(t *testing.T) {
	testCases := []struct {
		name    string
		path    string
		command string
	}{
		{
			name:    "empty path - upgrade",
			command: schema.CommandUpgrade,
		},
		{
			name:    "path not exists - upgrade",
			path:    "not_existing_path",
			command: schema.CommandUpgrade,
		},
		{
			name:    "empty path - revert",
			command: schema.CommandRevert,
		},
		{
			name:    "path not exists - revert",
			path:    "not_existing_path",
			command: schema.CommandRevert,
		},
		{
			name:    "empty path - recreate",
			command: schema.CommandRecreate,
		},
		{
			name:    "path not exists - recreate",
			path:    "not_existing_path",
			command: schema.CommandRecreate,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockScripter := schema_mock.NewMockScripter(ctrl)
			mockScripter.EXPECT().GetAll().Times(1).Return(store.SchemaScriptCollection{}, nil)

			s := schema.New(getMockLogger(ctrl, true), getMockDB(ctrl, true))
			s.Scripter = mockScripter

			if testCase.command == schema.CommandUpgrade {
				mockDB := getMockDB(ctrl, false)
				mockDB.EXPECT().Select(gomock.Any(), gomock.Any()).Times(1).Return(nil)

				s = schema.New(getMockLogger(ctrl, true), mockDB)
				s.Scripter = mockScripter
			}

			if err := s.Execute(testCase.path, testCase.command, ""); err == nil {
				t.Errorf("Expected an error on call with not existing path")
			}
		})
	}
}

func TestSchema_Execute_Unhappy_WrongCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := getSchemaWithDummies(ctrl)

	testCases := []struct {
		name    string
		command string
	}{
		{
			name: "empty command",
		},
		{
			name:    "wrong command",
			command: "god_command",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if err := s.Execute("./", testCase.command, ""); err == nil {
				t.Errorf("Expected an error on call with wrong command")
			}
		})
	}
}

func getMockDB(ctrl *gomock.Controller, dummy bool) *store_mock.MockDatabaseConnector {
	db := store_mock.NewMockDatabaseConnector(ctrl)
	if dummy {
		db.EXPECT().Select(gomock.Any(), gomock.Any()).Times(0)
	}
	return db
}

func getMockLogger(ctrl *gomock.Controller, dummy bool) *logrus_mock.MockFieldLogger {
	log := logrus_mock.NewMockFieldLogger(ctrl)
	if dummy {
		log.EXPECT().Error(gomock.Any()).Times(0)
		log.EXPECT().Errorf(gomock.Any()).Times(0)

		log.EXPECT().Fatal(gomock.Any()).Times(0)
		log.EXPECT().Fatalf(gomock.Any()).Times(0)
	}
	return log
}

func getSchemaWithDummies(ctrl *gomock.Controller) schema.Schema {
	return schema.New(getMockLogger(ctrl, true), getMockDB(ctrl, true))
}

func TestSchema_Execute_Integration_CommandUpgrade_Happy(t *testing.T) {
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

func TestSchema_Execute_Integration_CommandUpgrade_Happy_TwoSteps(t *testing.T) {
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
