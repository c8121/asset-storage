package rest_api

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/c8121/asset-storage/internal/collections"
	"github.com/c8121/asset-storage/internal/metadata_db"
	"github.com/c8121/go_vue_spa/util"
	"github.com/gin-gonic/gin"
)

type (
	UpdateCollectionRequest struct {
		UUID        string
		Name        string
		Description string
		Owner       string
		Assets      []string
	}
)

// GetCollection is a rest-api handler to send the collection content
func GetCollection(c *gin.Context) {

	uuid := c.Param("uuid")
	if len(uuid) < 32 {
		util.AbortWithStatus(c, http.StatusNotFound, "invalid uuid", nil)
		return
	}

	collection, err := collections.LoadById(uuid)
	if err != nil {
		util.AbortWithStatus(c, http.StatusNotFound, "invalid uuid", err)
		return
	}

	c.IndentedJSON(http.StatusOK, collection)
}

// AddCollection is a rest-api handler to create or update a collection
func AddCollection(c *gin.Context) {

	var req UpdateCollectionRequest
	err := c.BindJSON(&req)
	if err != nil {
		util.AbortWithStatus(c, http.StatusInternalServerError, "Failed to parse request", err)
		return
	}

	if len(req.Assets) == 0 {
		util.AbortWithStatus(c, http.StatusBadRequest, "No asset hashes given", nil)
		return
	}

	collection, err := collections.AddCollection(
		req.UUID,
		req.Name,
		req.Description,
		req.Owner,
		req.Assets)

	if err != nil {
		util.AbortWithStatus(c, http.StatusInternalServerError, "Failed to add collection", err)
		return
	}

	//Create/Update meta-data-database
	err = metadata_db.AddCollection(collection)
	if err != nil {
		slog.Warn("Error adding collection-data to database", "error", err)
	}

	c.IndentedJSON(http.StatusOK, collection)
}

// RemoveCollection is a rest-api handler to remove a collection
func RemoveCollection(c *gin.Context) {

	uuid := c.PostForm("UUID")
	if len(uuid) < 32 {
		util.AbortWithStatus(c, http.StatusNotFound, "invalid uuid", nil)
		return
	}

	// Remove from filesystem
	err := collections.RemoveCollection(uuid)
	if err != nil {
		util.LogError(err)
	}

	collection := &collections.JsonAssetCollection{UUID: uuid}

	// Remove from DB
	err = metadata_db.RemoveCollection(collection)
	if err != nil {
		slog.Warn("Error removing collection-data to database", "error", err)
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
		util.AbortWithStatus(c, http.StatusInternalServerError, "Failed to load collections", err)
		return
	}

	if len(items) > 0 {
		c.IndentedJSON(http.StatusOK, items)
	} else {
		//https://github.com/gin-gonic/gin/issues/125 ?
		c.Data(http.StatusOK, "application/json", []byte("[]"))
	}
}
