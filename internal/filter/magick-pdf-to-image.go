package filter

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/c8121/asset-storage/internal/config"
	filter_commands "github.com/c8121/asset-storage/internal/filter-commands"
	"github.com/c8121/asset-storage/internal/metadata"
	"github.com/c8121/asset-storage/internal/storage"
	"github.com/c8121/asset-storage/internal/util"
)

type ImageMagickPdfToImageFilter struct {
	DefaultWidth           int
	DefaultFileNamePattern string
	DefaultMimeType        string
}

func NewImageMagickPdfToImageFilter() *ImageMagickPdfToImageFilter {
	f := &ImageMagickPdfToImageFilter{}
	f.DefaultWidth = 400
	f.DefaultFileNamePattern = "asset-thumb*.png"
	f.DefaultMimeType = "image/png"
	return f
}

func (f ImageMagickPdfToImageFilter) Apply(assetHash string, meta *metadata.JsonAssetMetaData, params map[string]string) ([]byte, string, error) {

	check := strings.ToLower(meta.MimeType)
	if !strings.Contains(check, "pdf") {
		return nil, "", fmt.Errorf("mime-type not supported: %s", meta.MimeType)
	}

	width, _ := strconv.Atoi(util.GetOrDefault(params, "width", strconv.Itoa(f.DefaultWidth)))
	height, _ := strconv.Atoi(util.GetOrDefault(params, "height", "0"))
	tempFileNamePattern := util.GetOrDefault(params, "fileNamePattern", f.DefaultFileNamePattern)
	mimeType := util.GetOrDefault(params, "mimeType", f.DefaultMimeType)

	in, err := storage.FindByHash(assetHash)
	if err != nil {
		return nil, "", fmt.Errorf("cannot find asset: %w", err)
	}

	out, err := os.CreateTemp(config.AssetStorageTempDir, tempFileNamePattern)
	if err != nil {
		return nil, "", fmt.Errorf("Failed to create temp file: %w", err)
	}
	util.LogError(out.Close())

	err = f.imageMagickPdfToImage(in, out.Name(), width, height)
	if err != nil {
		util.LogError(os.Remove(out.Name()))
		return nil, "", fmt.Errorf("Failed to create thumbnail: %w", err)
	}

	bytes, err := os.ReadFile(out.Name())
	if err != nil {
		util.LogError(os.Remove(out.Name()))
		return nil, "", fmt.Errorf("Failed to read thumbnail: %w", err)
	}

	util.LogError(os.Remove(out.Name()))
	return bytes, mimeType, nil
}

// imageMagickPdfToImage executes ImageMagick for PDF to Image conversion ...
func (f ImageMagickPdfToImageFilter) imageMagickPdfToImage(input string, output string, width int, height int) error {

	binary := filter_commands.FindImageMagickBin()
	if binary == "" {
		return fmt.Errorf("ImageMagick not found (searching in %v)", filter_commands.ImageMagickBinPaths)
	}

	var args []string

	args = append(args, input+"[0]")

	if width > 0 || height > 0 {
		args = append(args, "-thumbnail")
		args = append(args, fmt.Sprintf("%sx%s", util.Iif(width > 0, strconv.Itoa(width), ""), util.Iif(height > 0, strconv.Itoa(height), "")))
	}
	args = append(args, "-thumbnail")
	if width > 0 && height > 0 {
		args = append(args, fmt.Sprintf("%dx%d", width, height))
	} else if width > 0 {
		args = append(args, fmt.Sprintf("%dx", width))
	} else if height > 0 {
		args = append(args, fmt.Sprintf("x%d", height))
	} else {
		args = append(args, fmt.Sprintf("%dx", f.DefaultWidth))
	}

	args = append(args, "-flatten")
	args = append(args, "-colorspace", "RGB")

	args = append(args, output)

	return util.RunSilent(binary, args...)
}
