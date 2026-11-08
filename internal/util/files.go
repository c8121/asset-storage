package util

import (
	"log/slog"
	"path/filepath"
)

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
