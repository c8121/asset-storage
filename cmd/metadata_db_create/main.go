package main

import (
	"fmt"

	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/metadata"
	"github.com/c8121/asset-storage/internal/metadata_db_conn"
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

	handler := func(path string) {

		tx, err := metadata_db_conn.BeginTransaction()
		if err != nil {
			util.LogError(err)
			return
		}

		if meta, err := metadata.LoadIfExists(path); err == nil {
			if err = metadata_db_entity.AddMetaData(tx, meta); err == nil {
				fmt.Printf("Added '%s'\n", path)
			} else {
				util.LogError(err)
			}
		} else {
			util.LogError(err)
		}

		util.LogError(metadata_db_conn.CommitOrLog(tx))
	}

	metadata.Walk(handler)
}
