package filter_commands

import (
	"github.com/c8121/asset-storage/internal/util"
)

var (
	TesseractBinPaths = []string{
		"/usr/bin/tesseract",
		"C:/Program Files/Tesseract*/tesseract.exe",
	}

	TesseractBinPath = ""
)

// FindTesseractBin checks if one of TesseractBinPaths exists
func FindTesseractBin() string {

	if TesseractBinPath != "" {
		return TesseractBinPath
	}
	TesseractBinPath = util.FindFile(TesseractBinPaths)
	return TesseractBinPath
}
