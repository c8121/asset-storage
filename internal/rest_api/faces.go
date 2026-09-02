package rest_api

import (
	"fmt"
	"net/http"

	"github.com/c8121/asset-storage/internal/faces"
	"github.com/c8121/asset-storage/internal/metadata_db_entity"
	"github.com/c8121/asset-storage/internal/util"
	"github.com/gin-gonic/gin"
)

// GetFaces is a rest-api endpoint for metadata_db_entity.GetFaces
// Loads faces that already have been detected and stored to database
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

// DetectFaces is a rest-api endpoint to trigger face detection
func DetectFaces(c *gin.Context) {

	hash := c.Param("hash")
	if len(hash) < 32 {
		util.LogError(c.AbortWithError(http.StatusNotFound, fmt.Errorf("invalid hash")))
		return
	}

	assetId := metadata_db_entity.GetAssetId(hash)
	facesFound, err := faces.GetFaces(hash)
	if err != nil {
		util.LogError(c.AbortWithError(http.StatusInternalServerError, err))
	} else {
		fmt.Printf("Found %d faces in %s\n", len(*facesFound), hash)

		if err = metadata_db_entity.RemoveFaces(assetId); err != nil {
			util.LogError(err)
		}

		if list, err := metadata_db_entity.AddFaces(assetId, facesFound); err != nil {
			util.LogError(c.AbortWithError(http.StatusInternalServerError, err))
		} else {
			c.IndentedJSON(http.StatusOK, list)
		}
	}

}
