package metadata_db

import (
	"strings"
	"time"

	"github.com/c8121/asset-storage/internal/fulltext"
)

type FinderByFulltext struct {
}

// Find searches all assets having the given face
func (f FinderByFulltext) Find(fulltextQuery any) (ScoredIdMap, error) {

	var sQuery = fulltextQuery.(string)
	if len(sQuery) == 0 {
		return nil, nil
	}

	hashes, err := fulltext.Find(sQuery)
	if err != nil {
		return nil, err
	}

	if len(hashes) == 0 {
		empty := make(ScoredIdMap)
		return empty, nil
	}

	var query = "SELECT a.id, a.fileTime FROM asset a " +
		" WHERE a.hash in(" +
		strings.Repeat("?,", len(hashes)-1) + "?" +
		");"

	args := make([]any, len(hashes))
	for i, hash := range hashes {
		args[i] = hash
	}

	return findAssetIds(func(id int64, match any, idMap *ScoredIdMap) {
		dt := match.(time.Time)
		score := float32(dt.Unix()) / float32(1000.0)
		idMap.Set(id, score)
	}, query, args...)
}
