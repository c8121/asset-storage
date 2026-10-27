package main

import (
	"fmt"

	"github.com/c8121/asset-storage/internal/collections"
	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/metadata"
	"github.com/c8121/asset-storage/internal/metadata_db"
	"github.com/c8121/asset-storage/internal/metadata_db_entity"
	"github.com/c8121/asset-storage/internal/metadata_sqlite"
	"github.com/c8121/asset-storage/internal/util"
)

/*
	Update meta-data-database by reading all meta-data JSON-files
	and writing contents to database.

	Not required if database is intact, because cmd/add also updates the database.
*/

func main() {

	config.LoadDefault()

	metadata_sqlite.Open()
	defer metadata_sqlite.Close()

	createFilesMetaDataDb()
	createCollectionsMetaDataDb()
}

func createFilesMetaDataDb() {
	handler := func(path string) {

		if meta, err := metadata.LoadIfExists(path); err == nil {
			if err = metadata_db_entity.AddAsset(meta); err == nil {
				fmt.Printf("Added '%s'\n", path)
			} else {
				util.LogError(err)
			}
		} else {
			util.LogError(err)
		}
	}

	metadata.Walk(handler)
}

func createCollectionsMetaDataDb() {
	handler := func(path string) {

		if collection, err := collections.LoadIfExists(path); err == nil {
			if err = metadata_db.AddCollection(collection); err == nil {
				fmt.Printf("Added '%s'\n", path)
			} else {
				util.LogError(err)
			}
		} else {
			util.LogError(err)
		}
	}

	collections.Walk(handler)
}
