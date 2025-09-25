package restapi

import (
	"fmt"
	"io"
	"net/http"

	"github.com/c8121/asset-storage/internal/metadata"
	mdsqlite "github.com/c8121/asset-storage/internal/metadata-sqlite"
	"github.com/c8121/asset-storage/internal/storage"
	"github.com/c8121/asset-storage/internal/util"
	"github.com/gin-gonic/gin"
)

// GetAsset is a rest-api handler to send the asset content
func GetAsset(c *gin.Context) {

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

	reader, err := storage.Open(hash)
	if err != nil {
		util.LogError(c.AbortWithError(http.StatusInternalServerError, err))
	}
	defer reader.Close()

	c.Status(http.StatusOK)
	c.Header("Content-Type", meta.MimeType)

	buf := make([]byte, 4096)
	for {
		n, err := reader.Read(buf)
		if err == io.EOF {
			break
		}
		c.Writer.Write(buf[:n])
		c.Writer.Flush()
	}
}

// ListAssets is a rest-api handler to send a list of assets
func ListAssets(c *gin.Context) {

	items, err := mdsqlite.ListAssets(
		util.Atoi(c.Param("offset"), 0),
		util.Atoi(c.Param("count"), 10))
	if err != nil {
		util.LogError(c.AbortWithError(http.StatusInternalServerError, err))
		return
	}

	c.IndentedJSON(http.StatusOK, items)
}
