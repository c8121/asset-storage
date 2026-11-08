package bootstrap

import (
	"flag"

	"github.com/c8121/asset-storage/internal/faces"
	"github.com/c8121/asset-storage/internal/filter/commands"
	"github.com/c8121/asset-storage/internal/fulltext"
	"github.com/c8121/asset-storage/shared/config"
	"github.com/c8121/go_vue_spa/util"
)

func InitAssetStorageApp() {

	config.AppConfig = config.NewConfigFromFlags(flag.ExitOnError)

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
	err := util.CreateDirIfNotExists(path, config.AppConfig.FilePermissions)
	if err != nil {
		panic(err)
	}
}

func FinalizeAssetStorageApp() {
	fulltext.Close()
	CloseDatabase()
}
