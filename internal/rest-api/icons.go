package restapi

import (
	"bufio"
	"bytes"
	"embed"
	_ "embed"
	"fmt"
	"net/http"
	"strings"

	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"

	"github.com/c8121/asset-storage/internal/metadata"
	"github.com/c8121/asset-storage/internal/storage"
	"github.com/c8121/asset-storage/internal/util"
	"github.com/gin-gonic/gin"

	"golang.org/x/image/draw"
)

var (
	//go:embed icons/file-regular.svg
	DefaultFileThumbnail         []byte
	DefaultFileThumbnailMimeType = "image/svg+xml"

	ByMimeTypeThumbnailExt      = ".svg"
	ByMimeTypeThumbnailMimeType = "image/svg+xml"

	//go:embed icons/*.svg
	ByMimeTypeThumbnails embed.FS
)

func Thumbnail(c *gin.Context) {

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

	if bytes, err := loadIconByMimeType(meta.MimeType); err == nil {
		c.Data(http.StatusOK, ByMimeTypeThumbnailExt, bytes)
		return
	} else {
		util.LogError(err)
	}

	c.Data(http.StatusOK, DefaultFileThumbnailMimeType, DefaultFileThumbnail)

}

func generateThumbnail(assetHash string, meta metadata.AssetMetadata) ([]byte, string, error) {

	check := strings.ToLower(meta.MimeType)
	if !strings.HasPrefix(check, "image/") {
		return nil, "", fmt.Errorf("mime-type not supported: %s", meta.MimeType)
	}

	buf, err := storage.LoadByHash(assetHash)
	if err != nil {
		return nil, "", fmt.Errorf("failed to load asset: %w", err)
	}

	img, _, err := image.Decode(bytes.NewReader(buf))
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode asset: %w", err)
	}

	destSize := image.Rect(0, 0, 400, 200)
	thumb := image.NewRGBA(destSize)

	draw.NearestNeighbor.Scale(thumb, destSize, img, destSize, draw.Over, nil)

	var outBuf bytes.Buffer
	writer := bufio.NewWriter(&outBuf)
	if err := png.Encode(writer, thumb); err != nil {
		return nil, "", fmt.Errorf("failed to encode png: %w", err)
	}

	fmt.Printf("Created thumbnail (%d bytes, %s)\n", outBuf.Len(), destSize)
	return outBuf.Bytes(), "image/png", nil
}

// loadIconByMimeType returns icon matching mime-type: icons/type[-subtype].svg
func loadIconByMimeType(mimeType string) ([]byte, error) {

	thumb := strings.Replace(strings.ToLower(mimeType), "/", "-", 1) + ByMimeTypeThumbnailExt
	if bytes, err := ByMimeTypeThumbnails.ReadFile("icons/" + thumb); err == nil {
		return bytes, nil
	}

	thumb = strings.ToLower(mimeType)
	if p := strings.Index(thumb, "/"); p > 1 {
		thumb = thumb[:p] + ByMimeTypeThumbnailExt
	} else {
		thumb += ByMimeTypeThumbnailExt
	}

	if bytes, err := ByMimeTypeThumbnails.ReadFile("icons/" + thumb); err == nil {
		return bytes, nil
	} else {
		return nil, err
	}
}
