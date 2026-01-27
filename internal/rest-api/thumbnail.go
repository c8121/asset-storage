package restapi

import (
	"fmt"
	"strconv"

	"github.com/c8121/asset-storage/internal/filter"
	"github.com/c8121/asset-storage/internal/metadata"
)

var (
	ThumbnailWidth = 150
)

// generateThumbnail returns a thumbnail image generate from an asset.
// Returns content, mimeType, error
func generateThumbnail(assetHash string, meta *metadata.JsonAssetMetaData) ([]byte, string, error) {

	var f = filter.GetFirstFilterByMimeType(meta.MimeType)
	if f == nil {
		return nil, "", fmt.Errorf("no filter available for mime-type: %s", meta.MimeType)
	}

	params := map[string]string{}
	params["width"] = strconv.Itoa(ThumbnailWidth)

	return f.Apply(assetHash, meta, params)
}
