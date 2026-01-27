package restapi

import (
	"fmt"

	"github.com/c8121/asset-storage/internal/filter"
	"github.com/c8121/asset-storage/internal/metadata"
)

// filterAsset converts/filters/modified an asset an returns the filtered content
// Returns content, mimeType, error
func filterAsset(assetHash string, meta *metadata.JsonAssetMetaData, filterName string, filterParams map[string]string) ([]byte, string, error) {

	var f filter.Filter

	if filterName == "image" {
		f = filter.NewImageMagickFilter()
	} else {
		return nil, "", fmt.Errorf("filter not found: %s", filterName)
	}

	return f.Apply(assetHash, meta, filterParams)

}
