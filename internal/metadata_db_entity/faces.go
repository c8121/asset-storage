package metadata_db_entity

import (
	"errors"

	"github.com/c8121/asset-storage/internal/faces"
	"github.com/c8121/asset-storage/internal/util"
	"github.com/c8121/go_db_entities"

	_ "gosqlite.org/vec" //Auto-registers the vec0 module on every connection
)

type FaceEmbedding struct {
	Id             int64
	AssetId        int64
	X1, Y1, X2, Y2 int
	Embedding      []float32
}

func (f *FaceEmbedding) GetEntityName() string {
	return "Faces"
}

func (f *FaceEmbedding) GetId() int64 {
	return f.Id
}

func (f *FaceEmbedding) SetId(id int64) {
	f.Id = id
}

// AddFaces
func AddFaces(assetId int64, faces *[]faces.RestApiFace) (*[]FaceEmbedding, error) {

	list := &[]FaceEmbedding{}

	for _, face := range *faces {
		faceEmbedding, err := AddFace(assetId, &face)
		if err != nil {
			util.LogError(err)
		} else {
			*list = append(*list, *faceEmbedding)
		}
	}

	return list, nil
}

// AddFace
func AddFace(assetId int64, face *faces.RestApiFace) (*FaceEmbedding, error) {

	var faceEmbedding = &FaceEmbedding{
		AssetId:   assetId,
		X1:        face.Image[0],
		Y1:        face.Image[1],
		X2:        face.Image[2],
		Y2:        face.Image[3],
		Embedding: face.Embedding,
	}

	err := go_db_entities.Save(DbContext, faceEmbedding)
	if err != nil {
		return nil, err
	}

	return faceEmbedding, nil
}

func RemoveFaces(assetId int64) error {

	/*stmt, err := tx.Prepare("DELETE FROM faces WHERE assetId=?;")
	if err != nil {
		return err
	}
	defer util.CloseOrLog(stmt)

	_, err = stmt.Exec(assetId)
	return err*/
	panic("Not implemented")
}

// GetFaces finds all faces that where detected in the given asset by face-api
// Loads FaceEmbedding-Struct without FaceEmbedding.Embedding (use GetFaceEmbedding to load that)
func GetFaces(hash string) (*[]FaceEmbedding, error) {

	assetMeta, err := LoadAsset(hash)
	if err != nil {
		return nil, err
	}

	itemSupplier := func() go_db_entities.Entity {
		return &FaceEmbedding{}
	}

	items, err := go_db_entities.Select(DbContext, itemSupplier, *go_db_entities.NewFilter("AssetId", go_db_entities.CompareEqual, assetMeta.Id))
	if err != nil && !errors.Is(go_db_entities.ErrNotFound, err) {
		return nil, err
	}

	faces := make([]FaceEmbedding, len(items))
	for i, item := range items {
		face := item.(*FaceEmbedding)
		faces[i] = *face
	}

	return &faces, nil
}

// GetFaceEmbedding loads the stored embedding
func GetFaceEmbedding(faceId int64) (*FaceEmbedding, error) {

	/*stmt, err := metadata_db_conn.GetDatabase().Prepare("SELECT id, assetId, vec_to_json(embedding) FROM faces WHERE id=?")
		if err != nil {
			return nil, err
		}
		defer util.CloseOrLog(stmt)

		rows, err := stmt.Query(faceId)
		if err != nil {
			log.Fatal("Query failed:", err)
		}
		defer util.CloseOrLog(rows)

		faceEmbedding := &FaceEmbedding{}
		if rows.Next() {
			var embeddingJsonString string
			if err := rows.Scan(&faceEmbedding.Id, &faceEmbedding.AssetId, &embeddingJsonString); err != nil {
				util.LogError(err)
			}
			err = json.Unmarshal([]byte(embeddingJsonString), &faceEmbedding.Embedding)
			if err != nil {
				log.Fatal("Failed to parse vector JSON string: ", err)
			}
			return faceEmbedding, nil
		}

		return nil, nil
	}*/
	panic("Not implemented")
}

func FindSimilarFaces(faceId int64, max int) (*[]FaceEmbedding, error) {

	faceEmbedding, err := GetFaceEmbedding(faceId)
	if err != nil {
		return nil, err
	}

	if faceEmbedding == nil {
		return nil, nil
	}

	return FindSimilarFacesByEmbedding(&faceEmbedding.Embedding, max)
}

func FindSimilarFacesByEmbedding(embedding *[]float32, max int) (*[]FaceEmbedding, error) {

	/*embeddingJson, err := json.Marshal(embedding)
	if err != nil {
		return nil, err
	}

	stmt, err := metadata_db_conn.GetDatabase().Prepare("WITH knn_matches AS (SELECT id, assetId, distance FROM faces WHERE embedding MATCH vec_f32(?) AND k = 99)" +
		"SELECT id, assetId, distance FROM knn_matches WHERE distance <= 20.0 ORDER BY distance;")
	if err != nil {
		return nil, err
	}
	defer util.CloseOrLog(stmt)

	rows, err := stmt.Query(string(embeddingJson), max)
	if err != nil {
		log.Fatal("Query failed:", err)
	}
	defer util.CloseOrLog(rows)

	list := &[]FaceEmbedding{}

	for rows.Next() {
		faceEmbedding := FaceEmbedding{}
		var distance float32
		if err := rows.Scan(&faceEmbedding.Id, &faceEmbedding.AssetId, &distance); err != nil {
			util.LogError(err)
		} else {
			//fmt.Printf("ID: %-5d | Asset ID: %-5d | Cosine Distance: %.4f\n", faceEmbedding.Id, faceEmbedding.AssetId, distance)
			*list = append(*list, faceEmbedding)
		}

	}
	return list, nil*/
	panic("Not implemented")
}
