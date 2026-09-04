package storage

import (
	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/util"
)

// CreateDirectories creates required directories
func CreateDirectories() {
	util.CreateDirIfNotExists(config.AssetStorageBaseDir, FilePermissions)
	util.CreateDirIfNotExists(config.AssetStorageTempDir, FilePermissions)
}
