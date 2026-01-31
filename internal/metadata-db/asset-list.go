package metadata_db

import (
	"fmt"
	"strings"
	"time"

	"github.com/c8121/asset-storage/internal/util"
)

type AssetListItem struct {
	Id       int64
	Hash     string
	Name     string
	MimeType string
	FileTime time.Time
}

type AssetListFilter struct {
	PathId   int64
	MimeType string
	FileName string
	PathName string
	Offset   int
	Count    int
}

func ListAssets(filter *AssetListFilter) ([]AssetListItem, error) {

	//Asset.Id -> Score
	var ids ScoredIdMap = nil

	//Finder -> value to use
	finders := map[Finder]any{
		FinderByPathId{}:   filter.PathId,
		FinderByMimeType{}: filter.MimeType,
		FinderByFileName{}: filter.FileName,
		FinderByPathName{}: filter.PathName,
	}

	for finder, value := range finders {
		foundIds, err := finder.Find(value)
		if err != nil {
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

	var query = "SELECT a.id, a.hash, m.name as mimeType, a.fileTime, f.name" +
		" FROM asset a " +
		" INNER JOIN mimeType m ON a.mimeType = m.id " +
		" INNER JOIN fileName f ON a.name = f.id "

	var params = make([]any, 0)

	if ids != nil {
		sorted := ids.Sort()
		endIdx := filter.Offset + filter.Count
		if endIdx >= len(sorted) {
			endIdx = len(sorted)
		}
		if filter.Offset < len(sorted) && endIdx > 0 {
			slice := sorted[filter.Offset:endIdx]
			//fmt.Printf("Found %d (%d - %d) items\n", len(slice), filter.Offset, endIdx)
			query += " WHERE a.id in(" +
				strings.Repeat("?,", len(slice)-1) + "?" +
				");"
			for _, id := range slice {
				//fmt.Printf("Found %d %f\n", id.Id, id.Score)
				params = append(params, id.Id)
			}
			items, err := loadAssetList(query, params...)
			if err != nil {
				return nil, err
			}
			//Retain sorting
			mapById := listToMap(items)
			list := make([]AssetListItem, 0, len(items))
			for _, id := range slice {
				item, ok := mapById[id.Id]
				if ok {
					list = append(list, item)
				}
			}
			return list, nil

		} else {
			//Nothing found
			empty := []AssetListItem{}
			return empty, nil
		}
	} else {
		//Nothing filtered
		query += "ORDER BY a.fileTime DESC, a.hash ASC LIMIT ? OFFSET ?;"
		params = append(params, filter.Count)
		params = append(params, filter.Offset)

		return loadAssetList(query, params...)
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
func loadAssetList(query string, params ...any) ([]AssetListItem, error) {

	fmt.Printf("Query: %s\n", query)
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer util.CloseOrLog(stmt)

	var items []AssetListItem

	if rows, err := stmt.Query(params...); err == nil {
		defer util.CloseOrLog(rows)
		for rows.Next() {
			var item AssetListItem
			if err := rows.Scan(&item.Id, &item.Hash, &item.MimeType, &item.FileTime, &item.Name); err != nil {
				return nil, err
			}
			items = append(items, item)
		}

	} else {
		return nil, err
	}

	return items, nil
}
