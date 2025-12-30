package restapi

import (
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	"os"
	"strings"

	_ "github.com/HugoSmits86/nativewebp"
	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/metadata"
	"github.com/c8121/asset-storage/internal/storage"
	"github.com/c8121/asset-storage/internal/thumbnails"
	"github.com/c8121/asset-storage/internal/util"
)

func generateThumbnailWithImageMagick(assetHash string, meta *metadata.JsonAssetMetaData) ([]byte, string, error) {

	in, err := storage.FindByHash(assetHash)
	if err != nil {
		return nil, "", fmt.Errorf("cannot find asset: %w", err)
	}

	out, err := os.CreateTemp(config.AssetStorageTempDir, "asset-thumb*.png")
	if err != nil {
		return nil, "", fmt.Errorf("Failed to create temp file: %w", err)
	}
	util.LogError(out.Close())

	mimeType := strings.ToLower(meta.MimeType)
	if strings.HasPrefix(mimeType, "application/pdf") {
		err = thumbnails.ImageMagickThumbFromPdf(in, out.Name(), ThumbnailWidth, -1)
	} else {
		err = thumbnails.ImageMagickThumb(in, out.Name(), ThumbnailWidth, -1)
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
	return bytes, ThumbnailMimeType, nil
}
