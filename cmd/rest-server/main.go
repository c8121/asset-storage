package main

import (
	"net/http"

	restapi "github.com/c8121/asset-storage/internal/rest-api"
	"github.com/c8121/asset-storage/internal/util"
	"github.com/gin-gonic/gin"

	mdsqlite "github.com/c8121/asset-storage/internal/metadata-sqlite"
)

func main() {
	mdsqlite.Open()
	defer mdsqlite.Close()

	router := gin.Default()
	router.GET("/assets/list", listAssets)
	router.GET("/assets/list/:offset", listAssets)
	router.GET("/assets/list/:offset/:count", listAssets)
	router.GET("/assets/thumbnail/:hash", restapi.Thumbnail)

	util.Check(router.Run("localhost:8080"), "Failed to start server")
}

func listAssets(c *gin.Context) {

	items, err := mdsqlite.ListAssets(
		util.Atoi(c.Param("offset"), 0),
		util.Atoi(c.Param("count"), 10))
	if err != nil {
		util.LogError(c.AbortWithError(500, err))
		return
	}

	c.IndentedJSON(http.StatusOK, items)
}
