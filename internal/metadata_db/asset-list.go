package metadata_db

import (
	"errors"
	"fmt"
	"time"

	"github.com/c8121/asset-storage/internal/metadata_db_entity"
	"github.com/c8121/go_db_entities"
)

type AssetListItem struct {
	Id       int64
	Hash     string
	Name     string
	MimeType string
	FileTime time.Time
	PathIds  []int64
}

type AssetListFilter struct {
	PathId     int64
	MimeType   string
	FileName   string
	PathName   string
	Collection string
	Fulltext   string
	Face       int64
	Offset     int
	Count      int
}

func ListAssets(filter *AssetListFilter) ([]AssetListItem, error) {

	//Asset.Id -> Score
	var ids ScoredIdMap = nil

	//Finder -> value to use
	finders := map[Finder]any{
		FinderByPathId{}:     filter.PathId,
		FinderByMimeType{}:   filter.MimeType,
		FinderByFileName{}:   filter.FileName,
		FinderByPathName{}:   filter.PathName,
		FinderByFace{}:       filter.Face,
		FinderByCollection{}: filter.Collection,
		FinderByFulltext{}:   filter.Fulltext,
	}

	for finder, value := range finders {
		foundIds, err := finder.Find(value)
		if err != nil && !errors.Is(go_db_entities.ErrNotFound, err) {
			return nil, err
		}
		if foundIds != nil {
			fmt.Printf("Found %d assets using %T with '%v'\n", len(foundIds), finder, value)
			if ids != nil {
				ids.Reduce(foundIds)
			} else {
				ids = foundIds
			}
		}
	}

	return ListAssetsByIds(ids, filter.Offset, filter.Count)
}

// ListAssetsByIds creates a list of AssetListItem's, loaded by given ScoredIdMap
func ListAssetsByIds(ids ScoredIdMap, offset int, count int) ([]AssetListItem, error) {

	var query = go_db_entities.TableQuery{
		FieldNames: []string{"Id", "Hash", "FileTime"},
		Entity:     &metadata_db_entity.Asset{},
	}
	query.AddJoin(go_db_entities.TableQueryJoin{
		TableQuery: &go_db_entities.TableQuery{
			FieldNames: []string{"Name"},
			Entity:     &metadata_db_entity.MimeType{},
		},
		OnFrom: "MimeType",
		OnTo:   "Id",
	})
	query.AddJoin(go_db_entities.TableQueryJoin{
		TableQuery: &go_db_entities.TableQuery{
			FieldNames: []string{"Name"},
			Entity:     &metadata_db_entity.FileName{},
		},
		OnFrom: "Name",
		OnTo:   "Id",
	})

	query.Filter = &go_db_entities.Filter{}

	if ids != nil {
		sorted := ids.Sort()
		endIdx := offset + count
		if endIdx >= len(sorted) {
			endIdx = len(sorted)
		}
		if offset < len(sorted) && endIdx > 0 {
			slice := sorted[offset:endIdx]

			query.Filter.FieldName = "Asset.Id"
			ids := make([]int64, len(slice))
			for i, id := range slice {
				ids[i] = id.Id
			}
			query.Filter.Value = ids

			items, err := loadAssetList(query)
			if err != nil {
				return nil, err
			}
			//Retain sorting
			mapById := listToMap(items)
			list := make([]AssetListItem, 0, len(items))
			for _, id := range sorted[offset:endIdx] {
				item, ok := mapById[id.Id]
				if ok {
					list = append(list, item)
				}
			}

			if err = assignPathIds(&list); err != nil {
				return nil, err
			}
			return list, nil

		} else {
			//Nothing found
			var empty []AssetListItem
			return empty, nil
		}
	} else {
		//Nothing filtered
		query.Filter.Offset = offset
		query.Filter.Count = count
		query.Filter.OrderBy = []string{"Asset.FileTime DESC, Asset.Id DESC"}

		list, err := loadAssetList(query)
		if err != nil {
			return nil, err
		}
		if err = assignPathIds(&list); err != nil {
			return nil, err
		}
		return list, nil
	}
}

// listToMap converts list of AssetListItem to a map with ID as key.
func listToMap(items []AssetListItem) map[int64]AssetListItem {
	mapById := make(map[int64]AssetListItem, len(items))
	for _, item := range items {
		mapById[item.Id] = item
	}
	return mapById
}

// loadAssetList queries the database
func loadAssetList(query go_db_entities.TableQuery) ([]AssetListItem, error) {

	rows, err := go_db_entities.Query(metadata_db_entity.DbContext, query)
	if err != nil {
		return nil, err
	}

	var items []AssetListItem

	for _, row := range rows {
		item := AssetListItem{
			Id:       row[0].(int64),
			Hash:     row[1].(string),
			FileTime: row[2].(time.Time),
			MimeType: row[3].(string),
			Name:     row[4].(string),
		}
		items = append(items, item)
	}

	return items, nil
}

// assignPathIds uses FindPathIds and assigns to AssetListItem.PathIds
func assignPathIds(items *[]AssetListItem) error {

	assetIds := make([]int64, len(*items))
	for i, item := range *items {
		assetIds[i] = item.Id
	}

	idMap, err := metadata_db_entity.FindAssetPathIds(assetIds)
	if err != nil {
		return err
	}

	for i := range *items {
		src := idMap[(*items)[i].Id]
		(*items)[i].PathIds = make([]int64, len(src))
		copy((*items)[i].PathIds, src)
	}

	return nil
}
