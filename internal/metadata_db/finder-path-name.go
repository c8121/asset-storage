package metadata_db

import (
	"fmt"
	"strings"
	"time"

	"github.com/c8121/asset-storage/internal/metadata_db/entity"
	"github.com/c8121/go_db_entities"
)

type FinderByPathName struct {
}

// Find searches all assets assigned to given path id
func (f FinderByPathName) Find(name any) (ScoredIdMap, error) {

	var sName = name.(string)
	if len(sName) == 0 {
		return nil, nil
	}

	//TODO search parents

	var findName = sName
	if strings.Contains(findName, "*") {
		//If an asterisk was given explicitly, only use this, just replace to sql-asterisk
		findName = strings.ReplaceAll(sName, "*", "%")
	} else {
		findName = "%" + strings.ReplaceAll(sName, " ", "%") + "%"
	}
	fmt.Printf("FinderByPathName: %s\n", findName)

	var originQuery = go_db_entities.TableQuery{
		FieldNames: []string{},
		Entity:     &metadata_db_entity.Origin{},
		Filter:     go_db_entities.NewFilter("PathItem.Name", go_db_entities.CompareContains, findName),
	}
	var pathItemQuery = go_db_entities.TableQuery{
		FieldNames: []string{},
		Entity:     &metadata_db_entity.PathItem{},
	}
	originQuery.AddJoin(go_db_entities.TableQueryJoin{TableQuery: &pathItemQuery, OnFrom: "Path", OnTo: "Id"})

	var assetQuery = go_db_entities.TableQuery{
		FieldNames: []string{"Id", "FileTime"},
		Entity:     &metadata_db_entity.Asset{},
	}
	assetQuery.AddJoin(go_db_entities.TableQueryJoin{TableQuery: &originQuery, OnFrom: "Id", OnTo: "Asset"})

	return findAssetIds(func(id int64, match any, idMap *ScoredIdMap) {
		dt := match.(time.Time)
		score := float32(dt.Unix()) / float32(1000.0)
		idMap.Set(id, score)
	}, assetQuery)

}
