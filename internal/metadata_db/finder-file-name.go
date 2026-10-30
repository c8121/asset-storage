package metadata_db

import (
	"fmt"
	"strings"
	"time"

	"github.com/c8121/asset-storage/internal/metadata_db/entity"
	"github.com/c8121/go_db_entities"
)

type FinderByFileName struct {
}

// Find searches all assets having the given name as origin.FileName
func (f FinderByFileName) Find(name any) (ScoredIdMap, error) {

	var sName = name.(string)
	if len(sName) == 0 {
		return nil, nil
	}

	var findName = sName
	if strings.Contains(findName, "*") {
		//If an asterisk was given explicitly, only use this, just replace to sql-asterisk
		findName = strings.ReplaceAll(sName, "*", "%")
	} else {
		findName = "%" + strings.ReplaceAll(sName, " ", "%") + "%"
	}
	fmt.Printf("FinderByFileName: %s\n", findName)

	var query = go_db_entities.TableQuery{
		FieldNames: []string{},
		Entity:     &metadata_db_entity.FileName{},
		Filter:     go_db_entities.NewFilter("FileName.Name", go_db_entities.CompareContains, findName),
	}
	query.AddJoin(go_db_entities.TableQueryJoin{
		TableQuery: &go_db_entities.TableQuery{
			FieldNames: []string{"Id", "FileTime"},
			Entity:     &metadata_db_entity.Asset{},
		},
		OnFrom: "Id",
		OnTo:   "Name",
	})

	return findAssetIds(func(id int64, match any, idMap *ScoredIdMap) {
		dt := match.(time.Time)
		score := float32(dt.Unix()) / float32(1000.0)
		idMap.Set(id, score)
	}, query)

}
