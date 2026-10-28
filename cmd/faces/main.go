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

	if !faces.CheckFaceRestServiceAvailable() {
		util.PrintNotifications()
		return
	}

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

		assetMeta, err := metadata_db_entity.LoadAsset(hash)
		if err != nil {
			util.LogError(err)
			return
		}

		facesFound, err := faces.DetectFaces(hash)
		if err != nil {
			fmt.Printf("Cannot get faces from %s: %s\n", hash, err)
		} else {
			fmt.Printf("Found %d faces in %s\n", len(*facesFound), hash)

			if err = metadata_db_entity.RemoveFaces(assetMeta.Id); err != nil {
				util.LogError(err)
			}

			if _, err = metadata_db_entity.AddFaces(assetMeta.Id, facesFound); err != nil {
				util.LogError(err)
			}
		}
	}

	storage.Walk(handler)
}
