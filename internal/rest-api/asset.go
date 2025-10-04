package restapi

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/c8121/asset-storage/internal/metadata"
	metadata_db "github.com/c8121/asset-storage/internal/metadata-db"
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
	defer util.CloseOrLog(reader)

	c.Status(http.StatusOK)
	c.Header("Content-Type", meta.MimeType)

	buf := make([]byte, storage.IoBufferSize)
	for {
		n, err := reader.Read(buf)
		if n == 0 && err == io.EOF {
			break
		}
		util.PanicOnIoError(err, "Failed to read file")

		if _, err = c.Writer.Write(buf[:n]); err != nil {
			util.LogError(c.AbortWithError(http.StatusInternalServerError, err))
			break
		}
		c.Writer.Flush()
	}
}

// ListAssets is a rest-api handler to send a list of assets
func ListAssets(c *gin.Context) {

	var filter = &metadata_db.AssetListFilter{
		MimeType: strings.ReplaceAll(
			strings.ReplaceAll(c.Param("mimetype"),
				"_", "/"),
			"*", "%"),
		Offset: util.Atoi(c.Param("offset"), 0),
		Count:  util.Atoi(c.Param("count"), 0),
	}

	items, err := metadata_db.ListAssets(filter)
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
