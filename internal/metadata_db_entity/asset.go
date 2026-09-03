package metadata_db_entity

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/c8121/asset-storage/internal/metadata"
	"github.com/c8121/asset-storage/internal/util"
)

type Asset struct {
	Id       int64
	Hash     string
	MimeType int64
	FileTime time.Time //Max of all origins
	Name     int64     //Latest name
}

// AddMetaData adds/updates meta-data in database
func AddMetaData(jsonMeta *metadata.JsonAssetMetaData) error {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer util.RollbackOrLog(tx)

	err = AddMetaDataTx(tx, jsonMeta)
	if err != nil {
		return err
	}

	return util.CommitOrLog(tx)
}

// AddMetaDataTx adds/updates meta-data in database
func AddMetaDataTx(tx *sql.Tx, jsonMeta *metadata.JsonAssetMetaData) error {

	var asset = &Asset{Hash: jsonMeta.Hash}
	err := LoadTx(tx, asset)
	if !errors.Is(err, ErrNotFound) && err != nil {
		return err
	}

	mimeType, err := GetMimeTypeTx(tx, jsonMeta.MimeType, true)
	if err != nil {
		return err
	}

	asset.MimeType = mimeType.Id

	latestOrigin := metadata.GetLatestOrigin(jsonMeta)
	if latestOrigin != nil {
		asset.FileTime = latestOrigin.FileTime
		asset.Name = GetFileNameIdTx(tx, latestOrigin.Name, true)
	}

	err = SaveTx(tx, asset)
	if err != nil {
		return err
	}

	err = RemoveOriginsTx(tx, asset)
	if err != nil {
		return err
	}

	for _, jsonOrigin := range jsonMeta.Origins {

		var origin = &Origin{
			Asset:    asset.Id,
			Name:     GetFileNameIdTx(tx, jsonOrigin.Name, true),
			Path:     GetPathItemIdTx(tx, jsonOrigin.Path, true),
			Owner:    GetOwnerIdTx(tx, jsonOrigin.Owner, true),
			FileTime: jsonOrigin.FileTime,
		}
		err = SaveTx(tx, origin)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetAssetId gets Asset-ID from db
func GetAssetId(hash string) int64 {

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Printf("Failed to begin transaction: %s\n", err)
		return 0
	}
	defer util.RollbackOrLog(tx)

	var asset = &Asset{Hash: hash}
	err = LoadTx(tx, asset)
	if !errors.Is(err, ErrNotFound) && err != nil {
		fmt.Printf("Failed load asset id: %s\n", err)
		return 0
	}
	return asset.Id
}

// GetAssetIdTx gets Asset-ID from db
func GetAssetIdTx(tx *sql.Tx, hash string) int64 {
	var asset = &Asset{Hash: hash}
	err := LoadTx(tx, asset)
	if !errors.Is(err, ErrNotFound) && err != nil {
		return 0
	}
	return asset.Id
}

// GetAssetHash gets hash from db
func GetAssetHash(assetId int64) string {

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Printf("Failed to begin transaction: %s\n", err)
		return ""
	}
	defer util.RollbackOrLog(tx)

	stmt, err := tx.Prepare("SELECT hash FROM asset WHERE id = ?;")
	if err != nil {
		fmt.Printf("Failed prepare statement: %s\n", err)
		return ""
	}
	defer util.CloseOrLog(stmt)

	if rows, err := stmt.Query(assetId); err == nil {
		defer util.CloseOrLog(rows)
		if rows.Next() {
			var hash string
			if err := rows.Scan(&hash); err != nil {
				fmt.Printf("Error scanning rows: %s\n", err)
				return ""
			}
			return hash
		}
	}

	return ""
}

func RemoveMetaData(assetId int64, pathId int64) (int, error) {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 9999, err
	}
	defer util.RollbackOrLog(tx)

	remainingOrigins, err := RemoveMetaDataTx(tx, assetId, pathId)
	if err != nil {
		return 9999, err
	}

	return remainingOrigins, util.CommitOrLog(tx)
}

func RemoveMetaDataTx(tx *sql.Tx, assetId int64, pathId int64) (int, error) {

	if pathId > 0 {
		if err := RemoveOriginsByAssetIdAndPathIdTx(tx, assetId, pathId); err != nil {
			return 9999, err
		}
	} else {
		if err := RemoveOriginsByAssetIdTx(tx, assetId); err != nil {
			return 9999, err
		}
	}

	remainingOrigins, err := GetOriginsTx(tx, assetId)
	if err != nil {
		return 9999, err
	}
	if len(*remainingOrigins) > 0 {
		fmt.Printf("Keep asset, it has more origins: %d\n", assetId)
		return len(*remainingOrigins), nil
	}

	stmt, err := tx.Prepare("DELETE FROM asset WHERE id=?;")
	if err != nil {
		return 9999, err
	}
	defer util.CloseOrLog(stmt)

	_, err = stmt.Exec(assetId)
	if err != nil {
		return 9999, err
	}

	return 0, nil
}

func (a *Asset) GetId() int64 {
	return a.Id
}

func (a *Asset) Load() error {
	return Load(a)
}

func (a *Asset) Save() error {
	return Save(a)
}

func (a *Asset) Get(insertIfNotExists bool) error {
	return Get(insertIfNotExists, a)
}

func (a *Asset) GetSelectQuery() string {
	return "SELECT id, hash, mimeType, fileTime, name FROM asset WHERE hash = ?;"
}

func (a *Asset) GetSelectQueryArgs() []any {
	return []any{a.Hash}
}

func (a *Asset) Scan(rows *sql.Rows) error {
	return rows.Scan(&a.Id, &a.Hash, &a.MimeType, &a.FileTime, &a.Name)
}

func (a *Asset) GetInsertQuery() string {
	return "INSERT INTO asset(hash, mimeType, fileTime, name) VALUES(?,?,?,?);"
}

func (a *Asset) GetUpdateQuery() string {
	return "UPDATE asset SET hash=?, mimeType=?, fileTime=?, name=? WHERE id = ?;"
}

func (a *Asset) GetUpdateQueryArgs() []any {
	return []any{&a.Hash, &a.MimeType, &a.FileTime, &a.Name, &a.Id}
}

func (a *Asset) Exec(stmt *sql.Stmt) (sql.Result, error) {
	return stmt.Exec(&a.Hash, &a.MimeType, &a.FileTime, &a.Name, &a.Id)
}

func (a *Asset) SetId(id int64) {
	a.Id = id
}

func (a *Asset) GetCreateQueries() []string {
	return []string{
		"CREATE TABLE IF NOT EXISTS asset(id integer PRIMARY KEY, hash TEXT(64), mimeType integer, fileTime DATETIME, name integer);",
		"CREATE INDEX IF NOT EXISTS idx_asset_hash on asset(hash);",
		"CREATE INDEX IF NOT EXISTS idx_asset_mimeType on asset(mimeType);",
		"CREATE INDEX IF NOT EXISTS idx_asset_fileTime on asset(fileTime);",
		"CREATE INDEX IF NOT EXISTS idx_asset_name on asset(name);",
	}
}
