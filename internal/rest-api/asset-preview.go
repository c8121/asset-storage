package restapi

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/c8121/asset-storage/internal/filter"
	"github.com/c8121/asset-storage/internal/metadata"
	"github.com/c8121/asset-storage/internal/util"
	"github.com/gin-gonic/gin"
)

var (
	ThumbnailWidth = 150
)

// GetPreview is a rest-api handler to generate a preview image
func GetPreview(c *gin.Context) {

	hash := c.Param("hash")
	if len(hash) < 32 {
		util.LogError(c.AbortWithError(http.StatusNotFound, fmt.Errorf("invalid hash")))
		return
	}

	meta, err := metadata.LoadByHash(hash)
	if err != nil {
		util.LogError(c.AbortWithError(http.StatusNotFound, fmt.Errorf("invalid hash (not found)")))
		return
	}

	if bytes, mimeType, err := generateThumbnail(hash, meta); err == nil {
		c.Data(http.StatusOK, mimeType, bytes)
		return
	} else {
		util.LogError(err)
	}

	if bytes, mimeType, err := loadIconByMimeType(meta.MimeType); err == nil {
		c.Data(http.StatusOK, mimeType, bytes)
		return
	} else {
		util.LogError(err)
	}

	if bytes, mimeType, err := loadDefaultIcon(); err == nil {
		c.Data(http.StatusOK, mimeType, bytes)
		return
	} else {
		util.LogError(err)
		util.LogError(c.AbortWithError(http.StatusInternalServerError, err))
	}
}

// generateThumbnail returns a thumbnail image generated from an asset.
// Returns content, mimeType, error
func generateThumbnail(assetHash string, meta *metadata.JsonAssetMetaData) ([]byte, string, error) {

	var f = filter.GetFirstFilterByMimeType(meta.MimeType)
	if f == nil {
		return nil, "", fmt.Errorf("no filter available for mime-type: %s", meta.MimeType)
	}

	params := map[string]string{}
	params["width"] = strconv.Itoa(ThumbnailWidth)

	return f.Apply(assetHash, meta, params)
}
