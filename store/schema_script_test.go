package store

import (
	"testing"
	"time"
)

func TestNewSchemaScriptSuccess(t *testing.T) {
	expected := SchemaScript{
		ScriptName: "something.sql",
		ExecutedAt: time.Now(),
		Status:     StatusSuccess,
		AppVersion: "0.1.3",
	}

	actual := NewSchemaScriptSuccess(expected.ScriptName, expected.AppVersion)

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
		t.Errorf("expected that status is automatically set to '%s' but got '%s'", expected.Status, actual.Status)
	}

	if actual.AppVersion != expected.AppVersion {
		t.Errorf("expected that appVersion is automatically set to '%s' but got '%s'", expected.AppVersion, actual.AppVersion)
	}

	if actual.ErrorMsg != "" {
		t.Errorf("expected error mesage to be empty but got '%s'", actual.ErrorMsg)
	}
}

func TestNewSchemaScriptError(t *testing.T) {
	expected := SchemaScript{
		ScriptName: "failing.sql",
		ExecutedAt: time.Now(),
		Status:     StatusError,
		ErrorMsg:   "houston we have a problem",
	}

	actual := NewSchemaScriptError(expected.ScriptName, expected.AppVersion, expected.ErrorMsg)

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

	if actual.AppVersion != expected.AppVersion {
		t.Errorf("expected that appVersion is automatically set to '%s' but got '%s'", expected.AppVersion, actual.AppVersion)
	}

	if actual.ErrorMsg != expected.ErrorMsg {
		t.Errorf("expected error mesage to be '%s' but got '%s'", expected.ErrorMsg, actual.ErrorMsg)
	}
}
