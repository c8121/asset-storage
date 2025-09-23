package restapi

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"

	"github.com/c8121/asset-storage/internal/metadata"
	"github.com/c8121/asset-storage/internal/storage"
	"golang.org/x/image/draw"
)

var (
	ThumbnailMimeType = "image/png"
)

// generateThumbnail returns a thumbnail image generate from an asset
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
	return outBuf.Bytes(), ThumbnailMimeType, nil
}
