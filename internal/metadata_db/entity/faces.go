package metadata_db_entity

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"

	"github.com/c8121/asset-storage/internal/faces"
	"github.com/c8121/asset-storage/internal/util"
	"github.com/c8121/go_db_entities"

	_ "gosqlite.org/vec" //Auto-registers the vec0 module on every connection
)

type FaceEmbeddingQueryBuilder interface {
	GetFindSimilarQuery(ctx *FacesContext, embedding *[]float32, max int) string
}

type FacesContext struct {
	DB                        *sql.DB
	FaceEmbeddingQueryBuilder FaceEmbeddingQueryBuilder
}

func NewContext(database *sql.DB, queryBuilder FaceEmbeddingQueryBuilder) *FacesContext {
	return &FacesContext{
		DB:                        database,
		FaceEmbeddingQueryBuilder: queryBuilder,
	}
}

var (
	FacesDbContext *FacesContext
)

type FaceEmbedding struct {
	Id             int64
	AssetId        int64 `dbe:"Index"`
	X1, Y1, X2, Y2 int
	Embedding      []float32 `dbe:"Type: Embedding, Size: 512"`
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

	face := &FaceEmbedding{}
	err := go_db_entities.DeleteBatch(DbContext, face, go_db_entities.Filter{FieldName: "AssetId", Comparison: go_db_entities.CompareEqual, Value: assetId})
	if err != nil {
		return err
	}

	return nil
}

// GetFaces finds all faces that where detected in the given asset by face-api
// Loads FaceEmbedding-Struct without FaceEmbedding.Embedding (use GetFaceEmbedding to load that)
func GetFaces(hash string) (*[]FaceEmbedding, error) {

	assetMeta, err := LoadAsset(hash)
	if err != nil {
		return nil, err
	}

	query := go_db_entities.TableQuery{
		FieldNames: []string{"Id", "AssetId", "X1", "Y1", "X2", "Y2"},
		Entity:     &FaceEmbedding{},
		Filter:     go_db_entities.NewFilter("AssetId", go_db_entities.CompareEqual, assetMeta.Id),
	}

	rows, err := go_db_entities.Query(DbContext, query)
	if err != nil && !errors.Is(go_db_entities.ErrNotFound, err) {
		return nil, err
	}

	faces := make([]FaceEmbedding, len(rows))
	for i, row := range rows {
		face := FaceEmbedding{
			Id:      row[0].(int64),
			AssetId: row[1].(int64),
			X1:      int(row[2].(int64)),
			Y1:      int(row[3].(int64)),
			X2:      int(row[4].(int64)),
			Y2:      int(row[5].(int64)),
		}
		faces[i] = face
	}

	return &faces, nil
}

// GetFaceEmbedding loads the stored embedding
func GetFaceEmbedding(faceId int64) (*FaceEmbedding, error) {

	faceEmbedding := &FaceEmbedding{}
	err := go_db_entities.Load(DbContext, faceEmbedding, *go_db_entities.NewIdFilter(faceId))
	if err != nil {
		return nil, err
	}

	return faceEmbedding, nil
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

	embeddingJson, err := json.Marshal(embedding)
	if err != nil {
		return nil, err
	}

	db := FacesDbContext.DB
	sql := FacesDbContext.FaceEmbeddingQueryBuilder.GetFindSimilarQuery(FacesDbContext, embedding, max)

	stmt, err := db.Prepare(sql)
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
	return list, nil
}
