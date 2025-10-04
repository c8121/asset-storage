package metadata_db

import (
	"fmt"
	"time"

	"github.com/c8121/asset-storage/internal/util"
)

type AssetListItem struct {
	Hash     string
	MimeType string
	FileTime time.Time
}

type AssetListFilter struct {
	MimeType string
	Offset   int
	Count    int
}

func ListAssets(filter *AssetListFilter) ([]AssetListItem, error) {

	var query = "SELECT a.hash, m.name as mimeType, a.fileTime FROM asset a " +
		" INNER JOIN mimeType m ON a.mimeType = m.id "
	var where = ""
	var limit = "ORDER BY fileTime DESC, hash ASC LIMIT ? OFFSET ?;"

	var params = make([]any, 0)

	if filter.MimeType != "" {
		params = append(params, filter.MimeType)
		where += "(m.name LIKE ?)"
	}

	params = append(params, filter.Count)
	params = append(params, filter.Offset)

	if where != "" {
		query += "WHERE " + where + " "
	}
	query += limit
	fmt.Printf("Query: %s %v\n", query, params)

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
			if err := rows.Scan(&item.Hash, &item.MimeType, &item.FileTime); err != nil {
				return items, err
			}
			items = append(items, item)
		}

	} else {
		return items, err
	}

	return items, nil
}
