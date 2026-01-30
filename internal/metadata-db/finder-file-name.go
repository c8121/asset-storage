package metadata_db

import (
	"fmt"
	"strings"
)

type FinderByFileName struct {
}

// Find searches all assets having the given name as origin.FileName
func (f FinderByFileName) Find(name any) (ScoredIdMap, error) {

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
