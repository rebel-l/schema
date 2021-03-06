package sqlfile_test

import (
	"testing"

	"github.com/rebel-l/go-utils/array"
	"github.com/rebel-l/schema/sqlfile"
)

func TestScanHappy(t *testing.T) {
	testCases := []struct {
		name     string
		path     string
		expected []string
	}{
		{
			name: "case1 with different files",
			path: "./testdata/case1",
			expected: []string{
				"./testdata/case1/001_with_content.sql",
				"./testdata/case1/004_some_more.sql",
			},
		},
		{
			name:     "case2 with no files",
			path:     "./testdata/case2",
			expected: []string{},
		},
	}

	for _, testCase := range testCases {
		path := testCase.path
		expected := testCase.expected
		t.Run(testCase.name, func(t *testing.T) {
			actual, err := sqlfile.Scan(path)
			if err != nil {
				t.Fatalf("scan shouldn't cause error: %s", err)
			}

			if !array.StringArrayEquals(expected, actual) {
				t.Errorf("Expected %#v but got %#v", expected, actual)
			}
		})
	}
}

func TestScanUnhappy(t *testing.T) {
	_, err := sqlfile.Scan("")
	if err == nil {
		t.Error("Scan empty path should cause an error")
	}
}

func TestScanReverseHappy(t *testing.T) {
	testCases := []struct {
		name     string
		path     string
		expected []string
	}{
		{
			name: "case1 with different files",
			path: "./testdata/case1",
			expected: []string{
				"./testdata/case1/004_some_more.sql",
				"./testdata/case1/001_with_content.sql",
			},
		},
		{
			name:     "case2 with no files",
			path:     "./testdata/case2",
			expected: []string{},
		},
	}

	for _, testCase := range testCases {
		path := testCase.path
		expected := testCase.expected
		t.Run(testCase.name, func(t *testing.T) {
			actual, err := sqlfile.ScanReverse(path)
			if err != nil {
				t.Fatalf("scan shouldn't cause error: %s", err)
			}

			if !array.StringArrayEquals(expected, actual) {
				t.Errorf("Expected %#v but got %#v", expected, actual)
			}
		})
	}
}

func TestScanReverseUnhappy(t *testing.T) {
	_, err := sqlfile.ScanReverse("")
	if err == nil {
		t.Error("Scan empty path should cause an error")
	}
}

func TestReadHappy(t *testing.T) {
	testCases := []struct {
		command  string
		expected string
	}{
		{
			command: sqlfile.CommandUpgrade,
			expected: `
CREATE TABLE IF NOT EXISTS test (
id INTEGER
);
CREATE TABLE IF NOT EXISTS another (
id INTEGER
);`,
		},
		{
			command: sqlfile.CommandDowngrade,
			expected: `
DROP TABLE IF EXISTS test;
DROP TABLE IF EXISTS another;`,
		},
	}

	for _, testCase := range testCases {
		command := testCase.command
		expected := testCase.expected
		t.Run(testCase.command, func(t *testing.T) {
			fileName := "./testdata/Read/test.sql"
			actual, err := sqlfile.Read(fileName, command)
			if err != nil {
				t.Errorf("Expected that file name %s is readable but got %s", fileName, err)
			}

			if expected != actual {
				t.Errorf("Expected file content '%s' but got '%s'", expected, actual)
			}
		})
	}
}

func TestReadUnhappy(t *testing.T) {
	content, err := sqlfile.Read("not_exist.sql", sqlfile.CommandUpgrade)
	if err == nil {
		t.Error("Expected that error is thrown for not existing file")
	}

	if content != "" {
		t.Errorf("Expected that content is empty on error but got %s", content)
	}
}
