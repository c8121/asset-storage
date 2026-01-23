package restapi

import "github.com/c8121/asset-storage/internal/metadata"

// filter converts/filters/modified an asset an returns the filtered content
// Returns content, mimeType, error
func filter(assetHash string, meta *metadata.JsonAssetMetaData, filterName string, filterParams string) ([]byte, string, error) {

	b := []byte(assetHash + ": Not implemented yet: filters (name='" + filterName + "', params='" + filterParams + "')")
	return b, "text/plain", nil
}
