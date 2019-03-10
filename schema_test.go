package schema_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/rebel-l/go-utils/osutils"

	"github.com/jmoiron/sqlx"
	"github.com/rebel-l/schema"
	"github.com/rebel-l/schema/mocks/schema_mock"
	"github.com/rebel-l/schema/mocks/store_mock"
	"github.com/rebel-l/schema/store"
	"github.com/rebel-l/schema/utils/testdb"

	"github.com/golang/mock/gomock"
)

func TestSchema_Upgrade_Happy(t *testing.T) {
	testCases := []struct {
		name            string
		withProgressBar bool
	}{
		{
			name: "without progress bar",
		},
		{
			name:            "with progress bar",
			withProgressBar: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDB := getMockDB(ctrl, false)
			mockDB.EXPECT().Select(gomock.Any(), gomock.Any()).Times(1).Return(errors.New("failed"))

			mockApplier := schema_mock.NewMockApplier(ctrl)
			mockApplier.EXPECT().Init().Times(1).Return(nil)
			mockApplier.EXPECT().ApplyScript(gomock.Eq("./testdata/unit/001.sql")).Times(1).Return(nil)
			mockApplier.EXPECT().ApplyScript(gomock.Eq("./testdata/unit/002.sql")).Times(1).Return(nil)

			mockScripter := schema_mock.NewMockScripter(ctrl)
			mockScripter.EXPECT().GetAll().Times(1).Return(store.SchemaScriptCollection{}, nil)
			mockScripter.EXPECT().Add(gomock.Any()).Times(2).Return(nil)

			s := schema.New(mockDB)
			if testCase.withProgressBar {
				s.WithProgressBar()
			}
			s.Applier = mockApplier
			s.Scripter = mockScripter

			if err := s.Upgrade("./testdata/unit", ""); err != nil {
				t.Errorf("Expected no errors but got %s", err)
			}
		})
	}
}

func TestSchema_Upgrade_Unhappy_GetAllError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := getMockDB(ctrl, false)
	mockDB.EXPECT().Select(gomock.Any(), gomock.Any()).Times(1).Return(nil)

	mockApplier := schema_mock.NewMockApplier(ctrl)
	mockApplier.EXPECT().Init().Times(1).Times(0)

	mockScripter := schema_mock.NewMockScripter(ctrl)
	mockScripter.EXPECT().GetAll().Times(1).Return(store.SchemaScriptCollection{}, errors.New("failed"))

	s := schema.New(mockDB)
	s.Applier = mockApplier
	s.Scripter = mockScripter

	if err := s.Upgrade("./testdata/unit", ""); err == nil {
		t.Error("Expected error is returned on failed database operation")
	}
}

func TestSchema_Upgrade_Unhappy_InitError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := getMockDB(ctrl, false)
	mockDB.EXPECT().Select(gomock.Any(), gomock.Any()).Return(errors.New("failed"))

	mockApplier := schema_mock.NewMockApplier(ctrl)
	mockApplier.EXPECT().Init().Return(errors.New("failed init"))

	mockScripter := schema_mock.NewMockScripter(ctrl)
	mockScripter.EXPECT().GetAll().Times(0)

	s := schema.New(mockDB)
	s.Applier = mockApplier
	s.Scripter = mockScripter

	if err := s.Upgrade("./testdata/unit", ""); err == nil {
		t.Error("Expected error is returned on failed database initialisation")
	}
}

func TestSchema_Upgrade_Unhappy_ApplyError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expected := "failed to execute script ./testdata/unit/001.sql: failed apply"

	mockDB := getMockDB(ctrl, false)
	mockDB.EXPECT().Select(gomock.Any(), gomock.Any()).Times(1).Return(nil)

	mockScripter := schema_mock.NewMockScripter(ctrl)
	res := store.SchemaScriptCollection{}
	mockScripter.EXPECT().GetAll().Times(1).Return(res, nil)
	mockScripter.EXPECT().Add(gomock.Any()).Return(nil)

	mockApplier := schema_mock.NewMockApplier(ctrl)

	mockApplier.EXPECT().ApplyScript("./testdata/unit/001.sql").Return(errors.New("failed apply"))

	s := schema.New(mockDB)
	s.Applier = mockApplier
	s.Scripter = mockScripter

	err := s.Upgrade("./testdata/unit", "")
	if err == nil {
		t.Error("Expected error is returned on failed apply")
	}

	if err.Error() != expected {
		t.Errorf("Expected error message '%s' but got '%s'", expected, err.Error())
	}
}

func TestSchema_Upgrade_Unhappy_AddError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expected := "failed add"

	mockDB := getMockDB(ctrl, false)
	mockDB.EXPECT().Select(gomock.Any(), gomock.Any()).Times(1).Return(nil)

	mockScripter := schema_mock.NewMockScripter(ctrl)
	res := store.SchemaScriptCollection{}
	mockScripter.EXPECT().GetAll().Times(1).Return(res, nil)
	mockScripter.EXPECT().Add(gomock.Any()).Return(errors.New(expected))

	mockApplier := schema_mock.NewMockApplier(ctrl)
	mockApplier.EXPECT().ApplyScript("./testdata/unit/001.sql").Return(nil)

	s := schema.New(mockDB)
	s.Applier = mockApplier
	s.Scripter = mockScripter

	err := s.Upgrade("./testdata/unit", "")
	if err == nil {
		t.Error("Expected error is not returned on failed add")
	}

	if err.Error() != expected {
		t.Errorf("Expected error message '%s' but got '%s'", expected, err.Error())
	}
}

func TestSchema_Upgrade_Unhappy_ApplyErrorAddError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	errMsg1 := "failed apply"
	errMsg2 := "failed add"
	expected := "original error: failed to execute script ./testdata/unit/001.sql: failed apply, following error: failed add"

	mockDB := getMockDB(ctrl, false)
	mockDB.EXPECT().Select(gomock.Any(), gomock.Any()).Times(1).Return(nil)

	mockScripter := schema_mock.NewMockScripter(ctrl)
	res := store.SchemaScriptCollection{}
	mockScripter.EXPECT().GetAll().Times(1).Return(res, nil)
	mockScripter.EXPECT().Add(gomock.Any()).Return(errors.New(errMsg2))

	mockApplier := schema_mock.NewMockApplier(ctrl)
	mockApplier.EXPECT().ApplyScript("./testdata/unit/001.sql").Return(errors.New(errMsg1))

	s := schema.New(mockDB)
	s.Applier = mockApplier
	s.Scripter = mockScripter

	err := s.Upgrade("./testdata/unit", "")
	if err == nil {
		t.Error("Expected error is returned on failed apply")
	}

	if err.Error() != expected {
		t.Errorf("Expected error message '%s' but got '%s'", expected, err.Error())
	}
}

func TestSchema_RevertLast_Unhappy_RevertError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockApplier := schema_mock.NewMockApplier(ctrl)
	mockApplier.EXPECT().RevertScript("./testdata/unit/002.sql").Return(errors.New("failed"))

	mockScripter := schema_mock.NewMockScripter(ctrl)
	res := store.SchemaScriptCollection{&store.SchemaScript{
		ScriptName: "./testdata/unit/002.sql",
		Status:     store.StatusSuccess,
	}}
	mockScripter.EXPECT().GetAll().Times(1).Return(res, nil)

	s := schema.New(getMockDB(ctrl, true))
	s.Applier = mockApplier
	s.Scripter = mockScripter

	if err := s.RevertLast("./testdata/unit"); err == nil {
		t.Error("Expected error is returned on failed revert")
	}
}

func TestSchema_RevertLast_Unhappy_RemoveError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockApplier := schema_mock.NewMockApplier(ctrl)
	mockApplier.EXPECT().RevertScript("./testdata/unit/002.sql").Return(nil)

	mockScripter := schema_mock.NewMockScripter(ctrl)
	res := store.SchemaScriptCollection{&store.SchemaScript{
		ScriptName: "./testdata/unit/002.sql",
		Status:     store.StatusSuccess,
	}}
	mockScripter.EXPECT().GetAll().Times(1).Return(res, nil)
	mockScripter.EXPECT().Remove("./testdata/unit/002.sql").Return(errors.New("failed"))

	s := schema.New(getMockDB(ctrl, true))
	s.Applier = mockApplier
	s.Scripter = mockScripter

	if err := s.RevertLast("./testdata/unit"); err == nil {
		t.Error("Expected error is returned on failed remove")
	}
}

func TestSchema_RevertLast_Unhappy_GetAllError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockApplier := schema_mock.NewMockApplier(ctrl)
	mockApplier.EXPECT().RevertScript("./testdata/unit/002.sql").Times(0)

	mockScripter := schema_mock.NewMockScripter(ctrl)
	mockScripter.EXPECT().GetAll().Return(nil, errors.New("failed getting data"))

	s := schema.New(getMockDB(ctrl, true))
	s.Applier = mockApplier
	s.Scripter = mockScripter

	if err := s.RevertLast("./testdata/unit"); err == nil {
		t.Error("Expected error is returned on failed operation to load data")
	}
}

func TestSchema_Recreate_Unhappy_ReInitError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockApplier := schema_mock.NewMockApplier(ctrl)
	mockApplier.EXPECT().RevertScript("./testdata/unit/002.sql").Return(nil)
	mockApplier.EXPECT().ReInit().Return(errors.New("failed to reinit db"))

	mockScripter := schema_mock.NewMockScripter(ctrl)
	res := store.SchemaScriptCollection{&store.SchemaScript{
		ScriptName: "./testdata/unit/002.sql",
		Status:     store.StatusSuccess,
	}}
	mockScripter.EXPECT().GetAll().Times(1).Return(res, nil)
	mockScripter.EXPECT().Remove(gomock.Any()).Return(nil)

	s := schema.New(getMockDB(ctrl, true))
	s.Applier = mockApplier
	s.Scripter = mockScripter

	if err := s.Recreate("./testdata/unit", ""); err == nil {
		t.Error("Expected error is returned on failed recreate")
	}
}

func TestSchema_Upgrade_Unhappy_NotExistingPath(t *testing.T) {
	testCases := []struct {
		name string
		path string
	}{
		{
			name: "empty path - upgrade",
		},
		{
			name: "path not exists - upgrade",
			path: "not_existing_path",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockScripter := schema_mock.NewMockScripter(ctrl)
			mockScripter.EXPECT().GetAll().Times(1).Return(store.SchemaScriptCollection{}, nil)

			s := schema.New(getMockDB(ctrl, true))
			s.Scripter = mockScripter

			mockDB := getMockDB(ctrl, false)
			mockDB.EXPECT().Select(gomock.Any(), gomock.Any()).Times(1).Return(nil)

			s = schema.New(mockDB)
			s.Scripter = mockScripter

			if err := s.Upgrade(testCase.path, ""); err == nil {
				t.Errorf("Expected an error on call with not existing path")
			}
		})
	}
}

func TestSchema_RevertLast_Unhappy_NotExistingPath(t *testing.T) {
	testCases := []struct {
		name string
		path string
	}{
		{
			name: "empty path - revert",
		},
		{
			name: "path not exists - revert",
			path: "not_existing_path",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockScripter := schema_mock.NewMockScripter(ctrl)
			mockScripter.EXPECT().GetAll().Times(1).Return(store.SchemaScriptCollection{}, nil)

			s := schema.New(getMockDB(ctrl, true))
			s.Scripter = mockScripter

			if err := s.RevertLast(testCase.path); err == nil {
				t.Errorf("Expected an error on call with not existing path")
			}
		})
	}
}

func TestSchema_Recreate_Unhappy_NotExistingPath(t *testing.T) {
	testCases := []struct {
		name string
		path string
	}{
		{
			name: "empty path - recreate",
		},
		{
			name: "path not exists - recreate",
			path: "not_existing_path",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockScripter := schema_mock.NewMockScripter(ctrl)
			mockScripter.EXPECT().GetAll().Times(1).Return(store.SchemaScriptCollection{}, nil)

			s := schema.New(getMockDB(ctrl, true))
			s.Scripter = mockScripter

			if err := s.Recreate(testCase.path, ""); err == nil {
				t.Errorf("Expected an error on call with not existing path")
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

func TestSchema_Upgrage_Integration_Happy(t *testing.T) {
	if testing.Short() {
		t.Skip("skipped because of long running")
	}

	db, err := testdb.GetDB("./testdata/tmp/schema_execute_upgrade.db")
	if err != nil {
		t.Fatalf("failed to init database: %s", err)
	}
	defer testdb.ShutdownDB(db, t)

	s := schema.New(db)
	err = s.Upgrade("./testdata/upgrade/happy", "")
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	}

	data, err := s.Scripter.GetAll()
	if err != nil {
		t.Fatalf("not able get rows from table: %s", err)
	}

	expected := store.SchemaScriptCollection{
		&store.SchemaScript{
			ScriptName: "./testdata/upgrade/happy/001.sql",
			Status:     store.StatusSuccess,
		},
		&store.SchemaScript{
			ScriptName: "./testdata/upgrade/happy/002.sql",
			Status:     store.StatusSuccess,
		},
	}

	checkScriptTable("TestSchema_Upgrage_Integration_Happy", expected, data, t)
	checkTable("something", db, t, 0)
	checkTable("something_new", db, t, 0)
}

func TestSchema_Upgrade_Integration_Happy_TwoSteps(t *testing.T) {
	if testing.Short() {
		t.Skip("skipped because of long running")
	}

	db, err := testdb.GetDB("./testdata/tmp/schema_execute_upgrade_twosteps.db")
	if err != nil {
		t.Fatalf("failed to init database: %s", err)
	}
	defer func() {
		testdb.ShutdownDB(db, t)

		// cleanup
		if err = os.Remove("./testdata/upgrade/two_steps/001.sql"); err != nil {
			t.Fatalf("Cleanup files failed: %s", err)
		}

		if err = os.Remove("./testdata/upgrade/two_steps/002.sql"); err != nil {
			t.Fatalf("Cleanup files failed: %s", err)
		}
	}()

	s := schema.New(db)

	expected := step1(t, db, s)
	step2(t, db, s, expected)
}

func step1(t *testing.T, db *sqlx.DB, s schema.Schema) store.SchemaScriptCollection {
	var err error
	/**
	STEP 1
	*/
	if err = osutils.CopyFile("./testdata/upgrade/step1/001.sql", "./testdata/upgrade/two_steps/001.sql"); err != nil {
		t.Fatalf("failed to copy file: %s", err)
	}

	err = s.Upgrade("./testdata/upgrade/two_steps", "")
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	}

	data, err := s.Scripter.GetAll()
	if err != nil {
		t.Fatalf("not able get rows from table: %s", err)
	}

	expected := store.SchemaScriptCollection{
		&store.SchemaScript{
			ScriptName: "./testdata/upgrade/two_steps/001.sql",
			Status:     store.StatusSuccess,
		},
	}

	checkScriptTable("TestSchema_Upgrade_Integration_Happy_TwoSteps - Step 1", expected, data, t)
	checkTable("something", db, t, 0)
	return expected
}

func step2(t *testing.T, db *sqlx.DB, s schema.Schema, expected store.SchemaScriptCollection) {
	var err error
	/**
	STEP 2
	*/
	if err = osutils.CopyFile("./testdata/upgrade/step2/002.sql", "./testdata/upgrade/two_steps/002.sql"); err != nil {
		t.Fatalf("failed to copy file: %s", err)
	}

	err = s.Upgrade("./testdata/upgrade/two_steps", "")
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	}

	data, err := s.Scripter.GetAll()
	if err != nil {
		t.Fatalf("not able get rows from table: %s", err)
	}

	expected = append(expected, &store.SchemaScript{
		ScriptName: "./testdata/upgrade/two_steps/002.sql",
		Status:     store.StatusSuccess,
	})
	checkScriptTable("TestSchema_Upgrade_Integration_Happy_TwoSteps - Step 2", expected, data, t)
	checkTable("something", db, t, 0)
	checkTable("something_new", db, t, 0)
}

func TestSchema_RevertLast_Integration_Happy(t *testing.T) {
	if testing.Short() {
		t.Skip("skipped because of long running")
	}

	// prepare
	db, err := testdb.GetDB("./testdata/tmp/schema_execute_revert.db")
	if err != nil {
		t.Fatalf("failed to init database: %s", err)
	}
	defer testdb.ShutdownDB(db, t)

	s := schema.New(db)
	if err = s.Upgrade("./testdata/revert", ""); err != nil {
		t.Fatalf("Expected no error but got %s", err)
	}

	data, err := s.Scripter.GetAll()
	if err != nil {
		t.Fatalf("not able get rows from table: %s", err)
	}

	expected := store.SchemaScriptCollection{
		&store.SchemaScript{
			ScriptName: "./testdata/revert/001.sql",
			Status:     store.StatusSuccess,
		},
		&store.SchemaScript{
			ScriptName: "./testdata/revert/002.sql",
			Status:     store.StatusSuccess,
		},
	}

	testName := "TestSchema_Execute_Integration_CommandRevert_Happy"
	checkScriptTable(testName, expected, data, t)
	checkTable("something", db, t, 0)
	checkTable("something_new", db, t, 0)

	// now the test
	if err = s.RevertLast("./testdata/revert"); err != nil {
		t.Fatalf("not able to revert: %s", err)
	}

	data, err = s.Scripter.GetAll()
	if err != nil {
		t.Fatalf("not able get rows from table: %s", err)
	}

	expected = expected[:1]
	checkScriptTable(testName, expected, data, t)
	checkTable("something", db, t, 0)
}

func TestSchema_Recreate_Integration_Happy(t *testing.T) {
	if testing.Short() {
		t.Skip("skipped because of long running")
	}

	db, err := testdb.GetDB("./testdata/tmp/schema_execute_recreate.db")
	if err != nil {
		t.Fatalf("failed to init database: %s", err)
	}
	defer testdb.ShutdownDB(db, t)

	// prepare
	s := schema.New(db)
	if err = s.Upgrade("./testdata/recreate", ""); err != nil {
		t.Fatalf("Prepare: failed to create data: %s", err)
	}

	expected := store.SchemaScriptCollection{
		&store.SchemaScript{
			ScriptName: "./testdata/recreate/001.sql",
			Status:     store.StatusSuccess,
		},
		&store.SchemaScript{
			ScriptName: "./testdata/recreate/002.sql",
			Status:     store.StatusSuccess,
		},
		&store.SchemaScript{
			ScriptName: "./testdata/recreate/003_fake.sql",
			Status:     store.StatusSuccess,
		},
	}

	m := store.NewSchemaScriptMapper(db)
	if err = m.Add(expected[2]); err != nil {
		t.Fatalf("Prepare: couldn't add fake script: %s", err)
	}

	data, err := s.Scripter.GetAll()
	if err != nil {
		t.Fatalf("Prepare: not able get rows from table: %s", err)
	}

	q := `INSERT INTO something (id) VALUES (1)`
	if _, err = db.Exec(q); err != nil {
		t.Fatalf("Prepare: failed to add data to table: %s", err)
	}

	testName := "TestSchema_Execute_Integration_CommandRecreate_Happy"
	checkScriptTable(testName, expected, data, t)
	checkTable("something", db, t, 1)
	checkTable("something_new", db, t, 0)

	// now the test
	if err = s.Recreate("./testdata/recreate", ""); err != nil {
		t.Fatalf("not able to recreate: %s", err)
	}

	data, err = s.Scripter.GetAll()
	if err != nil {
		t.Fatalf("not able get rows from table: %s", err)
	}

	expected = expected[0:2]
	checkScriptTable(testName, expected, data, t)
	checkTable("something", db, t, 0)
	checkTable("something_new", db, t, 0)
}

func checkScriptTable(testName string, expected store.SchemaScriptCollection, actual store.SchemaScriptCollection, t *testing.T) {
	if len(expected) != len(actual) {
		t.Fatalf("%s: Expeted %d rows in table but got %d", testName, len(expected), len(actual))
	}

	for i, v := range expected {
		w := actual[i]
		if v.ScriptName != w.ScriptName {
			t.Errorf("%s: Expected script name %s but got %s", testName, v.ScriptName, w.ScriptName)
		}

		if v.Status != w.Status {
			t.Errorf("%s: Expected status %s but got %s", testName, v.Status, w.Status)
		}

		if v.ErrorMsg != w.ErrorMsg {
			t.Errorf("%s Expected error message %s but got %s", testName, v.ErrorMsg, w.ErrorMsg)
		}

		if v.AppVersion != w.AppVersion {
			t.Errorf("%s: Expected app version %s but got %s", testName, v.AppVersion, w.AppVersion)
		}
	}
}

func checkTable(tableName string, db *sqlx.DB, t *testing.T, expectedCount uint32) {
	var counter []uint32
	q := db.Rebind(fmt.Sprintf("SELECT count(id) FROM %s;", tableName))
	err := db.Select(&counter, q)
	if err != nil {
		t.Fatalf("not able count rows in table: %s", err)
	}

	if len(counter) == 0 || counter[0] != expectedCount {
		t.Errorf("expected number of %d rows in table '%s' but got %d", expectedCount, tableName, counter[0])
	}
}
