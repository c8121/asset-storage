package rest_api

import (
	"fmt"
	"net/http"

	"github.com/c8121/asset-storage/internal/collections"
	"github.com/c8121/asset-storage/internal/metadata_db"
	"github.com/c8121/asset-storage/internal/metadata_db_conn"
	"github.com/c8121/asset-storage/internal/util"
	"github.com/gin-gonic/gin"
)

type (
	UpdateCollectionRequest struct {
		UUID        string
		Name        string
		Description string
		Owner       string
		AssetHashes []string
	}
)

// GetCollection is a rest-api handler to send the collection content
func GetCollection(c *gin.Context) {

	uuid := c.Param("uuid")
	if len(uuid) < 32 {
		util.LogError(c.AbortWithError(http.StatusNotFound, fmt.Errorf("invalid uuid")))
		return
	}

	collection, err := collections.LoadById(uuid)
	if err != nil {
		util.LogError(c.AbortWithError(http.StatusNotFound, fmt.Errorf("invalid uuid (not found)")))
		return
	}

	c.IndentedJSON(http.StatusOK, collection)
}

// AddCollection is a rest-api handler to create or update a collection
func AddCollection(c *gin.Context) {

	var req UpdateCollectionRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if len(req.AssetHashes) == 0 {
		c.JSON(http.StatusBadRequest, "No asset hashes given")
		return
	}

	collection, err := collections.AddCollection(
		req.UUID,
		req.Name,
		req.Description,
		req.Owner,
		req.AssetHashes)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	tx, err := metadata_db_conn.BeginTransaction()
	if err != nil {
		util.LogError(c.AbortWithError(http.StatusInternalServerError, err))
		return
	}

	//Create/Update meta-data-database
	err = metadata_db.AddCollection(tx, collection)
	if err != nil {
		fmt.Printf("Error adding collection-data to database: %s\n", err)
		metadata_db_conn.RollbackOrLog(tx)
	}

	util.LogError(metadata_db_conn.CommitOrLog(tx))

	c.IndentedJSON(http.StatusOK, collection)
}

// RemoveCollection is a rest-api handler to remove a collection
func RemoveCollection(c *gin.Context) {

	uuid := c.PostForm("UUID")
	if len(uuid) < 32 {
		c.JSON(http.StatusBadRequest, "UUID missing or invalid")
		return
	}

	// Remove from filesystem
	err := collections.RemoveCollection(uuid)
	if err != nil {
		util.LogError(err)
	}

	tx, err := metadata_db_conn.BeginTransaction()
	if err != nil {
		util.LogError(c.AbortWithError(http.StatusInternalServerError, err))
		return
	}

	collection := &collections.JsonAssetCollection{UUID: uuid}

	// Remove from DB
	err = metadata_db.RemoveCollection(tx, collection)
	if err != nil {
		fmt.Printf("Error removing collection-data to database: %s\n", err)
		metadata_db_conn.RollbackOrLog(tx)
	}

	util.LogError(metadata_db_conn.CommitOrLog(tx))

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
