package filter

import "github.com/c8121/asset-storage/internal/metadata"

type Filter interface {
	//Apply executes the filtering, returns content, mime-type, error
	Apply(assetHash string, meta *metadata.JsonAssetMetaData, params map[string]string) ([]byte, string, error)
}
