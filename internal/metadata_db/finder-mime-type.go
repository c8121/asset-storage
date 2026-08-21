package metadata_db

import (
	"strconv"
	"strings"
	"time"
)

type FinderByMimeType struct {
}

// Find searches all assets having the given mime-type
func (f FinderByMimeType) Find(name any) (ScoredIdMap, error) {

	if len(name.(string)) == 0 {
		return nil, nil
	}

	var query = "SELECT a.id, a.fileTime FROM asset a " +
		"INNER JOIN mimeType m ON m.id = a.mimeType WHERE "

	mimeTypeId, err := strconv.Atoi(name.(string))
	if err == nil {
		name = mimeTypeId
		query += "(m.id = ?)"
	} else {
		name = strings.ReplaceAll(name.(string), "*", "%")
		query += "(m.name LIKE ?)"
	}

	return findAssetIds(func(id int64, match any, idMap *ScoredIdMap) {
		dt := match.(time.Time)
		score := float32(dt.Unix()) / float32(1000.0)
		idMap.Set(id, score)
	}, query, name)

}
