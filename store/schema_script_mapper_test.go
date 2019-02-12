package store

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rebel-l/schema/mocks"
)

func TestSchemaScriptMapper_Add_Happy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedID := int64(101)
	script := NewSchemaScriptSuccess("my_sql_script.sql", "0.1.0")

	mockRes := mocks.NewMockResult(ctrl)
	mockRes.EXPECT().LastInsertId().Return(expectedID, nil)

	mockDb := mocks.NewMockDatabaseConnector(ctrl)
	mockDb.EXPECT().
		Exec(gomock.Any(), script.ScriptName, script.ExecutedAt.Format(dateTimeFormat), script.Status, script.ErrorMsg, script.AppVersion).
		Return(mockRes, nil)

	mapper := NewSchemaScriptMapper(mockDb)
	if err := mapper.Add(script); err != nil {
		t.Errorf("error is not expected but got: %s", err)
	}

	if script.ID != expectedID {
		t.Errorf("id was not set to entry, expected %d but got %d", expectedID, script.ID)
	}
}

func TestSchemaScriptMapper_Add_Unhappy_NilEntry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDb := mocks.NewMockDatabaseConnector(ctrl)
	mockDb.EXPECT().Exec(gomock.Any()).Times(0)
	mapper := NewSchemaScriptMapper(mockDb)
	if err := mapper.Add(nil); err == nil {
		t.Errorf("nil should be not allowed and throw an error")
	}
}

func TestSchemaScriptMapper_Add_Unhappy_InsertError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	script := NewSchemaScriptSuccess("my_sql_script.sql", "0.2.0")

	mockRes := mocks.NewMockResult(ctrl)
	mockRes.EXPECT().LastInsertId().Times(0)

	mockDb := mocks.NewMockDatabaseConnector(ctrl)
	mockDb.EXPECT().
		Exec(gomock.Any(), script.ScriptName, script.ExecutedAt.Format(dateTimeFormat), script.Status, script.ErrorMsg, script.AppVersion).
		Return(mockRes, errors.New("insert failed"))

	mapper := NewSchemaScriptMapper(mockDb)
	if err := mapper.Add(script); err == nil {
		t.Errorf("error is expected on failing insert")
	}
}

func TestSchemaScriptMapper_Add_Unhappy_LastInsertError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	script := NewSchemaScriptSuccess("my_sql_script.sql", "")

	mockRes := mocks.NewMockResult(ctrl)
	mockRes.EXPECT().LastInsertId().Return(int64(0), errors.New("last insert failed"))

	mockDb := mocks.NewMockDatabaseConnector(ctrl)
	mockDb.EXPECT().
		Exec(gomock.Any(), script.ScriptName, script.ExecutedAt.Format(dateTimeFormat), script.Status, script.ErrorMsg, script.AppVersion).
		Return(mockRes, nil)

	mapper := NewSchemaScriptMapper(mockDb)
	if err := mapper.Add(script); err == nil {
		t.Errorf("error is expected on failing insert")
	}
}

func TestSchemaScriptMapper_GetByID_Happy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := int64(203)

	mockDb := mocks.NewMockDatabaseConnector(ctrl)
	mockDb.EXPECT().Get(gomock.Any(), gomock.Any(), id).Return(nil)

	mapper := NewSchemaScriptMapper(mockDb)
	_, err := mapper.GetByID(id)
	if err != nil {
		t.Errorf("no error expected on reading")
	}
}

func TestSchemaScriptMapper_GetByID_Unhappy_ID(t *testing.T) {
	tcs := []struct {
		name string
		id   int64
	}{
		{
			name: "id is zero",
			id:   int64(0),
		},
		{
			name: "id is negative",
			id:   int64(-1),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockDb := mocks.NewMockDatabaseConnector(ctrl)
			mockDb.EXPECT().Select(gomock.Any(), gomock.Any()).Times(0)

			mapper := NewSchemaScriptMapper(mockDb)
			res, err := mapper.GetByID(tc.id)
			if err == nil {
				t.Errorf("get entry with zero or negative id should cause an error")
			}

			if res != nil {
				t.Errorf("in case of an error no schema version should be returned")
			}
			ctrl.Finish()
		})
	}
}

func TestSchemaScriptMapper_GetByID_Unhappy_SelectError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := int64(666)

	mockDb := mocks.NewMockDatabaseConnector(ctrl)
	mockDb.EXPECT().Get(gomock.Any(), gomock.Any(), id).Return(errors.New("select failed"))

	mapper := NewSchemaScriptMapper(mockDb)
	res, err := mapper.GetByID(id)
	if err == nil {
		t.Errorf("expected error on reading failure")
	}

	if res != nil {
		t.Errorf("returned schema version should be nil on error")
	}
}

func TestSchemaScriptMapper_GetAll_Happy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDb := mocks.NewMockDatabaseConnector(ctrl)
	empty := make([]interface{}, 0)
	mockDb.EXPECT().Select(gomock.Any(), gomock.Any(), gomock.Eq(empty))

	mapper := NewSchemaScriptMapper(mockDb)
	_, err := mapper.GetAll()
	if err != nil {
		t.Errorf("no error expected on getting all entries")
	}
}

func TestSchemaScriptMapper_GetAll_Unhappy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDb := mocks.NewMockDatabaseConnector(ctrl)
	empty := make([]interface{}, 0)
	mockDb.EXPECT().Select(gomock.Any(), gomock.Any(), gomock.Eq(empty)).Return(errors.New("select failed"))

	mapper := NewSchemaScriptMapper(mockDb)
	res, err := mapper.GetAll()
	if err == nil {
		t.Errorf("expected error on reading failure")
	}

	if res != nil {
		t.Errorf("returned list of schema versions should be nil on error")
	}
}
