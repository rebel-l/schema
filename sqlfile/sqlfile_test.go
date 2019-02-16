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
			path: "./../tests/data/sqlfile/case1",
			expected: []string{
				"./../tests/data/sqlfile/case1/001_with_content.sql",
				"./../tests/data/sqlfile/case1/004_some_more.sql",
			},
		},
		{
			name:     "case2 with no files",
			path:     "./../tests/data/sqlfile/case2",
			expected: []string{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual, err := sqlfile.Scan(testCase.path)
			if err != nil {
				t.Fatalf("scan shouldn't cause error: %s", err)
			}

			if !array.StringArrayEquals(testCase.expected, actual) {
				t.Errorf("Expected %#v but got %#v", testCase.expected, actual)
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
