package restapi

import (
	"fmt"
	"os"
	"strings"

	_ "image/gif"
	_ "image/jpeg"

	_ "github.com/HugoSmits86/nativewebp"
	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/metadata"
	"github.com/c8121/asset-storage/internal/storage"
	"github.com/c8121/asset-storage/internal/thumbnails"
	"github.com/c8121/asset-storage/internal/util"
)

func generateThumbnailFromVideo(assetHash string, meta *metadata.JsonAssetMetaData) ([]byte, string, error) {

	check := strings.ToLower(meta.MimeType)
	if !strings.HasPrefix(check, "video/") {
		return nil, "", fmt.Errorf("mime-type not supported: %s", meta.MimeType)
	}

	in, err := storage.FindByHash(assetHash)
	if err != nil {
		return nil, "", fmt.Errorf("cannot find asset: %w", err)
	}

	out, err := os.CreateTemp(config.AssetStorageTempDir, "asset-thumb*.png")
	if err != nil {
		return nil, "", fmt.Errorf("Failed to create temp file: %w", err)
	}
	util.LogError(out.Close())

	err = thumbnails.FFmpegThumb(in, out.Name(), ThumbnailWidth, -1)
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
