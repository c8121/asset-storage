package metadata_db

import (
	"github.com/c8121/asset-storage/internal/util"
)

func ListPathItems(parentId int) ([]PathItem, error) {

	var query = "SELECT id, parent, name " +
		" FROM pathItem WHERE parent = ? ORDER BY name, id asc LIMIT 9999;"

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer util.CloseOrLog(stmt)

	var items []PathItem

	if rows, err := stmt.Query(parentId); err == nil {
		defer util.CloseOrLog(rows)
		for rows.Next() {
			var item PathItem
			if err := rows.Scan(&item.Id, &item.Parent, &item.Name); err != nil {
				return items, err
			}
			items = append(items, item)
		}

	} else {
		return items, err
	}

	return items, nil
}
