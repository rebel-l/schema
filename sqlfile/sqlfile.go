package sqlfile

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

// Scan returns a sorted (asc) list of files (including path) ending with .sql
// and file size bigger than zero. It excludes directories.
func Scan(path string) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("scan files failed: %s", err)
	}

	cleaned := make([]string, 0)
	for _, v := range files {
		if v.IsDir() {
			continue
		}

		if filepath.Ext(v.Name()) != ".sql" {
			continue
		}

		if v.Size() == 0 {
			continue
		}

		cleaned = append(cleaned, path+"/"+v.Name())
	}
	return cleaned, nil
}
