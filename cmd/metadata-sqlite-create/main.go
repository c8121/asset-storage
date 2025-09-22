package main

import (
	"fmt"
	"os"
	"path/filepath"
	"encoding/json"

	"github.com/c8121/asset-storage/internal/util"
	"github.com/c8121/asset-storage/internal/storage"
	"github.com/c8121/asset-storage/internal/metadata"
	"github.com/c8121/asset-storage/internal/metadata-sqlite"
)

func main() {
	mdsqlite.Open()
	defer mdsqlite.Close()

	util.Check(readAllMetaData(metadata.BaseDir()), "Failed to read meta-data directory")
}

// readAllMetaData recurively find JSON meta-data an write to SQLite database
func readAllMetaData(path string) error {

	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, file := range entries {
		filePath := filepath.Join(path, file.Name())
		stat, statErr := os.Stat(filePath)
		if statErr != nil {
			return statErr
		}
		if stat.IsDir() {
			if err = readAllMetaData(filePath) ; err != nil {
				return err
			}
		} else {
			if hash, meta, err := readMetaData(filePath) ; err == nil {
				if err = mdsqlite.AddMetaData(hash, &meta) ; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}
	return nil
}

// readMetaData Read JSON-file, return hash + AssetMetadata
func readMetaData(jsonFile string) (hash string, meta metadata.AssetMetadata, err error) {
	fmt.Printf("%s\n", jsonFile)

	if buf, err := os.ReadFile(jsonFile); err == nil {
		err = json.Unmarshal(buf, &meta)
		return storage.HashFromPath(jsonFile), meta, err
	} else {
		return "", meta, err
	}
}
