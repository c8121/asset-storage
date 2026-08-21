package metadata_db

import (
	"fmt"
	"strings"
)

type FinderByPathName struct {
}

// Find searches all assets assigned to given path id
func (f FinderByPathName) Find(name any) (ScoredIdMap, error) {

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

	return findAssetIds(func(id int64, match any, idMap *ScoredIdMap) {
		score := float32(len(sName)) / float32(len(match.(string)))
		idMap.Add(id, score)
	}, query, findName)

}
