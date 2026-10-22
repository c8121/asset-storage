package metadata_db

import (
	"fmt"

	"github.com/c8121/asset-storage/internal/metadata_db_conn"
	"github.com/c8121/asset-storage/internal/metadata_db_entity"
	"github.com/c8121/asset-storage/internal/util"
)

func CanRead(assetId int64) bool {
	return true
}

func CanReadHash(hash string) bool {

	asset, err := metadata_db_entity.LoadAsset(metadata_db_conn.GetDatabase(), hash)
	if err != nil || asset == nil {
		util.LogError(err)
		return false
	}

	if !CanRead(asset.Id) {
		fmt.Printf("Access dinied to %s\n", hash)
		return false
	}

	return true
}
