package filter

import (
	"fmt"
	"strconv"
	"strings"

	"bufio"
	"bytes"

	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"

	"github.com/c8121/asset-storage/internal/metadata"
	"github.com/c8121/asset-storage/internal/storage"
	"github.com/c8121/asset-storage/internal/util"
	"golang.org/x/image/draw"

	_ "github.com/HugoSmits86/nativewebp"
)

type ImageFilter struct {
	DefaultWidth          string
	ThumbnailInterpolator draw.Interpolator
}

func NewImageFilter() *ImageFilter {
	f := &ImageFilter{}
	f.DefaultWidth = "400"
	f.ThumbnailInterpolator = draw.BiLinear
	return f
}

func (f ImageFilter) Apply(assetHash string, meta *metadata.JsonAssetMetaData, params map[string]string) ([]byte, string, error) {

	width, _ := strconv.Atoi(util.GetOrDefault(params, "width", f.DefaultWidth))

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
	case imgWidth < width:
		scaleToWidth = imgWidth
	default:
		scaleToWidth = width
	}

	if imgWidth > scaleToWidth {

		scaleToHeight := int(float64(img.Bounds().Dy()) * (float64(scaleToWidth) / float64(imgWidth)))
		fmt.Printf("scale to width: %d, height: %d (%d * (%d / %d))\n", scaleToWidth, scaleToHeight,
			img.Bounds().Dy(), scaleToWidth, imgWidth)

		destSize := image.Rect(0, 0, scaleToWidth, scaleToHeight)
		thumb := image.NewRGBA(destSize)

		f.ThumbnailInterpolator.Scale(thumb, destSize, img, img.Bounds(), draw.Over, nil)

		return encodePng(thumb)
	} else {
		return encodePng(img)
	}
}

func encodePng(img image.Image) ([]byte, string, error) {

	var outBuf bytes.Buffer

	writer := bufio.NewWriter(&outBuf)
	if err := png.Encode(writer, img); err != nil {
		return nil, "", fmt.Errorf("failed to encode png: %v", err)
	}
	util.LogError(writer.Flush())
	if outBuf.Len() == 0 {
		return nil, "", fmt.Errorf("failed to encode png, no bytes written")
	}

	return outBuf.Bytes(), "image/png", nil
}
