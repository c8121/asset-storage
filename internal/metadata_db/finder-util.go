package metadata_db

import (
	"github.com/c8121/asset-storage/internal/metadata_db_entity"
	"github.com/c8121/go_db_entities"
)

func findAssetIds(calcScore func(id int64, match any, idMap *ScoredIdMap), query go_db_entities.TableQuery) (ScoredIdMap, error) {

	rows, err := go_db_entities.Query(metadata_db_entity.DbContext, query)
	if err != nil {
		return nil, err
	}

	ids := make(ScoredIdMap)

	for _, row := range rows {
		id := row[0].(int64)
		match := row[1]
		calcScore(id, match, &ids)
	}

	return ids, nil
}
