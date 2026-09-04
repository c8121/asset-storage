package main

import (
	"flag"
	"fmt"
	"os/user"

	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/metadata"
	"github.com/c8121/asset-storage/internal/metadata_db"
	"github.com/c8121/asset-storage/internal/metadata_db_entity"
	"github.com/c8121/asset-storage/internal/metadata_sqlite"
	"github.com/c8121/asset-storage/internal/storage"
	"github.com/c8121/asset-storage/internal/util"
)

var (
	_, currentUserErr = user.Current()
	hash              = flag.String("hash", "", "Hash of asset to delete")
	path              = flag.String("path", "", "Path of assets to delete")
	force             = flag.Bool("f", false, "Force deletion")
)

func main() {

	flag.Parse()

	if currentUserErr != nil {
		panic(currentUserErr)
	}

	config.LoadDefault()

	metadata_sqlite.Open()
	defer metadata_sqlite.Close()

	if *hash != "" {
		removeAsset(*hash)
	} else if *path != "" {
		removePath(*path)
	} else {
		fmt.Println("Please specify a -hash or a -path flag")
	}
}

func removeAsset(hash string) {

	assetId := metadata_db_entity.GetAssetId(hash)
	if assetId == 0 {
		fmt.Printf("Asset '%s' not found\n", hash)
		return
	}

	fmt.Printf("Removing asset %s, id=%d\n", hash, assetId)

	if !*force {
		if !util.CliConfirm("Remove asset " + hash) {
			return
		}
	}

	if err := metadata_db_entity.RemoveFaces(assetId); err != nil {
		fmt.Printf("Error removing asset faces %s, id=%d: %s\n", hash, assetId, err)
	}

	if _, err := metadata_db_entity.RemoveMetaData(assetId, 0); err != nil {
		fmt.Printf("Error removing asset db-metadata %s, id=%d: %s\n", hash, assetId, err)
	}

	if err := metadata.RemoveMetaData(hash); err != nil {
		fmt.Printf("Error removing asset metadata %s, id=%d: %s\n", hash, assetId, err)
	}

	if err := storage.Remove(hash); err != nil {
		fmt.Printf("Error removing asset metadata %s, id=%d: %s\n", hash, assetId, err)
	}
}

func removePath(path string) {

	pathItem, err := metadata_db.FindPathItem(path)
	if err != nil {
		fmt.Printf("Error finding path item '%s': %s\n", path, err.Error())
		return
	}
	if pathItem == nil || pathItem.Id == 0 {
		fmt.Printf("Path '%s' not found\n", path)
		return
	}

	fmt.Printf("Removing path '%s'\n", path)
	removePathItem(pathItem)
}

func removePathItem(pathItem *metadata_db_entity.PathItem) {

	children, err := metadata_db.ListPathItems(int(pathItem.Id))
	if err != nil {
		fmt.Printf("Error listing path items: %s\n", err.Error())
	} else {
		for _, child := range children {
			removePathItem(&child)
		}
	}

	fmt.Printf("Removing path '%s', id=%d\n", pathItem.Name, pathItem.Id)

	if !*force {
		if !util.CliConfirm("Remove path '" + pathItem.Name + "' including all assets") {
			return
		}
	}

	assetIds, err := metadata_db_entity.GetAssetIdsFromPath(pathItem.Id)
	if err != nil {
		fmt.Printf("Error finding assets from path item '%s': %s\n", pathItem.Id, err.Error())
	} else {
		for _, assetId := range assetIds {
			hash := metadata_db_entity.GetAssetHash(assetId)
			if hash == "" {
				fmt.Printf("Asset %d not found\n", assetId)
			}

			fmt.Printf("Removing asset '%s', id=%d\n", hash, assetId)

			remainingOrigins, err := metadata_db_entity.RemoveMetaData(assetId, pathItem.Id)
			if remainingOrigins > 0 {
				fmt.Printf("Keep meta-data, it has remaining origins: hash=%s, id=%d, remains=%d\n", hash, assetId, remainingOrigins)
			} else if err != nil {
				fmt.Printf("Error removing asset db-metadata %s, id=%d: %s\n", hash, assetId, err)
			} else {

				if err := metadata_db_entity.RemoveFaces(assetId); err != nil {
					fmt.Printf("Error removing asset faces %s, id=%d: %s\n", hash, assetId, err)
				}

				if err := metadata.RemoveMetaData(hash); err != nil {
					fmt.Printf("Error removing asset metadata %s, id=%d: %s\n", hash, assetId, err)
				}

				if err := storage.Remove(hash); err != nil {
					fmt.Printf("Error removing asset metadata %s, id=%d: %s\n", hash, assetId, err)
				}
			}
		}

		if err = metadata_db_entity.RemovePathIfEmpty(pathItem.Id); err != nil {
			fmt.Printf("Error removing path item '%s', id=%d: %s\n", pathItem.Name, pathItem.Id, err.Error())
			return
		}
	}

}
