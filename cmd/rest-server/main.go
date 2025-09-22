package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	mdsqlite "github.com/c8121/asset-storage/internal/metadata-sqlite"
)

func main() {
	mdsqlite.Open()
	defer mdsqlite.Close()

	router := gin.Default()
	router.GET("/assets/list", listAssets)

	router.Run("localhost:8080")
}

func listAssets(c *gin.Context) {
	items, err := mdsqlite.ListAssets(0, 10)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.IndentedJSON(http.StatusOK, items)
}
