package main

import (
	"fmt"

	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/faces"
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
			fmt.Printf("Skip, %d faced stored for: %s\n", len(*facedStored), hash)
			return
		}

		assetId := metadata_db_entity.GetAssetId(hash)

		facesFound, err := faces.GetFaces(hash)
		if err != nil {
			fmt.Printf("Cannot get faces from %s: %s\n", hash, err)
		} else {
			fmt.Printf("Found %d faces in %s\n", len(*facesFound), hash)

			if err = metadata_db_entity.RemoveFaces(assetId); err != nil {
				util.LogError(err)
			}

			if _, err = metadata_db_entity.AddFaces(assetId, facesFound); err != nil {
				util.LogError(err)
			}
		}
	}

	storage.Walk(handler)

}
