package main

import (
	"github.com/c8121/asset-storage/internal/config"
	restapi "github.com/c8121/asset-storage/internal/rest-api"
	"github.com/c8121/asset-storage/internal/storage"
	"github.com/c8121/asset-storage/internal/util"
	"github.com/gin-gonic/gin"

	mdsqlite "github.com/c8121/asset-storage/internal/metadata-sqlite"
)

func main() {

	config.LoadDefault()
	storage.Init()

	mdsqlite.Open()
	defer mdsqlite.Close()

	router := gin.Default()
	restapi.CreateRoutes(router)

	util.PanicOnError(router.Run("localhost:8080"), "Failed to start server")
}
