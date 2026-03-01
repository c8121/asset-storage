package metadata_db

import (
	"fmt"
	"strconv"
	"strings"

	metadata_db_entity "github.com/c8121/asset-storage/internal/metadata-db-entity"
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
		return nil, fmt.Errorf("invalid Face query")
	}
	hash := sFace[:p]
	assetId := metadata_db_entity.GetAssetId(hash)
	faceIdx, _ := strconv.Atoi(sFace[p+1:])

	var query = "SELECT a.id, 2.0 as score FROM asset a WHERE id = ? " +
		"UNION " +
		"SELECT asset_b, score FROM faceSimilarity " +
		"WHERE asset_a = ? AND face_a = ? " +
		"UNION " +
		"SELECT asset_a, score FROM faceSimilarity " +
		"WHERE asset_b = ? AND face_b = ?;"

	fmt.Printf("findAssetIdsByFace: %s\n", sFace)

	return findAssetIds(func(id int64, match any, idMap *ScoredIdMap) {
		score := float32(match.(float64))
		idMap.Add(id, score)
	}, query, hash, assetId, faceIdx, assetId, faceIdx)

}
