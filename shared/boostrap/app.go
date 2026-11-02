package boostrap

import (
	"github.com/c8121/asset-storage/internal/faces"
	"github.com/c8121/asset-storage/internal/filter/commands"
	"github.com/c8121/asset-storage/internal/fulltext"
	"github.com/c8121/asset-storage/internal/metadata_db/sqlite"
	apputil "github.com/c8121/asset-storage/internal/util"
	"github.com/c8121/asset-storage/shared/config"
	"github.com/c8121/go_vue_spa/util"
)

func InitAssetStorageApp() {

	config.AppConfig = config.NewConfigFromFlags()

	apputil.CreateDirIfNotExists(config.AppConfig.AssetStorageBaseDir, config.AppConfig.FilePermissions)
	apputil.CreateDirIfNotExists(config.AppConfig.AssetStorageTempDir, config.AppConfig.FilePermissions)

	apputil.CreateDirIfNotExists(config.AppConfig.AssetMetaDataBaseDir, config.AppConfig.FilePermissions)

	apputil.CreateDirIfNotExists(config.AppConfig.AssetCollectionsBaseDir, config.AppConfig.FilePermissions)

	apputil.CreateDirIfNotExists(config.AppConfig.AssetFulltextIndexBaseDir, config.AppConfig.FilePermissions)

	filter_commands.CheckAndNotify()
	faces.CheckFaceRestServiceAvailable()

	util.LogError(fulltext.Open())

	metadata_db_sqlite.Open()
}

func FinalizeAssetStorageApp() {
	fulltext.Close()
	metadata_db_sqlite.Close()
}
