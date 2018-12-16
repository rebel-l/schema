package store

import (
	"testing"
	"time"
)

func TestNewSchemaVersionSuccess(t *testing.T) {
	expected := SchemaVersion{
		ScriptName: "something.sql",
		ExecutedAt: time.Now(),
		Status:     StatusSuccess,
	}

	actual := NewSchemaVersionSuccess(expected.ScriptName)

	if actual.ID > 0 {
		t.Errorf("expected id to be 0 but got %d", actual.ID)
	}

	if actual.ScriptName != expected.ScriptName {
		t.Errorf("expected scriptname '%s' but got '%s'", expected.ScriptName, actual.ScriptName)
	}

	if actual.ExecutedAt.Before(expected.ExecutedAt) {
		t.Errorf("expected executionAt would be later or equal than '%s' but got '%s'", expected.ExecutedAt.String(), actual.ExecutedAt.String())
	}

	if actual.Status != expected.Status {
		t.Errorf("expected that satus is automatically set to '%s' but got '%s'", expected.Status, actual.Status)
	}

	if actual.ErrorMsg != "" {
		t.Errorf("expected error mesage to be empty but got '%s'", actual.ErrorMsg)
	}
}

func TestNewSchemaVersionError(t *testing.T) {
	expected := SchemaVersion{
		ScriptName: "failing.sql",
		ExecutedAt: time.Now(),
		Status:     StatusError,
		ErrorMsg:   "houston we have a problem",
	}

	actual := NewSchemaVersionError(expected.ScriptName, expected.ErrorMsg)

	if actual.ID > 0 {
		t.Errorf("expected id to be 0 but got %d", actual.ID)
	}

	if actual.ScriptName != expected.ScriptName {
		t.Errorf("expected scriptname '%s' but got '%s'", expected.ScriptName, actual.ScriptName)
	}

	if actual.ExecutedAt.Before(expected.ExecutedAt) {
		t.Errorf("expected executionAt would be later or equal than '%s' but got '%s'", expected.ExecutedAt.String(), actual.ExecutedAt.String())
	}

	if actual.Status != expected.Status {
		t.Errorf("expected that satus is automatically set to '%s' but got '%s'", expected.Status, actual.Status)
	}

	if actual.ErrorMsg != expected.ErrorMsg {
		t.Errorf("expected error mesage to be '%s' but got '%s'", expected.ErrorMsg, actual.ErrorMsg)
	}
}
