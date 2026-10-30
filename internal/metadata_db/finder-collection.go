package metadata_db

import (
	"time"

	"github.com/c8121/asset-storage/internal/collections"
	"github.com/c8121/asset-storage/internal/metadata_db/entity"
	"github.com/c8121/go_db_entities"
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

	var query = go_db_entities.TableQuery{
		FieldNames: []string{"Id", "FileTime"},
		Entity:     &metadata_db_entity.Asset{},
		Filter:     go_db_entities.NewFilter("Hash", go_db_entities.CompareIn, collection.Assets),
	}

	return findAssetIds(func(id int64, match any, idMap *ScoredIdMap) {
		dt := match.(time.Time)
		score := float32(dt.Unix()) / float32(1000.0)
		idMap.Set(id, score)
	}, query)
}
