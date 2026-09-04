package main

import (
	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/faces"
	"github.com/c8121/asset-storage/internal/filter_commands"
	"github.com/c8121/asset-storage/internal/metadata"
	"github.com/c8121/asset-storage/internal/rest_api"
	"github.com/c8121/asset-storage/internal/spa_server"
	"github.com/c8121/asset-storage/internal/storage"
	"github.com/c8121/asset-storage/internal/util"
	"github.com/gin-gonic/gin"

	"github.com/c8121/asset-storage/internal/metadata_sqlite"
)

func main() {

	config.LoadDefault()
	storage.CreateDirectories()
	metadata.CreateDirectories()

	filter_commands.CheckAndNotify()
	faces.CheckFaceRestServiceAvailable()

	metadata_sqlite.Open()
	defer metadata_sqlite.Close()

	router := gin.Default()
	router.SetTrustedProxies(nil)

	rest_api.CreateRoutes(router)

	spa_server.CreateRoutes(router)

	if config.CertFile != "" {
		util.PanicOnError(router.RunTLS(config.ListenAddress, config.CertFile, config.KeyFile), "Failed to start server")
	} else {
		util.PanicOnError(router.Run(config.ListenAddress), "Failed to start server")
	}
}
