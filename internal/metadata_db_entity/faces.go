package metadata_db_entity

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/c8121/asset-storage/internal/faces"
	"github.com/c8121/asset-storage/internal/util"

	_ "gosqlite.org/vec" //Auto-registers the vec0 module on every connection
)

type FaceEmbedding struct {
	Id             int64
	AssetId        int64
	X1, Y1, X2, Y2 int
	Embedding      []float32
}

// AddFace
func AddFace(assetId int64, face *faces.RestApiFace) (*FaceEmbedding, error) {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer util.RollbackOrLog(tx)

	faceEmbedding, err := AddFaceTx(tx, assetId, face)
	if err != nil {
		return nil, err
	}

	return faceEmbedding, util.CommitOrLog(tx)
}

// AddFaces
func AddFaces(assetId int64, faces *[]faces.RestApiFace) (*[]FaceEmbedding, error) {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer util.RollbackOrLog(tx)

	list := &[]FaceEmbedding{}

	for _, face := range *faces {
		faceEmbedding, err := AddFaceTx(tx, assetId, &face)
		if err != nil {
			util.LogError(err)
		} else {
			*list = append(*list, *faceEmbedding)
		}
	}

	return list, util.CommitOrLog(tx)
}

// AddFaceTx
func AddFaceTx(tx *sql.Tx, assetId int64, face *faces.RestApiFace) (*FaceEmbedding, error) {

	var faceEmbedding = &FaceEmbedding{
		AssetId:   assetId,
		X1:        face.Image[0],
		Y1:        face.Image[1],
		X2:        face.Image[2],
		Y2:        face.Image[3],
		Embedding: face.Embedding,
	}

	err := InsertTx(tx, faceEmbedding)
	if err != nil {
		return nil, err
	}

	return faceEmbedding, nil
}

func RemoveFaces(assetId int64) error {

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer util.RollbackOrLog(tx)

	stmt, err := tx.Prepare("DELETE FROM faces WHERE assetId=?;")
	if err != nil {
		return err
	}
	defer util.CloseOrLog(stmt)

	_, err = stmt.Exec(assetId)
	if err != nil {
		return err
	}

	return util.CommitOrLog(tx)
}

// GetFaces finds all faces that where detected in the given asset by face-api
// Loads FaceEmbedding-Struct without FaceEmbedding.Embedding (use GetFaceEmbedding to load that)
func GetFaces(hash string) (*[]FaceEmbedding, error) {

	assetId := GetAssetId(hash)

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer util.RollbackOrLog(tx)

	query := "SELECT id, assetId, x1, y1, x2, y2 FROM faces WHERE assetId=?"
	rows, err := db.Query(query, assetId)
	if err != nil {
		log.Fatal("Query failed:", err)
	}
	defer util.CloseOrLog(rows)

	list := &[]FaceEmbedding{}

	for rows.Next() {
		faceEmbedding := FaceEmbedding{}
		if err := rows.Scan(&faceEmbedding.Id, &faceEmbedding.AssetId, &faceEmbedding.X1, &faceEmbedding.Y1, &faceEmbedding.X2, &faceEmbedding.Y2); err != nil {
			util.LogError(err)
		} else {
			*list = append(*list, faceEmbedding)
		}
	}

	return list, nil
}

// GetFaceEmbedding loads the stored embedding
func GetFaceEmbedding(faceId int64) (*FaceEmbedding, error) {

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer util.RollbackOrLog(tx)

	query := "SELECT id, assetId, vec_to_json(embedding) FROM faces WHERE id=?"
	rows, err := db.Query(query, faceId)
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
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer util.RollbackOrLog(tx)

	embeddingJson, err := json.Marshal(embedding)
	if err != nil {
		return nil, err
	}
	query := "WITH knn_matches AS (SELECT id, assetId, distance FROM faces WHERE embedding MATCH vec_f32(?) AND k = 99)" +
		"SELECT id, assetId, distance FROM knn_matches WHERE distance <= 20.0 ORDER BY distance;"
	rows, err := db.Query(query, string(embeddingJson), max)
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
			fmt.Printf("ID: %-5d | Asset ID: %-5d | Cosine Distance: %.4f\n", faceEmbedding.Id, faceEmbedding.AssetId, distance)
			*list = append(*list, faceEmbedding)
		}

	}
	return list, nil
}

func (f *FaceEmbedding) GetId() int64 {
	return f.Id
}

func (f *FaceEmbedding) Save() error {
	return Save(f)
}

func (f *FaceEmbedding) GetSelectQuery() string {
	return "SELECT id, assetId FROM faces WHERE id = ?;"
}

func (f *FaceEmbedding) GetSelectQueryArgs() []any {
	return []any{f.Id}
}

func (f *FaceEmbedding) Scan(rows *sql.Rows) error {
	return rows.Scan(&f.Id, &f.AssetId)
}

func (f *FaceEmbedding) GetInsertQuery() string {
	return "INSERT INTO faces(assetId, x1, y1, x2, y2, embedding) VALUES(?,?,?,?,?,vec_f32(?));"
}

func (f *FaceEmbedding) Exec(stmt *sql.Stmt) (sql.Result, error) {
	embeddingJson, err := json.Marshal(f.Embedding)
	if err != nil {
		return nil, err
	}
	return stmt.Exec(&f.AssetId, &f.X1, &f.Y1, &f.X2, &f.Y2, string(embeddingJson))
}

func (f *FaceEmbedding) SetId(id int64) {
	f.Id = id
}

func (a *FaceEmbedding) GetCreateQueries() []string {
	return []string{
		//"DROP TABLE IF EXISTS faces",
		"CREATE VIRTUAL TABLE IF NOT EXISTS faces USING vec0(id integer PRIMARY KEY, assetId integer, " +
			"x1 integer, y1 integer, x2 integer, y2 integer, " +
			"embedding FLOAT[512]);",
	}
}
