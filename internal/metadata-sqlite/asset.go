package mdsqlite

import (
	"time"

	"github.com/c8121/asset-storage/internal/util"
)

type Asset struct {
	Hash     string
	Name     string
	FileTime time.Time
}

type AssetFilter struct {
	MimeType string
	Offset   int
	Count    int
}

// ListAssets loads a list of Asset's from database, matching given filter
func ListAssets(filter *AssetFilter) ([]Asset, error) {

	var query = "SELECT o.hash, name, filetime FROM origin o INNER JOIN asset a ON o.hash = a.hash "
	var where = ""
	var limit = "ORDER BY filetime DESC, o.hash ASC LIMIT ? OFFSET ?;"

	var params = make([]any, 0)

	if filter.MimeType != "" {
		params = append(params, filter.MimeType)
		where += "(mimetype LIKE ?)"
	}

	params = append(params, filter.Count)
	params = append(params, filter.Offset)

	if where != "" {
		query += "WHERE " + where + " "
	}
	query += limit
	//fmt.Printf("Query: %s %v\n", query, params)

	stmt, err := DB.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer util.CloseOrLog(stmt)

	var items []Asset

	if rows, err := stmt.Query(params...); err == nil {
		defer util.CloseOrLog(rows)
		for rows.Next() {
			var item Asset
			if err := rows.Scan(&item.Hash, &item.Name, &item.FileTime); err != nil {
				return items, err
			}
			items = append(items, item)
		}

	} else {
		return items, err
	}

	return items, nil
}
