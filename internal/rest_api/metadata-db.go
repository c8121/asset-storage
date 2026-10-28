package rest_api

import (
	"net/http"
	"strconv"

	"github.com/c8121/asset-storage/internal/metadata_db"
	"github.com/c8121/asset-storage/internal/metadata_db_entity"
	"github.com/c8121/asset-storage/internal/util"
	"github.com/gin-gonic/gin"
)

func ListMimeTypes(c *gin.Context) {

	items, err := metadata_db_entity.ListMimeTypes()
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

// ListPathItems is a rest-api handler
func ListPathItems(c *gin.Context) {

	parentId, err := strconv.ParseInt(c.Param("parent"), 10, 64)
	if err != nil {
		parentId = 0
	}

	items, err := metadata_db_entity.ListPathItems(parentId)
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

// FindPathItem is a rest-api handler
func FindPathItem(c *gin.Context) {

	item, err := metadata_db.FindPathItem(c.PostForm("path"))
	if err != nil {
		util.LogError(c.AbortWithError(http.StatusInternalServerError, err))
		return
	}

	if item != nil {
		c.IndentedJSON(http.StatusOK, item)
	} else {
		//https://github.com/gin-gonic/gin/issues/125 ?
		c.Data(http.StatusOK, "application/json", []byte("[]"))
	}
}
