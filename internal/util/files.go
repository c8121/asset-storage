package util

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// CreateDirIfNotExists checks if directory exists, if not it creates it
func CreateDirIfNotExists(path string, perm fs.FileMode) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(path, perm); err != nil {
			panic(fmt.Errorf("failed to create database directory"))
		} else {
			fmt.Printf("Created directory: %s\n", path)
		}
	}
}

// FindFile takes path that might contain wildcards and return the first that exists
func FindFile(paths []string) string {

	for _, path := range paths {
		l, err := filepath.Glob(path)
		if err != nil {
			continue
		}
		if len(l) > 0 {
			return l[0]
		}
	}

	fmt.Printf("No files found in %s\n", paths)
	return ""
}

func SplitPath(path string) []string {

	items := make([]string, 0)

	s := 0
	p := 0
	for p < len(path) {
		if path[p] == '/' || path[p] == '\\' {
			items = append(items, path[s:p])
			s = p + 1
		}
		p++
	}
	if p > 0 {
		items = append(items, path[s:p])
	}

	return items
}
