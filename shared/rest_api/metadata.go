package rest_api

import (
	"net/http"

	"github.com/c8121/asset-storage/internal/access_control"
	"github.com/c8121/asset-storage/internal/metadata"
	"github.com/c8121/go_vue_spa/util"
	"github.com/gin-gonic/gin"
)

// GetMetaData is a rest-api handler to send the asset meta-data
func GetMetaData(c *gin.Context) {

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

	meta, err := metadata.LoadByHash(hash)
	if err != nil {
		util.AbortWithStatus(c, http.StatusNotFound, "invalid hash", err)
		return
	}

	c.IndentedJSON(http.StatusOK, meta)
}
