package main

import (
	restapi "github.com/c8121/asset-storage/internal/rest-api"
	spaserver "github.com/c8121/asset-storage/internal/spa-server"
	"github.com/c8121/asset-storage/internal/util"
	"github.com/gin-gonic/gin"

	mdsqlite "github.com/c8121/asset-storage/internal/metadata-sqlite"
)

func main() {
	mdsqlite.Open()
	defer mdsqlite.Close()

	router := gin.Default()

	restapi.CreateRoutes(router)

	spaserver.CreateRoutes(router)

	util.Check(router.Run("localhost:8080"), "Failed to start server")
}
