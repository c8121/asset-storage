package restapi

import (
	"fmt"
	"net/http"

	"github.com/c8121/asset-storage/internal/collections"
	metadata_db "github.com/c8121/asset-storage/internal/metadata-db"
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

	//Create/Update meta-data-database
	err = metadata_db.AddCollection(collection)
	if err != nil {
		fmt.Printf("Error adding collecton-data to database: %s\n", err)
	}

	c.IndentedJSON(http.StatusOK, collection)
}

// ListCollections is a rest-api handler to send a list of collections
func ListCollections(c *gin.Context) {

	var listFilter *metadata_db.CollectionListFilter = nil
	err := c.ShouldBind(&listFilter)
	if err != nil || listFilter == nil {
		util.LogError(fmt.Errorf("failed to parse request: %w", err))
		listFilter = &metadata_db.CollectionListFilter{
			Offset: 0,
			Count:  DefaultListItemCount,
		}
	}
	//fmt.Printf("Filter: %v\n", listFilter)

	items, err := metadata_db.ListCollections(listFilter)
	if err != nil {
		util.LogError(c.AbortWithError(http.StatusInternalServerError, err))
		return
	}

	if len(items) > 0 {
		c.IndentedJSON(http.StatusOK, items)
	} else {
		//https://github.com/gin-gonic/gin/issues/125 ?
		c.Data(http.StatusOK, "application/json", []byte("[]"))
	}
}
