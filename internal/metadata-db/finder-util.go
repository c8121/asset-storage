package metadata_db

import "github.com/c8121/asset-storage/internal/util"

func findAssetIds(calcScore func(id int64, match any, idMap *ScoredIdMap), query string, args ...any) (ScoredIdMap, error) {

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer util.CloseOrLog(stmt)

	if rows, err := stmt.Query(args...); err == nil {
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
