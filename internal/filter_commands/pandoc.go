package filter_commands

import (
	"github.com/c8121/asset-storage/internal/util"
)

var (
	PandocBinPaths = []string{
		"/usr/bin/pandoc",
	}

	PandocBinPath = ""
)

// FindPandocBin checks if one of PandocBinPaths exists
func FindPandocBin() string {

	if PandocBinPath != "" {
		return PandocBinPath
	}
	PandocBinPath = util.FindFile(PandocBinPaths)
	return PandocBinPath
}
