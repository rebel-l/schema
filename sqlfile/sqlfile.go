package sqlfile

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
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

// ScanReverse does the same as Scan but returns the filenames in reverse order
func ScanReverse(path string) ([]string, error) {
	var files sort.StringSlice
	var err error
	files, err = Scan(path)
	if err != nil {
		return nil, err
	}
	sort.Sort(sort.Reverse(files))
	return files, nil
}

// Read returns the content of a file
func Read(fileName string) (string, error) {
	file, err := os.Open(filepath.Clean(fileName))

	if err != nil {
		return "", err
	}
	defer func() {
		err = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	var buffer string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) > 0 {
			buffer += "\n" + line
		}
	}
	return buffer, nil
}
