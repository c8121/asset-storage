package main

import (
	"fmt"

	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/faces"
	"github.com/c8121/asset-storage/internal/metadata_db_entity"
	"github.com/c8121/asset-storage/internal/metadata_sqlite"
	"github.com/c8121/asset-storage/internal/storage"
)

func main() {

	config.LoadDefault()
	storage.CreateDirectories()

	metadata_sqlite.Open()
	defer metadata_sqlite.Close()

	handler := func(path string) {
		hash := storage.HashFromStoragePath(path)
		assetId := metadata_db_entity.GetAssetId(hash)

		facesFound, err := faces.GetFaces(hash)
		if err != nil {
			fmt.Printf("Cannot get faces from %s: %s\n", hash, err)
		} else {
			fmt.Printf("Found %d faces in %s\n", len(facesFound.Faces), hash)

			for _, face := range facesFound.Faces {

				fmt.Printf("  %s\n", face.Image)

				err := metadata_db_entity.AddFace(assetId, &face)
				if err != nil {
					fmt.Printf("Cannot add face %s: %s\n", hash, err)
				}
			}
		}
	}

	storage.Walk(handler)

}
