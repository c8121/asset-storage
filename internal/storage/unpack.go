package storage

import (
	"archive/zip"
	"fmt"
	"strings"

	"github.com/c8121/asset-storage/internal/util"
)

func IsUnpackable(path string, mimeType string) bool {
	if strings.HasSuffix(strings.ToLower(mimeType), "zip") {
		return true
	}
	return false
}

// Unpack deflates files directly to storage
func Unpack(path string, mimeType string) ([]AddedFileInfo, error) {

	if !strings.HasSuffix(strings.ToLower(mimeType), "zip") {
		fmt.Printf("Not an archive: %s, %s", path, mimeType)
		return nil, nil
	}

	unpacked := make([]AddedFileInfo, 0)

	reader, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer util.CloseOrLog(reader)

	for _, file := range reader.File {

		if file.FileInfo().IsDir() {
			continue
		}

		reader, err := file.Open()
		if err != nil {
			fmt.Printf("Error opening file %s: %s\n", file.Name, err)
			continue
		}

		writer, err := newTempWriter(file.FileInfo().Size())
		if err != nil {
			fmt.Printf("Error creating temp writer: %s\n", err)
			continue
		}

		info, err := moveToStorage(reader, writer)
		if err != nil {
			fmt.Printf("Error copying file %s: %s\n", file.Name, err)
		} else {
			info.SourcePath = file.Name
			unpacked = append(unpacked, *info)
		}

		util.CloseOrLog(writer)
		util.CloseOrLog(reader)
	}

	return unpacked, nil
}
