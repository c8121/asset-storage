package metadata_db

import (
	"strconv"
	"strings"
	"time"

	"github.com/c8121/asset-storage/internal/metadata_db_entity"
	"github.com/c8121/go_db_entities"
)

type FinderByMimeType struct {
}

// Find searches all assets having the given mime-type
func (f FinderByMimeType) Find(name any) (ScoredIdMap, error) {

	if len(name.(string)) == 0 {
		return nil, nil
	}

	filter := &go_db_entities.Filter{}

	mimeTypeId, err := strconv.Atoi(name.(string))
	if err == nil {
		filter.FieldName = "MimeType.Id"
		filter.Value = mimeTypeId
	} else {
		filter.FieldName = "MimeType.Name"
		filter.Comparison = go_db_entities.CompareContains
		filter.Value = strings.ReplaceAll(name.(string), "*", "%")
	}

	var query = go_db_entities.TableQuery{
		FieldNames: []string{},
		Entity:     &metadata_db_entity.MimeType{},
		Filter:     filter,
	}
	query.AddJoin(go_db_entities.TableQueryJoin{
		TableQuery: &go_db_entities.TableQuery{
			FieldNames: []string{"Id", "FileTime"},
			Entity:     &metadata_db_entity.Asset{},
		},
		OnFrom: "Id",
		OnTo:   "MimeType",
	})

	return findAssetIds(func(id int64, match any, idMap *ScoredIdMap) {
		dt := match.(time.Time)
		score := float32(dt.Unix()) / float32(1000.0)
		idMap.Set(id, score)
	}, query)

}
