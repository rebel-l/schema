package store

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rebel-l/schema/mocks"
)

func TestSchemaVersionMapper_Add_Happy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedID := int64(101)
	version := NewSchemaVersionSuccess("my_sql_script.sql")

	mockRes := mocks.NewMockResult(ctrl)
	mockRes.EXPECT().LastInsertId().Return(expectedID, nil)

	mockDb := mocks.NewMockDatabaseConnector(ctrl)
	mockDb.EXPECT().
		Exec(gomock.Any(), version.ScriptName, version.ExecutedAt.Format(dateTimeFormat), version.Status, version.ErrorMsg).
		Return(mockRes, nil)

	mapper := NewSchemaVersionMapper(mockDb)
	if err := mapper.Add(version); err != nil {
		t.Errorf("error is not expected but got: %s", err)
	}

	if version.ID != expectedID {
		t.Errorf("id was not set to entry, expected %d but got %d", expectedID, version.ID)
	}
}

func TestSchemaVersionMapper_Add_Unhappy_NilEntry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDb := mocks.NewMockDatabaseConnector(ctrl)
	mockDb.EXPECT().Exec(gomock.Any()).Times(0)
	mapper := NewSchemaVersionMapper(mockDb)
	if err := mapper.Add(nil); err == nil {
		t.Errorf("nil should be not allowed and throw an error")
	}
}

func TestSchemaVersionMapper_Add_Unhappy_InsertError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	version := NewSchemaVersionSuccess("my_sql_script.sql")

	mockRes := mocks.NewMockResult(ctrl)
	mockRes.EXPECT().LastInsertId().Times(0)

	mockDb := mocks.NewMockDatabaseConnector(ctrl)
	mockDb.EXPECT().
		Exec(gomock.Any(), version.ScriptName, version.ExecutedAt.Format(dateTimeFormat), version.Status, version.ErrorMsg).
		Return(mockRes, errors.New("insert failed"))

	mapper := NewSchemaVersionMapper(mockDb)
	if err := mapper.Add(version); err == nil {
		t.Errorf("error is expected on failing insert")
	}
}

func TestSchemaVersionMapper_Add_Unhappy_LastInsertError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	version := NewSchemaVersionSuccess("my_sql_script.sql")

	mockRes := mocks.NewMockResult(ctrl)
	mockRes.EXPECT().LastInsertId().Return(int64(0), errors.New("last insert failed"))

	mockDb := mocks.NewMockDatabaseConnector(ctrl)
	mockDb.EXPECT().
		Exec(gomock.Any(), version.ScriptName, version.ExecutedAt.Format(dateTimeFormat), version.Status, version.ErrorMsg).
		Return(mockRes, nil)

	mapper := NewSchemaVersionMapper(mockDb)
	if err := mapper.Add(version); err == nil {
		t.Errorf("error is expected on failing insert")
	}
}

func TestSchemaVersionMapper_GetByID_Happy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := int64(203)

	mockDb := mocks.NewMockDatabaseConnector(ctrl)
	mockDb.EXPECT().Select(gomock.Any(), gomock.Any(), id).Return(nil)

	mapper := NewSchemaVersionMapper(mockDb)
	res, err := mapper.GetByID(id)
	if err != nil {
		t.Errorf("no error expected on reading")
	}

	if res == nil {
		t.Errorf("returned schema version should not be nil")
	}
}

func TestSchemaVersionMapper_GetByID_Unhappy_ID(t *testing.T) {
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

			mapper := NewSchemaVersionMapper(mockDb)
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

func TestSchemaVersionMapper_GetByID_Unhappy_SelectError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := int64(666)

	mockDb := mocks.NewMockDatabaseConnector(ctrl)
	mockDb.EXPECT().Select(gomock.Any(), gomock.Any(), id).Return(errors.New("select failed"))

	mapper := NewSchemaVersionMapper(mockDb)
	res, err := mapper.GetByID(id)
	if err == nil {
		t.Errorf("expected error on reading failure")
	}

	if res != nil {
		t.Errorf("returned schema version should be nil on error")
	}
}
