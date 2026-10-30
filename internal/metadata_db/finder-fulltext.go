package metadata_db

import (
	"time"

	"github.com/c8121/asset-storage/internal/fulltext"
	"github.com/c8121/asset-storage/internal/metadata_db/entity"
	"github.com/c8121/go_db_entities"
)

type FinderByFulltext struct {
}

// Find searches all assets having the given face
func (f FinderByFulltext) Find(fulltextQuery any) (ScoredIdMap, error) {

	var sQuery = fulltextQuery.(string)
	if len(sQuery) == 0 {
		return nil, nil
	}

	hashes, err := fulltext.Find(sQuery, 0, 999)
	if err != nil {
		return nil, err
	}

	if len(hashes) == 0 {
		empty := make(ScoredIdMap)
		return empty, nil
	}

	var query = go_db_entities.TableQuery{
		FieldNames: []string{"Id", "FileTime"},
		Entity:     &metadata_db_entity.Asset{},
		Filter:     go_db_entities.NewFilter("Hash", go_db_entities.CompareIn, hashes),
	}

	return findAssetIds(func(id int64, match any, idMap *ScoredIdMap) {
		dt := match.(time.Time)
		score := float32(dt.Unix()) / float32(1000.0)
		idMap.Set(id, score)
	}, query)
}
