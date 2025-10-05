package storage

import (
	"archive/zip"
	"fmt"
	"io"
	"strings"

	"github.com/c8121/asset-storage/internal/util"
)

type Unpacked struct {
	Source      string //Original path (zip-file for example)
	ArchivePath string //Path of file in archive
	TempPath    string //Unpacked temp file
}

func IsUnpackable(path string, mimeType string) bool {
	if strings.HasSuffix(strings.ToLower(mimeType), "zip") {
		return true
	}
	return false
}

func Unpack(path string, mimeType string) ([]Unpacked, error) {

	if !strings.HasSuffix(strings.ToLower(mimeType), "zip") {
		fmt.Printf("Not an archive: %s, %s", path, mimeType)
		return nil, nil
	}

	unpacked := make([]Unpacked, 0)

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

		writer, err := NewTempFileWriter()
		if err != nil {
			fmt.Printf("Error creating temp writer: %s\n", err)
			continue
		}

		_, err = io.Copy(writer, reader)
		if err != nil {
			fmt.Printf("Error copying file %s: %s\n", file.Name, err)
		} else {

			var item Unpacked
			item.Source = path
			item.ArchivePath = file.Name
			item.TempPath = writer.Name()

			unpacked = append(unpacked, item)
		}

		util.CloseOrLog(writer)
		util.CloseOrLog(reader)
	}

	return unpacked, nil
}
