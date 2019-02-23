package schema

import (
	"errors"
	"testing"

	"github.com/rebel-l/schema/mocks/logrus_mock"
	"github.com/rebel-l/schema/mocks/schema_mock"
	"github.com/rebel-l/schema/mocks/store_mock"
	"github.com/rebel-l/schema/store"

	"github.com/golang/mock/gomock"
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

	schema := Schema{
		Logger:   getMockLogger(ctrl, true),
		db:       mockDB,
		applier:  mockApplier,
		scripter: mockScripter,
	}

	if err := schema.Execute("./tests/data/schema/unit", CommandUpgrade, ""); err != nil {
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

	schema := Schema{
		Logger:   getMockLogger(ctrl, true),
		db:       mockDB,
		applier:  mockApplier,
		scripter: mockScripter,
	}

	if err := schema.Execute("./tests/data/schema/unit", CommandUpgrade, ""); err == nil {
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

	schema := Schema{
		Logger:   mockLogger,
		applier:  mockApplier,
		scripter: mockScripter,
		db:       mockDB,
	}

	err := schema.Execute("./tests/data/schema/unit", CommandUpgrade, "")
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

	schema := Schema{
		Logger:   mockLogger,
		applier:  mockApplier,
		scripter: mockScripter,
		db:       mockDB,
	}

	err := schema.Execute("./tests/data/schema/unit", CommandUpgrade, "")
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

	schema := Schema{
		Logger:   mockLogger,
		applier:  mockApplier,
		scripter: mockScripter,
		db:       mockDB,
	}

	err := schema.Execute("./tests/data/schema/unit", CommandUpgrade, "")
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

	schema := Schema{
		Logger:   getMockLogger(ctrl, true),
		applier:  mockApplier,
		scripter: mockScripter,
	}

	if err := schema.Execute("./tests/data/schema/unit", CommandRevert, ""); err == nil {
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

	schema := Schema{
		Logger:   getMockLogger(ctrl, true),
		applier:  mockApplier,
		scripter: mockScripter,
	}

	if err := schema.Execute("./tests/data/schema/unit", CommandRevert, ""); err == nil {
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
			command: CommandUpgrade,
		},
		{
			name:    "path not exists - upgrade",
			path:    "not_existing_path",
			command: CommandUpgrade,
		},
		{
			name:    "empty path - revert",
			command: CommandRevert,
		},
		{
			name:    "path not exists - revert",
			path:    "not_existing_path",
			command: CommandRevert,
		},
		{
			name:    "empty path - recreate",
			command: CommandRecreate,
		},
		{
			name:    "path not exists - recreate",
			path:    "not_existing_path",
			command: CommandRecreate,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockScripter := schema_mock.NewMockScripter(ctrl)
			mockScripter.EXPECT().GetAll().Times(1).Return(store.SchemaScriptCollection{}, nil)

			schema := Schema{
				Logger:   getMockLogger(ctrl, true),
				scripter: mockScripter,
			}
			if testCase.command == CommandUpgrade {
				mockDB := getMockDB(ctrl, false)
				mockDB.EXPECT().Select(gomock.Any(), gomock.Any()).Times(1).Return(nil)
				schema.db = mockDB
			}

			if err := schema.Execute(testCase.path, testCase.command, ""); err == nil {
				t.Errorf("Expected an error on call with not existing path")
			}
		})
	}
}

func TestSchema_Execute_Unhappy_WrongCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	schema := getSchemaWithDummies(ctrl)

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
			if err := schema.Execute("./", testCase.command, ""); err == nil {
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

func getSchemaWithDummies(ctrl *gomock.Controller) Schema {
	return Schema{
		Logger: getMockLogger(ctrl, true),
		db:     getMockDB(ctrl, true),
	}
}
