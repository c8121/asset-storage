package restapi

import (
	"fmt"
	"io"
	"net/http"

	"github.com/c8121/asset-storage/internal/metadata"
	"github.com/c8121/asset-storage/internal/util"
	"github.com/gin-gonic/gin"
)

// GetFiltered is a rest-api handler to filterAsset/convert an asset
func GetFiltered(c *gin.Context) {

	hash := c.Param("hash")
	if len(hash) < 32 {
		util.LogError(c.AbortWithError(http.StatusNotFound, fmt.Errorf("invalid hash")))
		return
	}

	filterName := c.Param("filterAsset")
	if len(filterName) < 1 {
		util.LogError(c.AbortWithError(http.StatusNotFound, fmt.Errorf("no filterAsset name given")))
		return
	}

	filterParamsReader := c.Request.Body
	b, err := io.ReadAll(filterParamsReader)
	if err != nil {
		util.LogError(c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Failed to read request body")))
		return
	}
	filterParams := string(b)

	meta, err := metadata.LoadByHash(hash)
	if err != nil {
		util.LogError(c.AbortWithError(http.StatusNotFound, fmt.Errorf("invalid hash (not found)")))
		return
	}

	if bytes, mimeType, err := filterAsset(hash, meta, filterName, filterParams); err == nil {
		c.Data(http.StatusOK, mimeType, bytes)
		return
	} else {
		util.LogError(err)
		util.LogError(c.AbortWithError(http.StatusInternalServerError, err))
	}
}
