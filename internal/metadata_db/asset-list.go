package metadata_db

import (
	"fmt"
	"strings"
	"time"

	"github.com/c8121/asset-storage/internal/metadata_db_conn"
	"github.com/c8121/asset-storage/internal/util"
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

	return ListAssetsByIds(ids, filter.Offset, filter.Count)
}

// ListAssetsByIds creates a list of AssetListItem's, loaded by given ScoredIdMap
func ListAssetsByIds(ids ScoredIdMap, offset int, count int) ([]AssetListItem, error) {
	var query = "SELECT a.id, a.hash, m.name as mimeType, a.fileTime, f.name" +
		" FROM asset a " +
		" INNER JOIN mimeType m ON a.mimeType = m.id " +
		" INNER JOIN fileName f ON a.name = f.id "

	var params = make([]any, 0)

	if ids != nil {
		sorted := ids.Sort()
		endIdx := offset + count
		if endIdx >= len(sorted) {
			endIdx = len(sorted)
		}
		if offset < len(sorted) && endIdx > 0 {
			slice := sorted[offset:endIdx]
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

			if err = assignPathIds(&list); err != nil {
				return nil, err
			}
			for _, item := range list {
				fmt.Printf("%d items\n", len(item.PathIds))
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
		params = append(params, count)
		params = append(params, offset)

		list, err := loadAssetList(query, params...)
		if err != nil {
			return nil, err
		}
		if err = assignPathIds(&list); err != nil {
			return nil, err
		}
		for _, item := range list {
			fmt.Printf("%d items\n", len(item.PathIds))
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
func loadAssetList(query string, params ...any) ([]AssetListItem, error) {

	//fmt.Printf("Query: %s\n", query)

	stmt, err := metadata_db_conn.GetDatabase().Prepare(query)
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

// assignPathIds uses FindPathIds and assigns to AssetListItem.PathIds
func assignPathIds(items *[]AssetListItem) error {

	idMap, err := FindPathIds(*items)
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

// FindPathIds looks for origin-paths in database and
// returns it as map: asset.Id -> [path.Id,...]
func FindPathIds(items []AssetListItem) (map[int64][]int64, error) {

	var idMap = make(map[int64][]int64, len(items))
	var chunk = make([]AssetListItem, 0)

	for i, item := range items {

		chunk = append(chunk, item)
		if len(chunk) == 4 || (i+1) == len(items) {

			var query = "SELECT a.id, o.path" +
				" FROM origin o " +
				" INNER JOIN asset a ON o.asset = a.id "

			query += " WHERE a.id in(" +
				strings.Repeat("?,", len(chunk)-1) + "?" +
				");"

			var params = make([]any, 0)
			for _, id := range chunk {
				params = append(params, id.Id)
			}

			stmt, err := metadata_db_conn.GetDatabase().Prepare(query)
			if err != nil {
				return nil, err
			}
			defer util.CloseOrLog(stmt)

			if rows, err := stmt.Query(params...); err == nil {
				defer util.CloseOrLog(rows)
				for rows.Next() {
					var assetId int64
					var pathId int64
					if err := rows.Scan(&assetId, &pathId); err != nil {
						return nil, err
					}

					list, exists := idMap[assetId]
					if exists {
						idMap[assetId] = append(list, pathId)
					} else {
						list = make([]int64, 0)
						idMap[assetId] = append(list, pathId)
					}
				}

			} else {
				return nil, err
			}

			chunk = make([]AssetListItem, 0)
		}
	}

	return idMap, nil
}
