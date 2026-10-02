package filter

import (
	"fmt"
	"io"
	"os"
	"strings"

	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"

	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/filter_commands"
	"github.com/c8121/asset-storage/internal/metadata"
	"github.com/c8121/asset-storage/internal/storage"
	"github.com/c8121/asset-storage/internal/util"
	"github.com/disintegration/imaging"
	"golang.org/x/image/draw"

	_ "github.com/HugoSmits86/nativewebp"
)

type OcrImagePrepareFilter struct {
	DefaultWidth      string
	ImageInterpolator draw.Interpolator
}

func NewOcrImagePrepareFilter() *OcrImagePrepareFilter {
	f := &OcrImagePrepareFilter{}
	f.DefaultWidth = "400"
	f.ImageInterpolator = draw.BiLinear
	return f
}

func (f OcrImagePrepareFilter) Apply(assetHash string, meta *metadata.JsonAssetMetaData, params map[string]string) ([]byte, string, error) {

	tmpInputFile := ""
	var reader io.ReadCloser
	var err error

	check := strings.ToLower(meta.MimeType)
	if strings.HasPrefix(check, "application/pdf") {

		//Convert PDF to image

		out, err := os.CreateTemp(config.AssetStorageTempDir, "asset-ocr*.png")
		if err != nil {
			return nil, "", fmt.Errorf("failed to create temp file: %w", err)
		}
		util.LogError(out.Close())
		tmpInputFile = out.Name()

		path, err := storage.FindByHash(assetHash)
		if err != nil {
			return nil, "", fmt.Errorf("failed to load asset: %w", err)
		}

		err = imageMagickPdfToOcrImage(path, tmpInputFile)
		if err != nil {
			return nil, "", fmt.Errorf("failed to convert pdf: %w", err)
		}

		reader, err = os.Open(tmpInputFile)
		if err != nil {
			return nil, "", fmt.Errorf("failed to open converted file: %w", err)
		}
	} else {
		reader, err = storage.Open(assetHash)
		if err != nil {
			return nil, "", fmt.Errorf("failed to load asset: %w", err)
		}
	}

	defer util.CloseOrLog(reader)

	src, err := imaging.Decode(reader)
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode asset: %w", err)
	}

	// 2. Grayscale & Rescale (Double the width/height cleanly via CatmullRom)
	img := imaging.Grayscale(src)
	img = imaging.Resize(img, img.Bounds().Dx()*2, 0, imaging.CatmullRom)

	// 3. Simple Binarization (Thresholding)
	// For complex shading, an adaptive library is ideal, but a simple 128 mid-point
	// loop works perfectly for clean pages without using OpenCV.
	bounds := img.Bounds()
	binarized := image.NewGray(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			oldColor := img.At(x, y)
			grayColor := color.GrayModel.Convert(oldColor).(color.Gray)
			if grayColor.Y > 128 { // Threshold cutoff
				binarized.SetGray(x, y, color.Gray{Y: 255}) // White
			} else {
				binarized.SetGray(x, y, color.Gray{Y: 0}) // Black
			}
		}
	}

	// 4. Add White Border padding
	canvasWidth := binarized.Bounds().Dx() + 40
	canvasHeight := binarized.Bounds().Dy() + 40
	backgroundCanvas := imaging.New(canvasWidth, canvasHeight, color.White)

	imgFinal := imaging.Paste(backgroundCanvas, binarized, image.Pt(20, 20))

	if tmpInputFile != "" {
		util.LogError(os.Remove(tmpInputFile))
	}

	return encodePng(imgFinal)

}

// TODO currently converts first page only
// imageMagickPdfToOcrImage executes ImageMagick for PDF to Image conversion with parameters optimized for OCR ...
func imageMagickPdfToOcrImage(input string, output string) error {

	binary := filter_commands.FindImageMagickBin()
	if binary == "" {
		return fmt.Errorf("ImageMagick not found (searching in %v)", filter_commands.ImageMagickBinPaths)
	}

	var args []string

	args = append(args, "-density", "600") //must be before input

	args = append(args, input+"[0]")

	args = append(args, "-thumbnail")
	args = append(args, "1200x")

	args = append(args, "-flatten")
	args = append(args, "-colorspace", "RGB")

	args = append(args, output)

	return util.RunSilent(binary, args...)
}
