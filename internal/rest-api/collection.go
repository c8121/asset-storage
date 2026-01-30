package restapi

import (
	"fmt"
	"net/http"

	"github.com/c8121/asset-storage/internal/collections"
	"github.com/c8121/asset-storage/internal/util"
	"github.com/gin-gonic/gin"
)

type (
	CreateCollectionRequest struct {
		Name        string
		Description string
		Owner       string
		AssetHashes []string
	}
)

// GetCollection is a rest-api handler to send the collection content
func GetCollection(c *gin.Context) {

	hash := c.Param("hash")
	if len(hash) < 32 {
		util.LogError(c.AbortWithError(http.StatusNotFound, fmt.Errorf("invalid hash")))
		return
	}

	collection, err := collections.LoadByHash(hash)
	if err != nil {
		util.LogError(c.AbortWithError(http.StatusNotFound, fmt.Errorf("invalid hash (not found)")))
		return
	}

	c.IndentedJSON(http.StatusOK, collection)
}

// AddCollection is a rest-api handler to create collections
func AddCollection(c *gin.Context) {

	var req CreateCollectionRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if len(req.AssetHashes) == 0 {
		c.JSON(http.StatusBadRequest, "No asset hashes given")
		return
	}

	collection, err := collections.AddCollection(req.Name, req.Description, req.Owner, req.AssetHashes)
	if err != nil {
		util.LogError(c.AbortWithError(http.StatusNotFound, fmt.Errorf("invalid hash (not found)")))
		return
	}

	c.IndentedJSON(http.StatusOK, collection)
}
