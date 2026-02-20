package metadata_db

import (
	"fmt"
	"strconv"
	"strings"
)

type FinderByFace struct {
}

// Find searches all assets having the given face
func (f FinderByFace) Find(face any) (ScoredIdMap, error) {

	var sFace = face.(string)
	if len(sFace) == 0 {
		return nil, nil
	}

	p := strings.Index(sFace, "/")
	if p == -1 {
		return nil, fmt.Errorf("Invalid Face query")
	}
	hash := sFace[:p]
	faceIdx, _ := strconv.Atoi(sFace[p+1:])

	var query = "SELECT a.id, 2 as score FROM asset a WHERE a.hash = ? " +
		"UNION " +
		"SELECT asset_b, f.id as score FROM faceSimilarity f " +
		"INNER JOIN asset a ON f.asset_a = a.id " +
		"WHERE a.hash = ? AND f.face_a = ? " +
		"UNION " +
		"SELECT asset_a, f.id as score FROM faceSimilarity f " +
		"INNER JOIN asset a ON f.asset_b = a.id " +
		"WHERE a.hash = ? AND f.face_b = ?;"

	fmt.Printf("findAssetIdsByFace: %s\n", sFace)

	return findAssetIds(func(id int64, match any, idMap *ScoredIdMap) {
		idMap.Add(id, 1)
	}, query, hash, hash, faceIdx, hash, faceIdx)

}
