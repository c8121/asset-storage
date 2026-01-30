package metadata_db

import "time"

type FinderByPathId struct {
}

// Find searches all assets assigned to given path id
func (f FinderByPathId) Find(pathId any) (ScoredIdMap, error) {

	if pathId.(int64) == 0 {
		return nil, nil
	}

	var query = "SELECT a.id, a.fileTime FROM origin o " +
		"INNER JOIN asset a ON o.asset = a.id " +
		"WHERE path = ?;"

	return findAssetIds(query, pathId, func(id int64, match any, idMap *ScoredIdMap) {
		dt := match.(time.Time)
		score := float32(dt.Unix()) / float32(1000.0)
		idMap.Set(id, score)
	})

}
