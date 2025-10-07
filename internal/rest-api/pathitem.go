package restapi

import (
	"net/http"

	metadata_db "github.com/c8121/asset-storage/internal/metadata-db"
	"github.com/c8121/asset-storage/internal/util"
	"github.com/gin-gonic/gin"
)

// ListPathItems is a rest-api handler
func ListPathItems(c *gin.Context) {

	items, err := metadata_db.ListPathItems(util.Atoi(c.Param("parent"), 0))
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
