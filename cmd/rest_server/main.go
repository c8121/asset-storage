package main

import (
	"github.com/c8121/asset-storage/internal/metadata_db"
	"github.com/c8121/asset-storage/internal/users"
	"github.com/c8121/asset-storage/shared/bootstrap"
	"github.com/c8121/asset-storage/shared/config"
	"github.com/c8121/asset-storage/shared/rest_api"
	"github.com/c8121/go_vue_spa/util"
	"github.com/gin-gonic/gin"
)

func main() {

	bootstrap.InitAssetStorageApp()
	defer bootstrap.FinalizeAssetStorageApp()

	if users.UsersFileExists() {
		util.PanicOnError(users.InitLocalFileAppAuthentication(users.DefaultPasswordCost), "failed to init app authentication")
		//If user file does not exist, app can be used without auth
	}
	metadata_db.InitAccessControl()

	router := gin.Default()
	util.LogError(router.SetTrustedProxies(nil))

	rest_api.CreateRoutes(router)

	if config.AppConfig.CertFile != "" {
		util.PanicOnError(router.RunTLS(config.AppConfig.ListenAddress, config.AppConfig.CertFile, config.AppConfig.KeyFile), "Failed to start server")
	} else {
		util.PanicOnError(router.Run(config.AppConfig.ListenAddress), "Failed to start server")
	}
}
