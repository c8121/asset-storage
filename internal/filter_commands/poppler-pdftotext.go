package filter_commands

import (
	"github.com/c8121/asset-storage/internal/util"
)

var (
	PopplerPdfBinPaths = []string{
		"/usr/bin/pdftotext",
	}

	PopplerPdfBinPath = ""
)

// FindPopplerPdfBin checks if one of PopplerPdfBinPaths exists
func FindPopplerPdfBin() string {

	if PopplerPdfBinPath != "" {
		return PopplerPdfBinPath
	}
	PopplerPdfBinPath = util.FindFile(PopplerPdfBinPaths)
	return PopplerPdfBinPath
}
