package rest_api

import (
	"fmt"
	"io"
	"net/http"

	"github.com/c8121/asset-storage/internal/access_control"
	"github.com/c8121/asset-storage/internal/metadata_db"
	"github.com/c8121/asset-storage/internal/metadata_db_entity"
	"github.com/c8121/asset-storage/internal/storage"
	"github.com/c8121/asset-storage/internal/util"
	"github.com/gin-gonic/gin"
)

var (
	DefaultListItemCount = 30
)

// GetAsset is a rest-api handler to send the asset content
func GetAsset(c *gin.Context) {

	hash := c.Param("hash")
	if len(hash) < 32 {
		util.LogError(c.AbortWithError(http.StatusNotFound, fmt.Errorf("invalid hash")))
		return
	}

	asset, err := metadata_db_entity.LoadAsset(hash)
	if err != nil || asset == nil {
		util.LogError(err)
		util.LogError(c.AbortWithError(http.StatusNotFound, fmt.Errorf("invalid hash (not found)")))
		return
	}

	acc := GetAccessContext(c)
	ctx := access_control.GetAccessController()
	if ctx != nil && ctx.CanRead(hash, acc) != nil {
		fmt.Printf("Access dinied to %s\n", hash)
		util.LogError(c.AbortWithError(http.StatusForbidden, fmt.Errorf("Forbidden")))
		return
	}

	fileName, _ := metadata_db_entity.LoadFileNameById(asset.Name)
	mimeType, _ := metadata_db_entity.LoadMimeTypeById(asset.MimeType)

	reader, err := storage.Open(hash)
	if err != nil {
		util.LogError(c.AbortWithError(http.StatusInternalServerError, err))
		return
	}
	defer util.CloseOrLog(reader)

	c.Status(http.StatusOK)

	if mimeType != nil && len(mimeType.Name) > 0 {
		c.Header("Content-Type", mimeType.Name)
	} else {
		c.Header("Content-Type", "application/octed-stream")
	}

	if fileName != nil && len(fileName.Name) > 0 {
		c.Header("Content-Disposition", "inline; filename=\""+fileName.Name+"\"")
	} else {
		c.Header("Content-Disposition", "inline")
	}

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

	var listFilter *metadata_db.AssetListFilter = nil
	err := c.ShouldBind(&listFilter)
	if err != nil || listFilter == nil {
		util.LogError(fmt.Errorf("failed to parse request: %w", err))
		listFilter = &metadata_db.AssetListFilter{
			Offset: 0,
			Count:  DefaultListItemCount,
		}
	}
	//fmt.Printf("Filter: %v\n", listFilter)

	items, err := metadata_db.ListAssets(listFilter)
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
