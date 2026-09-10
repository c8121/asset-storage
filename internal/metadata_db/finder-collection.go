package metadata_db

import (
	"strings"
	"time"

	"github.com/c8121/asset-storage/internal/collections"
)

type FinderByCollection struct {
}

// Find searches all assets having the given face
func (f FinderByCollection) Find(uuid any) (ScoredIdMap, error) {

	var sUuid = uuid.(string)
	if len(sUuid) == 0 {
		return nil, nil
	}

	collection, err := collections.LoadById(sUuid)
	if err != nil {
		return nil, err
	}

	var query = "SELECT a.id, a.fileTime FROM asset a " +
		" WHERE a.hash in(" +
		strings.Repeat("?,", len(collection.Assets)-1) + "?" +
		");"

	args := make([]any, len(collection.Assets))
	for i, asset := range collection.Assets {
		args[i] = asset
	}

	return findAssetIds(func(id int64, match any, idMap *ScoredIdMap) {
		dt := match.(time.Time)
		score := float32(dt.Unix()) / float32(1000.0)
		idMap.Set(id, score)
	}, query, args...)
}
