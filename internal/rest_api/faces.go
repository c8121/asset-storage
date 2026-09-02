package rest_api

import (
	"fmt"
	"net/http"

	"github.com/c8121/asset-storage/internal/metadata_db_entity"
	"github.com/c8121/asset-storage/internal/util"
	"github.com/gin-gonic/gin"
)

// GetFaces is a rest-api endpoint for metadata_db_entity.GetFaces
func GetFaces(c *gin.Context) {

	hash := c.Param("hash")
	if len(hash) < 32 {
		util.LogError(c.AbortWithError(http.StatusNotFound, fmt.Errorf("invalid hash")))
		return
	}

	faces, err := metadata_db_entity.GetFaces(hash)
	if err != nil {
		util.LogError(c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to get faces: %s", err)))
		return
	}

	if len(*faces) > 0 {
		c.IndentedJSON(http.StatusOK, faces)
	} else {
		//https://github.com/gin-gonic/gin/issues/125 ?
		c.Data(http.StatusOK, "application/json", []byte("[]"))
	}
}
