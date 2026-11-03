package bootstrap

import (
	"github.com/c8121/asset-storage/internal/faces"
	filter_commands "github.com/c8121/asset-storage/internal/filter/commands"
	"github.com/c8121/asset-storage/internal/fulltext"
	apputil "github.com/c8121/asset-storage/internal/util"
	"github.com/c8121/asset-storage/shared/config"
	"github.com/c8121/go_vue_spa/util"
)

func InitAssetStorageApp() {

	config.AppConfig = config.NewConfigFromFlags()

	CreateAppDir(config.AppConfig.AssetStorageBaseDir)
	CreateAppDir(config.AppConfig.AssetStorageTempDir)

	CreateAppDir(config.AppConfig.AssetMetaDataBaseDir)

	CreateAppDir(config.AppConfig.AssetCollectionsBaseDir)

	CreateAppDir(config.AppConfig.AssetFulltextIndexBaseDir)

	filter_commands.CheckAndNotify()
	faces.CheckFaceRestServiceAvailable()

	util.LogError(fulltext.Open())

	OpenDatabase()
}

func CreateAppDir(path string) {
	apputil.CreateDirIfNotExists(path, config.AppConfig.FilePermissions)
}

func FinalizeAssetStorageApp() {
	fulltext.Close()
	CloseDatabase()
}
