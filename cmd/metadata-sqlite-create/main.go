package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/metadata"
	mdsqlite "github.com/c8121/asset-storage/internal/metadata-sqlite"
	"github.com/c8121/asset-storage/internal/util"
)

/*
	Update meta-data-database by reading all meta-data JSON-files
	and writing contents to database.

	Not required if database is intact, because cmd/add also updates the database.
*/

func main() {

	config.LoadDefault()

	mdsqlite.Open()
	defer mdsqlite.Close()

	util.PanicOnError(readAllMetaData(config.AssetMetaDataBaseDir), "Failed to read meta-data directory")
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
			if err = readAllMetaData(filePath); err != nil {
				return err
			}
		} else {
			if hash, meta, err := metadata.LoadIfExists(filePath); err == nil {
				if err = mdsqlite.AddMetaData(hash, &meta); err != nil {
					return err
				} else {
					fmt.Printf("Added '%s'\n", filePath)
				}
			} else {
				return err
			}
		}
	}
	return nil
}
