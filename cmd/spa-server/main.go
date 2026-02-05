package main

import (
	"github.com/c8121/asset-storage/internal/config"
	restapi "github.com/c8121/asset-storage/internal/rest-api"
	spaserver "github.com/c8121/asset-storage/internal/spa-server"
	"github.com/c8121/asset-storage/internal/storage"
	"github.com/c8121/asset-storage/internal/util"
	"github.com/gin-gonic/gin"

	mdsqlite "github.com/c8121/asset-storage/internal/metadata-sqlite"
)

func main() {

	config.LoadDefault()
	storage.CreateDirectories()

	mdsqlite.Open()
	defer mdsqlite.Close()

	router := gin.Default()
	router.SetTrustedProxies(nil)

	restapi.CreateRoutes(router)

	spaserver.CreateRoutes(router)

	util.PanicOnError(router.Run(config.ListenAddress), "Failed to start server")
}
