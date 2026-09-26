package filter

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/filter_commands"
	"github.com/c8121/asset-storage/internal/metadata"
	"github.com/c8121/asset-storage/internal/storage"
	"github.com/c8121/asset-storage/internal/util"
)

type ImageMagickDocToImageFilter struct {
	DefaultWidth           int
	DefaultFileNamePattern string
	DefaultMimeType        string
}

func NewImageMagickDocToImageFilter() *ImageMagickDocToImageFilter {
	f := &ImageMagickDocToImageFilter{}
	f.DefaultWidth = 400
	f.DefaultFileNamePattern = "asset-thumb*.png"
	f.DefaultMimeType = "image/png"
	return f
}

func (f ImageMagickDocToImageFilter) Apply(assetHash string, meta *metadata.JsonAssetMetaData, params map[string]string) ([]byte, string, error) {

	check := strings.ToLower(meta.MimeType)
	if !strings.Contains(check, "document") && !strings.Contains(check, "application/vnd.ms-") {
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
		return nil, "", fmt.Errorf("failed to create temp file: %w", err)
	}
	util.LogError(out.Close())

	var fileExtension string
	for _, origin := range meta.Origins {
		if origin.Name != "" {
			fileExtension = strings.ToLower(filepath.Ext(origin.Name))
			if strings.HasPrefix(fileExtension, ".") {
				fileExtension = fileExtension[1:]
			}
			break
		}
	}

	if fileExtension == "" {
		fileExtension = "doc"
	}

	err = f.imageMagickDocToImage(fileExtension+":"+in, out.Name(), width, height)
	if err != nil {
		util.LogError(os.Remove(out.Name()))
		return nil, "", fmt.Errorf("failed to create thumbnail: %w", err)
	}

	bytes, err := os.ReadFile(out.Name())
	if err != nil {
		util.LogError(os.Remove(out.Name()))
		return nil, "", fmt.Errorf("failed to read thumbnail: %w", err)
	}

	util.LogError(os.Remove(out.Name()))
	return bytes, mimeType, nil
}

// imageMagickDocToImage executes ImageMagick for Document to Image conversion ...
func (f ImageMagickDocToImageFilter) imageMagickDocToImage(input string, output string, width int, height int) error {

	binary := filter_commands.FindImageMagickBin()
	if binary == "" {
		return fmt.Errorf("ImageMagick not found (searching in %v)", filter_commands.ImageMagickBinPaths)
	}

	var args []string

	args = append(args, input+"[0]")

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
