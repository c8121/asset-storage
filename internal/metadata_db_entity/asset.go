package metadata_db_entity

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/c8121/asset-storage/internal/metadata"
	"github.com/c8121/asset-storage/internal/metadata_db_conn"
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
func AddMetaData(tx *sql.Tx, jsonMeta *metadata.JsonAssetMetaData) error {

	var asset = &Asset{Hash: jsonMeta.Hash}
	err := Load(tx, asset)
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
		asset.Name = GetFileNameId(tx, latestOrigin.Name, true)
	}

	err = Save(tx, asset)
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
			Name:     GetFileNameId(tx, jsonOrigin.Name, true),
			Path:     GetPathItemIdTx(tx, jsonOrigin.Path, true),
			Owner:    GetOwnerIdTx(tx, jsonOrigin.Owner, true),
			FileTime: jsonOrigin.FileTime,
		}
		err = Save(tx, origin)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetMetaData loads metadata from database
func GetMetaData(hash string) (*Asset, error) {

	var asset = &Asset{Hash: hash}
	err := Load(metadata_db_conn.GetDatabase(), asset)
	if !errors.Is(err, ErrNotFound) && err != nil {
		fmt.Printf("Failed load asset: %s\n", err)
		return nil, err
	}
	return asset, nil
}

// GetMetaDataById loads metadata from database
func GetMetaDataById(assetId int64) (*Asset, error) {

	stmt, err := metadata_db_conn.GetDatabase().Prepare("SELECT id, hash, mimeType, fileTime, name FROM asset WHERE id = ?;")
	if err != nil {
		return nil, err
	}
	defer util.CloseOrLog(stmt)

	if rows, err := stmt.Query(assetId); err == nil {
		defer util.CloseOrLog(rows)
		if rows.Next() {
			asset := Asset{}
			if err := rows.Scan(&asset.Id, &asset.Hash, &asset.MimeType, &asset.FileTime, &asset.Name); err != nil {
				fmt.Printf("Error scanning rows: %s\n", err)
				return nil, err
			}
			return &asset, nil
		}

		return nil, ErrNotFound

	} else {
		return nil, err
	}
}

func RemoveMetaData(tx *sql.Tx, assetId int64, pathId int64) (int, error) {

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
	return Load(metadata_db_conn.GetDatabase(), a)
}

func (a *Asset) Save(tx *sql.Tx) error {
	return Save(tx, a)
}

func (a *Asset) Get(tx *sql.Tx, insertIfNotExists bool) error {
	return Get(tx, insertIfNotExists, a)
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
