package restapi

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/c8121/asset-storage/internal/filter"
	"github.com/c8121/asset-storage/internal/metadata"
)

var (
	ThumbnailWidth = 150
)

// generateThumbnail returns a thumbnail image generate from an asset.
// Returns content, mimeType, error
func generateThumbnail(assetHash string, meta *metadata.JsonAssetMetaData) ([]byte, string, error) {

	//TODO make mime-type->converter mapping in configurable

	params := map[string]string{}
	params["width"] = strconv.Itoa(ThumbnailWidth)

	var f filter.Filter

	check := strings.ToLower(meta.MimeType)
	if check == "image/bmp" || check == "image/tiff" {
		f = filter.NewImageMagickFilter()
	} else if strings.HasPrefix(check, "application/pdf") {
		f = filter.NewImageMagickFilter()
	} else if strings.Contains(check, "photoshop") {
		f = filter.NewImageMagickFilter()
	} else if strings.HasPrefix(check, "image/") {
		f = filter.NewImageFilter()
	} else if strings.HasPrefix(check, "video/") {
		f = filter.NewFFmpegFilter()
	} else {
		return nil, "", fmt.Errorf("mime-type not supported: %s", meta.MimeType)
	}

	return f.Apply(assetHash, meta, params)
}
