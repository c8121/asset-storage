package metadata_db

import (
	"fmt"
	"strconv"
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

	//Find value -> function to use
	filterFunctions := map[any]interface{}{
		filter.PathId:   findAssetIdsByPathId,
		filter.MimeType: findAssetIdsByMimeType,
		filter.FileName: findAssetIdsByFileName,
		filter.PathName: findAssetIdsByPathName,
	}

	for filterValue, filterFunction := range filterFunctions {
		foundIds, err := filterFunction.(func(any) (ScoredIdMap, error))(filterValue)
		if err != nil {
			return nil, err
		}
		if foundIds != nil {
			fmt.Printf("Found %d assets by: %v\n", len(foundIds), filterValue)
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
			items, err := loadList(query, params...)
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

		return loadList(query, params...)
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

// loadList queries the database
func loadList(query string, params ...any) ([]AssetListItem, error) {

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

func findAssetIdsByPathName(name any) (ScoredIdMap, error) {

	var sName = name.(string)
	if len(sName) == 0 {
		return nil, nil
	}

	//TODO search parents

	var query = "SELECT a.id, p.name FROM origin o " +
		"INNER JOIN pathItem p ON p.id = o.path " +
		"INNER JOIN asset a ON o.asset = a.id " +
		"WHERE p.name like ?;"

	var findName = sName
	if strings.Contains(findName, "*") {
		//If an asterisk was given explicitly, only use this, just replace to sql-asterisk
		findName = strings.ReplaceAll(sName, "*", "%")
	} else {
		findName = "%" + strings.ReplaceAll(sName, " ", "%") + "%"
	}
	fmt.Printf("findAssetIdsByPathName: %s\n", findName)

	return findAssetIds(query, findName, func(id int64, match any, idMap *ScoredIdMap) {
		score := float32(len(sName)) / float32(len(match.(string)))
		idMap.Add(id, score)
	})
}

func findAssetIdsByFileName(name any) (ScoredIdMap, error) {

	var sName = name.(string)
	if len(sName) == 0 {
		return nil, nil
	}

	var query = "SELECT a.id, f.name FROM origin o " +
		"INNER JOIN fileName f ON f.id = o.name " +
		"INNER JOIN asset a ON o.asset = a.id " +
		"WHERE f.name like ?;"

	var findName = sName
	if strings.Contains(findName, "*") {
		//If an asterisk was given explicitly, only use this, just replace to sql-asterisk
		findName = strings.ReplaceAll(sName, "*", "%")
	} else {
		findName = "%" + strings.ReplaceAll(sName, " ", "%") + "%"
	}
	fmt.Printf("findAssetIdsByFileName: %s\n", findName)

	return findAssetIds(query, findName, func(id int64, match any, idMap *ScoredIdMap) {
		score := float32(len(sName)) / float32(len(match.(string)))
		//fmt.Printf("Match: %s, Score: %f\n", match, score)
		idMap.Add(id, score)
	})
}

func findAssetIdsByMimeType(name any) (ScoredIdMap, error) {

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

	return findAssetIds(query, name, func(id int64, match any, idMap *ScoredIdMap) {
		dt := match.(time.Time)
		score := float32(dt.Unix()) / float32(1000.0)
		idMap.Set(id, score)
	})
}

func findAssetIdsByPathId(pathId any) (ScoredIdMap, error) {

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

func findAssetIds(query string, name any, calcScore func(id int64, match any, idMap *ScoredIdMap)) (ScoredIdMap, error) {

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer util.CloseOrLog(stmt)

	if rows, err := stmt.Query(name); err == nil {
		defer util.CloseOrLog(rows)

		ids := make(ScoredIdMap)

		for rows.Next() {
			var id int64
			var match any
			if err := rows.Scan(&id, &match); err != nil {
				return nil, err
			}
			calcScore(id, match, &ids)
		}
		return ids, nil
	} else {
		return nil, err
	}
}
