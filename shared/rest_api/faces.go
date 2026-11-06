package rest_api

import (
	"log/slog"
	"net/http"

	"github.com/c8121/asset-storage/internal/access_control"
	"github.com/c8121/asset-storage/internal/faces"
	metadata_db_entity "github.com/c8121/asset-storage/internal/metadata_db/entity"
	"github.com/c8121/go_vue_spa/util"
	"github.com/gin-gonic/gin"
)

// GetFaces is a rest-api endpoint for metadata_db_entity.GetFaces
// Loads faces that already have been detected and stored to database
func GetFaces(c *gin.Context) {

	hash := c.Param("hash")
	if len(hash) < 32 {
		util.AbortWithStatus(c, http.StatusNotFound, "invalid hash", nil)
		return
	}

	acc := GetAccessContext(c)
	ctx := access_control.GetAccessController()
	if ctx != nil && ctx.CanRead(hash, acc) != nil {
		util.AbortWithStatus(c, http.StatusForbidden, "forbidden", nil)
		return
	}

	facesFound, err := metadata_db_entity.GetFaces(hash)
	if err != nil {
		util.AbortWithStatus(c, http.StatusInternalServerError, "failed to get faces", err)
		return
	}

	if len(*facesFound) > 0 {
		c.IndentedJSON(http.StatusOK, facesFound)
	} else {
		//https://github.com/gin-gonic/gin/issues/125 ?
		c.Data(http.StatusOK, "application/json", []byte("[]"))
	}
}

// DetectFaces is a rest-api endpoint to trigger face detection
func DetectFaces(c *gin.Context) {

	hash := c.Param("hash")
	if len(hash) < 32 {
		util.AbortWithStatus(c, http.StatusNotFound, "invalid hash", nil)
		return
	}

	acc := GetAccessContext(c)
	ctx := access_control.GetAccessController()
	if ctx != nil && ctx.CanRead(hash, acc) != nil {
		util.AbortWithStatus(c, http.StatusForbidden, "forbidden", nil)
		return
	}

	assetMeta, err := metadata_db_entity.LoadAsset(hash)
	if err != nil {
		util.AbortWithStatus(c, http.StatusInternalServerError, "failed to load face", err)
		return
	}

	facesFound, err := faces.DetectFaces(hash)
	if err != nil {
		util.AbortWithStatus(c, http.StatusInternalServerError, "failed to detect faces", err)
	} else {
		slog.Info("Found faces", "count", len(*facesFound), "hash", hash)

		if err = metadata_db_entity.RemoveFaces(assetMeta.Id); err != nil {
			util.LogError(err)
		}

		if list, err := metadata_db_entity.AddFaces(assetMeta.Id, facesFound); err != nil {
			util.LogError(c.AbortWithError(http.StatusInternalServerError, err))
		} else {
			c.IndentedJSON(http.StatusOK, list)
		}
	}
}
