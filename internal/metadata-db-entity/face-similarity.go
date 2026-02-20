package metadata_db_entity

import (
	"context"
	"database/sql"

	"github.com/c8121/asset-storage/internal/util"
)

type FaceSimilarity struct {
	Id     int64
	AssetA int64
	FaceA  int
	AssetB int64
	FaceB  int
}

// AddFaceSimilarity adds/updates face-similarity in database
func AddFaceSimilarity(hashA string, faceA int, hashB string, faceB int) error {

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer util.RollbackOrLog(tx)

	err = AddFaceSimilarityTx(tx, hashA, faceA, hashB, faceB)
	if err != nil {
		return err
	}

	return util.CommitOrLog(tx)
}

// AddFaceSimilarityTx adds/updates face-similarity in database
func AddFaceSimilarityTx(tx *sql.Tx, hashA string, faceA int, hashB string, faceB int) error {

	var similarity = &FaceSimilarity{
		AssetA: GetAssetIdTx(tx, hashA),
		FaceA:  faceA,
		AssetB: GetAssetIdTx(tx, hashB),
		FaceB:  faceB,
	}

	err := InsertTx(tx, similarity)
	if err != nil {
		return err
	}

	return nil
}

func (f *FaceSimilarity) GetId() int64 {
	return f.Id
}

func (f *FaceSimilarity) Save() error {
	return Save(f)
}

func (f *FaceSimilarity) GetSelectQuery() string {
	return "SELECT id, asset_a, face_a, asset_b, face_b FROM faceSimilarity WHERE asset_a = ? and face_a = ? and asset_b = ? and face_b = ?;"
}

func (f *FaceSimilarity) GetSelectQueryArgs() []any {
	return []any{f.AssetA, f.FaceA, f.AssetB, f.FaceB}
}

func (f *FaceSimilarity) Scan(rows *sql.Rows) error {
	return rows.Scan(&f.Id, &f.AssetA, &f.FaceA, &f.AssetB, &f.FaceB)
}

func (f *FaceSimilarity) GetInsertQuery() string {
	return "INSERT INTO faceSimilarity(asset_a, face_a, asset_b, face_b) VALUES(?,?,?,?);"
}

func (f *FaceSimilarity) Exec(stmt *sql.Stmt) (sql.Result, error) {
	return stmt.Exec(&f.AssetA, &f.FaceA, &f.AssetB, &f.FaceB, &f.Id)
}

func (f *FaceSimilarity) SetId(id int64) {
	f.Id = id
}

func (a *FaceSimilarity) GetCreateQueries() []string {
	return []string{
		"CREATE TABLE IF NOT EXISTS faceSimilarity(id integer PRIMARY KEY, asset_a integer, face_a integer, asset_b integer, face_b integer);",
		"CREATE INDEX IF NOT EXISTS idx_faceSimilarity_asset_a on faceSimilarity(asset_a);",
		"CREATE INDEX IF NOT EXISTS idx_faceSimilarity_face_a on faceSimilarity(face_a);",
		"CREATE INDEX IF NOT EXISTS idx_faceSimilarity_asset_b on faceSimilarity(asset_b);",
		"CREATE INDEX IF NOT EXISTS idx_faceSimilarity_face_b on faceSimilarity(face_b);",
	}
}
