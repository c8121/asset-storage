package main

import (
	restapi "github.com/c8121/asset-storage/internal/rest-api"
	"github.com/c8121/asset-storage/internal/util"
	"github.com/gin-gonic/gin"

	mdsqlite "github.com/c8121/asset-storage/internal/metadata-sqlite"
)

func main() {
	mdsqlite.Open()
	defer mdsqlite.Close()

	router := gin.Default()
	router.GET("/assets/:hash", restapi.GetAsset)
	router.GET("/assets/list", restapi.ListAssets)
	router.GET("/assets/list/:offset", restapi.ListAssets)
	router.GET("/assets/list/:offset/:count", restapi.ListAssets)
	router.GET("/assets/thumbnail/:hash", restapi.GetPreview)

	util.Check(router.Run("localhost:8080"), "Failed to start server")
}
