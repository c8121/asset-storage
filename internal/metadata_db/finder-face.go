package metadata_db

import (
	"github.com/c8121/asset-storage/internal/metadata_db_entity"
)

type FinderByFace struct {
}

// Find searches all assets having the given face
func (f FinderByFace) Find(faceId any) (ScoredIdMap, error) {

	if faceId.(int64) == 0 {
		return nil, nil
	}

	similarFaces, err := metadata_db_entity.FindSimilarFaces(faceId.(int64), 5)
	if err != nil {
		return nil, err
	}

	if similarFaces == nil {
		return nil, nil
	}

	ids := ScoredIdMap{}
	for _, face := range *similarFaces {
		ids.Add(face.AssetId, 0)
	}

	return ids, nil
}
