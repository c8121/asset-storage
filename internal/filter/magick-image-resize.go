package filter

import (
	"fmt"
	"os"
	"strconv"

	"github.com/c8121/asset-storage/internal/config"
	filter_commands "github.com/c8121/asset-storage/internal/filter-commands"
	"github.com/c8121/asset-storage/internal/metadata"
	"github.com/c8121/asset-storage/internal/storage"
	"github.com/c8121/asset-storage/internal/util"
)

type ImageMagickResizeFilter struct {
	DefaultWidth           string
	DefaultFileNamePattern string
	DefaultMimeType        string
}

func NewImageMagickResizeFilter() *ImageMagickResizeFilter {
	f := &ImageMagickResizeFilter{}
	f.DefaultWidth = "400"
	f.DefaultFileNamePattern = "asset-thumb*.png"
	f.DefaultMimeType = "image/png"
	return f
}

func (f ImageMagickResizeFilter) Apply(assetHash string, meta *metadata.JsonAssetMetaData, params map[string]string) ([]byte, string, error) {

	width, _ := strconv.Atoi(util.GetOrDefault(params, "width", f.DefaultWidth))
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

	err = imageMagickResize(in, out.Name(), width, height)
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

// imageMagickResize executes ImageMagick for Image conversion ...
func imageMagickResize(inputFilePath string, outputFilePath string, width int, height int) error {

	binary := filter_commands.FindImageMagickBin()
	if binary == "" {
		return fmt.Errorf("ImageMagick not found (searching in %v)", filter_commands.ImageMagickBinPaths)
	}

	var args []string

	args = append(args, inputFilePath)

	if width > 0 || height > 0 {
		args = append(args, "-geometry")
		if width > 0 && height > 0 {
			args = append(args, fmt.Sprintf("%d:%d", width, height))
		} else if width > 0 {
			args = append(args, fmt.Sprintf("%d", width))
		} else {
			args = append(args, fmt.Sprintf("x%d", height))
		}
	}

	args = append(args, "-flatten")
	args = append(args, "-colorspace", "RGB")

	args = append(args, outputFilePath)

	return util.RunSilent(binary, args...)
}
