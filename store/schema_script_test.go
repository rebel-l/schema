package store_test

import (
	"testing"
	"time"

	"github.com/rebel-l/schema/store"
)

func TestNewSchemaScriptSuccess(t *testing.T) {
	expected := store.SchemaScript{
		ScriptName: "something.sql",
		ExecutedAt: time.Now(),
		Status:     store.StatusSuccess,
		AppVersion: "0.1.3",
	}

	actual := store.NewSchemaScriptSuccess(expected.ScriptName, expected.AppVersion)

	if actual.ID > 0 {
		t.Errorf("expected id to be 0 but got %d", actual.ID)
	}

	if actual.ScriptName != expected.ScriptName {
		t.Errorf("expected scriptname '%s' but got '%s'", expected.ScriptName, actual.ScriptName)
	}

	if actual.ExecutedAt.Before(expected.ExecutedAt) {
		t.Errorf(
			"expected executionAt would be later or equal than '%s' but got '%s'",
			expected.ExecutedAt.String(),
			actual.ExecutedAt.String(),
		)
	}

	if actual.Status != expected.Status {
		t.Errorf("expected that status is automatically set to '%s' but got '%s'", expected.Status, actual.Status)
	}

	if actual.AppVersion != expected.AppVersion {
		t.Errorf(
			"expected that appVersion is automatically set to '%s' but got '%s'",
			expected.AppVersion,
			actual.AppVersion,
		)
	}

	if actual.ErrorMsg != "" {
		t.Errorf("expected error mesage to be empty but got '%s'", actual.ErrorMsg)
	}
}

func TestNewSchemaScriptError(t *testing.T) {
	expected := store.SchemaScript{
		ScriptName: "failing.sql",
		ExecutedAt: time.Now(),
		Status:     store.StatusError,
		ErrorMsg:   "houston we have a problem",
	}

	actual := store.NewSchemaScriptError(expected.ScriptName, expected.AppVersion, expected.ErrorMsg)

	if actual.ID > 0 {
		t.Errorf("expected id to be 0 but got %d", actual.ID)
	}

	if actual.ScriptName != expected.ScriptName {
		t.Errorf("expected scriptname '%s' but got '%s'", expected.ScriptName, actual.ScriptName)
	}

	if actual.ExecutedAt.Before(expected.ExecutedAt) {
		t.Errorf(
			"expected executionAt would be later or equal than '%s' but got '%s'",
			expected.ExecutedAt.String(),
			actual.ExecutedAt.String(),
		)
	}

	if actual.Status != expected.Status {
		t.Errorf("expected that satus is automatically set to '%s' but got '%s'", expected.Status, actual.Status)
	}

	if actual.AppVersion != expected.AppVersion {
		t.Errorf(
			"expected that appVersion is automatically set to '%s' but got '%s'",
			expected.AppVersion,
			actual.AppVersion,
		)
	}

	if actual.ErrorMsg != expected.ErrorMsg {
		t.Errorf("expected error mesage to be '%s' but got '%s'", expected.ErrorMsg, actual.ErrorMsg)
	}
}

func TestSchemaScriptCollection_ScriptExecuted(t *testing.T) {
	testCases := []struct {
		name       string
		scriptName string
		collection store.SchemaScriptCollection
		expected   bool
	}{
		{
			name:       "empty collection",
			scriptName: "something.sql",
			expected:   false,
		},
		{
			name:       "one item in collection, no hit",
			scriptName: "something.sql",
			collection: store.SchemaScriptCollection{
				&store.SchemaScript{ScriptName: "else.sql", Status: store.StatusSuccess},
			},
			expected: false,
		},
		{
			name:       "two items in collection, no hit",
			scriptName: "something.sql",
			collection: store.SchemaScriptCollection{
				&store.SchemaScript{ScriptName: "else.sql", Status: store.StatusSuccess},
				&store.SchemaScript{ScriptName: "else2.sql", Status: store.StatusSuccess},
			},
			expected: false,
		},
		{
			name: "empty script name",
			collection: store.SchemaScriptCollection{
				&store.SchemaScript{ScriptName: "else.sql", Status: store.StatusSuccess},
				&store.SchemaScript{ScriptName: "else2.sql", Status: store.StatusSuccess},
			},
			expected: false,
		},
		{
			name:       "one item in collection, hit",
			scriptName: "hit.sql",
			collection: store.SchemaScriptCollection{
				&store.SchemaScript{ScriptName: "hit.sql", Status: store.StatusSuccess},
			},
			expected: true,
		},
		{
			name:       "two items in collection, hit",
			scriptName: "hit.sql",
			collection: store.SchemaScriptCollection{
				&store.SchemaScript{ScriptName: "hit1.sql", Status: store.StatusSuccess},
				&store.SchemaScript{ScriptName: "hit.sql", Status: store.StatusSuccess},
			},
			expected: true,
		},
		{
			name:       "one error item in collection",
			scriptName: "hit.sql",
			collection: store.SchemaScriptCollection{
				&store.SchemaScript{ScriptName: "hit.sql", Status: store.StatusError},
			},
			expected: false,
		},
		{
			name:       "one error item, one success item in collection",
			scriptName: "hit.sql",
			collection: store.SchemaScriptCollection{
				&store.SchemaScript{ScriptName: "hit1.sql", Status: store.StatusSuccess},
				&store.SchemaScript{ScriptName: "hit.sql", Status: store.StatusError},
			},
			expected: false,
		},
	}

	for _, testCase := range testCases {
		collection := testCase.collection
		scriptName := testCase.scriptName
		expected := testCase.expected
		t.Run(testCase.name, func(t *testing.T) {
			actual := collection.ScriptExecuted(scriptName)
			if expected != actual {
				t.Errorf(
					"Expected for script '%s' and collection %v result is %t but got %t",
					scriptName,
					collection,
					expected,
					actual,
				)
			}
		})
	}
}
