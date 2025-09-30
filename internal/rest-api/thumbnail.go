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
	ThumbnailWidth    = 400

	//ThumbnailInterpolator = draw.NearestNeighbor
	//ThumbnailInterpolator = draw.ApproxBiLinear
	ThumbnailInterpolator = draw.BiLinear
	//ThumbnailInterpolator = draw.CatmullRom
)

// generateThumbnail returns a thumbnail image generate from an asset
func generateThumbnail(assetHash string, meta metadata.AssetMetadata) ([]byte, string, error) {

	check := strings.ToLower(meta.MimeType)
	if !strings.HasPrefix(check, "image/") {
		return nil, "", fmt.Errorf("mime-type not supported: %s", meta.MimeType)
	}

	reader, err := storage.Open(assetHash)
	if err != nil {
		return nil, "", fmt.Errorf("failed to load asset: %w", err)
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode asset: %w", err)
	}

	imgWidth := img.Bounds().Dx()
	var scaleToWidth int
	switch {
	case imgWidth < ThumbnailWidth:
		scaleToWidth = imgWidth
	default:
		scaleToWidth = ThumbnailWidth
	}

	var outBuf bytes.Buffer
	if imgWidth > scaleToWidth {

		scaleToHeight := int(float64(img.Bounds().Dy()) * (float64(scaleToWidth) / float64(imgWidth)))
		fmt.Printf("scale to width: %d, height: %d (%d * (%d / %d))\n", scaleToWidth, scaleToHeight,
			img.Bounds().Dy(), scaleToWidth, imgWidth)

		destSize := image.Rect(0, 0, scaleToWidth, scaleToHeight)
		thumb := image.NewRGBA(destSize)

		ThumbnailInterpolator.Scale(thumb, destSize, img, img.Bounds(), draw.Over, nil)

		writer := bufio.NewWriter(&outBuf)
		if err := png.Encode(writer, thumb); err != nil || outBuf.Len() == 0 {
			return nil, "", fmt.Errorf("failed to encode png: %v", err)
		}

	} else {
		writer := bufio.NewWriter(&outBuf)
		if err := png.Encode(writer, img); err != nil || outBuf.Len() == 0 {
			return nil, "", fmt.Errorf("failed to encode png: %v", err)
		}
	}

	//fmt.Printf("Created thumbnail (%d bytes, %s)\n", outBuf.Len(), destSize)
	return outBuf.Bytes(), ThumbnailMimeType, nil
}
