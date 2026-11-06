package util

import (
	"errors"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
)

// CreateDirIfNotExists checks if directory exists, if not it creates it
func CreateDirIfNotExists(path string, perm fs.FileMode) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(path, perm); err != nil {
			slog.Error("failed to create database directory", "error", err)
			return err
		} else {
			slog.Info("Created directory", "path", path)
		}
	}
	return nil
}

// FindFile takes path that might contain wildcards and return the first that exists.
// Returns an empty string "", if no existing path found
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

	slog.Info("No files found", "paths", paths)
	return ""
}
