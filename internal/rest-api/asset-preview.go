package restapi

import (
	"fmt"
	"net/http"

	"github.com/c8121/asset-storage/internal/metadata"
	"github.com/c8121/asset-storage/internal/util"
	"github.com/gin-gonic/gin"
)

// GetPreview is a rest-api handler to generate a preview image
func GetPreview(c *gin.Context) {

	hash := c.Param("hash")
	if len(hash) < 32 {
		util.LogError(c.AbortWithError(404, fmt.Errorf("invalid hash")))
		return
	}

	meta, err := metadata.LoadByHash(hash)
	if err != nil {
		util.LogError(c.AbortWithError(404, fmt.Errorf("invalid hash (not found)")))
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
		util.LogError(c.AbortWithError(500, err))
	}
}
