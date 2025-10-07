package restapi

import (
	"fmt"
	"strings"

	_ "image/gif"
	_ "image/jpeg"

	_ "github.com/HugoSmits86/nativewebp"
	"github.com/c8121/asset-storage/internal/metadata"
)

var (
	ThumbnailWidth = 150
)

// generateThumbnail returns a thumbnail image generate from an asset
func generateThumbnail(assetHash string, meta *metadata.JsonAssetMetaData) ([]byte, string, error) {

	check := strings.ToLower(meta.MimeType)
	if strings.HasPrefix(check, "image/") {
		return generateThumbnailFromImage(assetHash, meta)
	} else if strings.HasPrefix(check, "video/") {
		return generateThumbnailFromVideo(assetHash, meta)
	} else {
		return nil, "", fmt.Errorf("mime-type not supported: %s", meta.MimeType)
	}
}
