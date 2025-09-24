package util

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

// CreateDirIfNotExists checks if directory exists, if not it creates it
func CreateDirIfNotExists(path string, perm fs.FileMode) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(path, perm); err != nil {
			panic(fmt.Errorf("failed to create database directory"))
		}
	}
}
