package ssh_server

import (
	"path/filepath"
	"strings"
)

func resolve(rootDirectory string, requestedPath string) (string, error) {

	//fmt.Printf(" - Path requested: '%s'\n", path)
	if strings.Contains(requestedPath, ":") {
		return "", ErrorInvalidPathRequest
	}

	resolved := requestedPath
	if resolved == string(filepath.Separator) || resolved == "/" {
		resolved = ""
	}

	resolved = filepath.FromSlash(resolved)
	resolved = strings.TrimPrefix(resolved, string(filepath.Separator))
	resolved = strings.TrimPrefix(resolved, "/")
	resolved = filepath.Join(rootDirectory, resolved)

	//fmt.Printf(" - Path resolved: '%s'\n", resolved)
	return resolved, nil
}
