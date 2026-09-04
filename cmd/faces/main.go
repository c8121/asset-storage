package main

import (
	"fmt"

	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/faces"
	"github.com/c8121/asset-storage/internal/metadata_db_conn"
	"github.com/c8121/asset-storage/internal/metadata_db_entity"
	"github.com/c8121/asset-storage/internal/metadata_sqlite"
	"github.com/c8121/asset-storage/internal/storage"
	"github.com/c8121/asset-storage/internal/util"
)

func main() {

	config.LoadDefault()

	metadata_sqlite.Open()
	defer metadata_sqlite.Close()

	handler := func(path string) {
		hash := storage.HashFromStoragePath(path)

		facedStored, err := metadata_db_entity.GetFaces(hash)
		if err != nil {
			util.LogError(err)
		}
		if facedStored != nil && len(*facedStored) > 0 {
			fmt.Printf("Skip, %d faces stored for: %s\n", len(*facedStored), hash)
			return
		}

		assetMeta, err := metadata_db_entity.GetMetaData(hash)
		if err != nil {
			util.LogError(err)
			return
		}

		tx, err := metadata_db_conn.BeginTransaction()
		if err != nil {
			util.LogError(err)
			return
		}

		facesFound, err := faces.DetectFaces(hash)
		if err != nil {
			fmt.Printf("Cannot get faces from %s: %s\n", hash, err)
		} else {
			fmt.Printf("Found %d faces in %s\n", len(*facesFound), hash)

			if err = metadata_db_entity.RemoveFaces(tx, assetMeta.Id); err != nil {
				util.LogError(err)
			}

			if _, err = metadata_db_entity.AddFaces(tx, assetMeta.Id, facesFound); err != nil {
				util.LogError(err)
			}
		}

		util.LogError(metadata_db_conn.CommitOrLog(tx))
	}

	storage.Walk(handler)
}
