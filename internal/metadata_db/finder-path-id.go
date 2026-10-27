package metadata_db

import (
	"time"

	"github.com/c8121/asset-storage/internal/metadata_db_entity"
	"github.com/c8121/go_db_entities"
)

type FinderByPathId struct {
}

// Find searches all assets assigned to given path id
func (f FinderByPathId) Find(pathId any) (ScoredIdMap, error) {

	if pathId.(int64) == 0 {
		return nil, nil
	}

	var query = go_db_entities.TableQuery{
		FieldNames: []string{},
		Entity:     &metadata_db_entity.Origin{},
		Filter:     go_db_entities.NewFilter("Path", go_db_entities.CompareEqual, pathId),
	}
	query.AddJoin(go_db_entities.TableQueryJoin{
		TableQuery: &go_db_entities.TableQuery{
			FieldNames: []string{"Id", "FileTime"},
			Entity:     &metadata_db_entity.Asset{},
		},
		OnFrom: "Asset",
		OnTo:   "Id",
	})

	return findAssetIds(func(id int64, match any, idMap *ScoredIdMap) {
		dt := match.(time.Time)
		score := float32(dt.Unix()) / float32(1000.0)
		idMap.Set(id, score)
	}, query)

}
