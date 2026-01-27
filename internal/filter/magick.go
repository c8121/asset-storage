package filter

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/metadata"
	shell_command "github.com/c8121/asset-storage/internal/shell-command"
	"github.com/c8121/asset-storage/internal/storage"
	"github.com/c8121/asset-storage/internal/util"
)

type ImageMagickFilter struct {
	DefaultWidth           string
	DefaultFileNamePattern string
	DefaultMimeType        string
}

func NewImageMagickFilter() *ImageMagickFilter {
	f := &ImageMagickFilter{}
	f.DefaultWidth = "400"
	f.DefaultFileNamePattern = "asset-thumb*.png"
	f.DefaultMimeType = "image/png"
	return f
}

func (f ImageMagickFilter) Apply(assetHash string, meta *metadata.JsonAssetMetaData, params map[string]string) ([]byte, string, error) {

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

	checkMimeType := strings.ToLower(meta.MimeType)
	if strings.HasPrefix(checkMimeType, "application/pdf") {
		err = shell_command.ImageMagickThumbFromPdf(in, out.Name(), width, height)
	} else {
		err = shell_command.ImageMagickThumb(in, out.Name(), width, height)
	}

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
