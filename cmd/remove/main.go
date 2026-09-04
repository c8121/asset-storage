package main

import (
	"flag"
	"fmt"
	"os/user"

	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/metadata"
	"github.com/c8121/asset-storage/internal/metadata_db"
	"github.com/c8121/asset-storage/internal/metadata_db_conn"
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

	assetMeta, err := metadata_db_entity.GetMetaData(hash)
	if err != nil {
		util.LogError(err)
		return
	}

	fmt.Printf("Removing asset %s, id=%d\n", hash, assetMeta.Id)

	if !*force {
		if !util.CliConfirm("Remove asset " + hash) {
			return
		}
	}

	tx, err := metadata_db_conn.BeginTransaction()
	if err != nil {
		util.LogError(err)
		return
	}

	if err := metadata_db_entity.RemoveFaces(tx, assetMeta.Id); err != nil {
		fmt.Printf("Error removing asset faces %s, id=%d: %s\n", hash, assetMeta.Id, err)
	}

	if _, err := metadata_db_entity.RemoveMetaData(tx, assetMeta.Id, 0); err != nil {
		fmt.Printf("Error removing asset db-metadata %s, id=%d: %s\n", hash, assetMeta.Id, err)
	}

	if err := metadata.RemoveMetaData(hash); err != nil {
		fmt.Printf("Error removing asset metadata %s, id=%d: %s\n", hash, assetMeta.Id, err)
	}

	if err := storage.Remove(hash); err != nil {
		fmt.Printf("Error removing asset metadata %s, id=%d: %s\n", hash, assetMeta.Id, err)
	}

	util.LogError(metadata_db_conn.CommitOrLog(tx))
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

			assetMeta, err := metadata_db_entity.GetMetaDataById(assetId)
			if err != nil {
				util.LogError(err)
				continue
			}

			tx, err := metadata_db_conn.BeginTransaction()
			if err != nil {
				util.LogError(err)
				return
			}

			fmt.Printf("Removing asset '%s', id=%d\n", assetMeta.Hash, assetId)

			remainingOrigins, err := metadata_db_entity.RemoveMetaData(tx, assetId, pathItem.Id)
			if remainingOrigins > 0 {
				fmt.Printf("Keep meta-data, it has remaining origins: hash=%s, id=%d, remains=%d\n", assetMeta.Hash, assetId, remainingOrigins)
			} else if err != nil {
				fmt.Printf("Error removing asset db-metadata %s, id=%d: %s\n", assetMeta.Hash, assetId, err)
			} else {

				if err := metadata_db_entity.RemoveFaces(tx, assetId); err != nil {
					fmt.Printf("Error removing asset faces %s, id=%d: %s\n", assetMeta.Hash, assetId, err)
				}

				if err := metadata.RemoveMetaData(assetMeta.Hash); err != nil {
					fmt.Printf("Error removing asset metadata %s, id=%d: %s\n", assetMeta.Hash, assetId, err)
				}

				if err := storage.Remove(assetMeta.Hash); err != nil {
					fmt.Printf("Error removing asset metadata %s, id=%d: %s\n", assetMeta.Hash, assetId, err)
				}
			}

			util.LogError(metadata_db_conn.CommitOrLog(tx))
		}

		if err = metadata_db_entity.RemovePathIfEmpty(pathItem.Id); err != nil {
			fmt.Printf("Error removing path item '%s', id=%d: %s\n", pathItem.Name, pathItem.Id, err.Error())
			return
		}
	}

}
