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
	DefaultWidth      string
	ImageInterpolator draw.Interpolator
}

func NewImageFilter() *ImageFilter {
	f := &ImageFilter{}
	f.DefaultWidth = "400"
	f.ImageInterpolator = draw.BiLinear
	return f
}

func (f ImageFilter) Apply(assetHash string, meta *metadata.JsonAssetMetaData, params map[string]string) ([]byte, string, error) {

	width, _ := strconv.Atoi(util.GetOrDefault(params, "width", f.DefaultWidth))

	cropX1, _ := strconv.Atoi(util.GetOrDefault(params, "x1", "0"))
	cropY1, _ := strconv.Atoi(util.GetOrDefault(params, "y1", "0"))
	cropX2, _ := strconv.Atoi(util.GetOrDefault(params, "x2", "0"))
	cropY2, _ := strconv.Atoi(util.GetOrDefault(params, "y2", "0"))

	check := strings.ToLower(meta.MimeType)
	if !strings.HasPrefix(check, "image/") {
		return nil, "", fmt.Errorf("mime-type not supported: %s", meta.MimeType)
	}

	reader, err := storage.Open(assetHash)
	if err != nil {
		return nil, "", fmt.Errorf("failed to load asset: %w", err)
	}
	defer util.CloseOrLog(reader)

	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode asset: %w", err)
	}

	if cropX1 > 0 || cropY1 > 0 || cropX2 > 0 || cropY2 > 0 {
		img, err = cropImage(img, image.Rect(cropX1, cropY1, cropX2, cropY2))
		if err != nil {
			return nil, "", fmt.Errorf("failed to crop image: %w", err)
		}
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
		destSize := image.Rect(0, 0, scaleToWidth, scaleToHeight)
		thumb := image.NewRGBA(destSize)

		f.ImageInterpolator.Scale(thumb, destSize, img, img.Bounds(), draw.Over, nil)

		return encodePng(thumb)
	}

	return encodePng(img)
}

// cropImage takes an image and crops it to the specified rectangle.
func cropImage(img image.Image, crop image.Rectangle) (image.Image, error) {
	type subImager interface {
		SubImage(r image.Rectangle) image.Image
	}

	// img is an Image interface. This checks if the underlying value has a
	// method called SubImage. If it does, then we can use SubImage to crop the
	// image.
	subImage, ok := img.(subImager)
	if !ok {
		return nil, fmt.Errorf("image does not support cropping")
	}

	return subImage.SubImage(crop), nil
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
