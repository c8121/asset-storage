package rest_api

import (
	"errors"
	"net/http"
	"strconv"

	metadata_db_entity "github.com/c8121/asset-storage/internal/metadata_db/entity"
	"github.com/c8121/go_db_entities"
	"github.com/c8121/go_vue_spa/util"
	"github.com/gin-gonic/gin"
)

// ListPathItems is a rest-api handler
func ListPathItems(c *gin.Context) {

	parentId, err := strconv.ParseInt(c.Param("parent"), 10, 64)
	if err != nil {
		parentId = 0
	}

	items, err := metadata_db_entity.ListPathItems(parentId)
	if err != nil && !errors.Is(go_db_entities.ErrNotFound, err) {
		util.AboutWithStatus(c, http.StatusInternalServerError, "failed to load path-items", err)
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

	item, err := metadata_db_entity.FindPathItem(c.PostForm("path"))
	if err != nil {
		util.AboutWithStatus(c, http.StatusInternalServerError, "failed to find path-items", err)
		return
	}

	if item != nil {
		c.IndentedJSON(http.StatusOK, item)
	} else {
		//https://github.com/gin-gonic/gin/issues/125 ?
		c.Data(http.StatusOK, "application/json", []byte("[]"))
	}
}
