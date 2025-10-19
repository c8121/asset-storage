package metadata_db

import (
	"fmt"
	"strings"
	"time"

	"github.com/c8121/asset-storage/internal/util"
)

type AssetListItem struct {
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

	if filter.PathId != 0 {
		var err error = nil
		ids, err = findAssetIdsByPathId(filter.PathId)
		if err != nil {
			return nil, err
		}
		fmt.Printf("Found %d assets by path id: %d\n", len(ids), filter.PathId)
	}

	if filter.MimeType != "" {
		mimeTypeIds, err := findAssetIdsByMimeType(filter.MimeType)
		if err != nil {
			return nil, err
		}
		fmt.Printf("Found %d assets by mime-type: %s\n", len(mimeTypeIds), filter.MimeType)
		if ids != nil {
			ids.Reduce(mimeTypeIds)
		} else {
			ids = mimeTypeIds
		}
	}

	if filter.FileName != "" {
		fileNameIds, err := findAssetIdsByFileName(filter.FileName)
		if err != nil {
			return nil, err
		}
		fmt.Printf("Found %d assets by file name: %s\n", len(fileNameIds), filter.FileName)
		if ids != nil {
			ids.Reduce(fileNameIds)
		} else {
			ids = fileNameIds
		}
	}

	if filter.PathName != "" {
		pathNameIds, err := findAssetIdsByPathName(filter.PathName)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
		}
		fmt.Printf("Found %d assets by path name: %s\n", len(pathNameIds), filter.PathName)
		if ids != nil {
			ids.Reduce(pathNameIds)
		} else {
			ids = pathNameIds
		}
	}

	var query = "SELECT a.hash, m.name as mimeType, a.fileTime, " +
		" (SELECT name FROM origin where origin.asset = a.id LIMIT 1) as name " +
		" FROM asset a " +
		" INNER JOIN mimeType m ON a.mimeType = m.id "

	var params = make([]any, 0)

	if ids != nil {
		sorted := ids.Sort()
		endIdx := filter.Offset + filter.Count
		if endIdx >= len(sorted) {
			endIdx = len(sorted)
		}
		if endIdx > 0 {
			slice := sorted[filter.Offset:endIdx]
			fmt.Printf("Found %d (%d - %d) items\n", len(slice), filter.Offset, endIdx)
			query += " WHERE a.id in(" +
				strings.Repeat("?,", len(slice)-1) + "?" +
				");"
			for _, id := range slice {
				params = append(params, id.Id)
			}
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
	}

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
			if err := rows.Scan(&item.Hash, &item.MimeType, &item.FileTime, &item.Name); err != nil {
				return items, err
			}
			items = append(items, item)
		}

	} else {
		return items, err
	}

	return items, nil
}

func addWhere(sql string, add string) string {
	if len(sql) > 0 {
		return sql + " AND " + add
	} else {
		return add
	}
}

func findAssetIdsByPathName(name string) (ScoredIdMap, error) {

	//TODO search parents

	var query = "SELECT a.id, p.name FROM origin o " +
		"INNER JOIN pathItem p ON p.id = o.path " +
		"INNER JOIN asset a ON o.asset = a.id " +
		"WHERE p.name like ?;"

	var findName = "%" + strings.ReplaceAll(name, " ", "%") + "%"

	return findAssetIds(query, findName, func(match any) float32 {
		return float32(len(name)) / float32(len(match.(string)))
	})
}

func findAssetIdsByFileName(name string) (ScoredIdMap, error) {

	var query = "SELECT a.id, f.name FROM origin o " +
		"INNER JOIN fileName f ON f.id = o.name " +
		"INNER JOIN asset a ON o.asset = a.id " +
		"WHERE f.name like ?;"

	var findName = "%" + strings.ReplaceAll(name, " ", "%") + "%"

	return findAssetIds(query, findName, func(match any) float32 {
		return float32(len(name)) / float32(len(match.(string)))
	})
}

func findAssetIdsByMimeType(name string) (ScoredIdMap, error) {

	var query = "SELECT a.id, m.name FROM asset a " +
		"INNER JOIN mimeType m ON m.id = a.mimeType " +
		"WHERE m.name like ?;"

	return findAssetIds(query, name, func(match any) float32 {
		return 0.0
	})
}

func findAssetIdsByPathId(pathId int64) (ScoredIdMap, error) {

	var query = "SELECT a.id, a.id FROM origin o " +
		"INNER JOIN asset a ON o.asset = a.id " +
		"WHERE path = ?;"

	return findAssetIds(query, pathId, func(match any) float32 {
		return 0.0
	})
}

func findAssetIds(query string, name any, calcScore func(match any) float32) (ScoredIdMap, error) {

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
			score := calcScore(match)
			ids.Add(id, score)
		}
		return ids, nil
	} else {
		return nil, err
	}
}
