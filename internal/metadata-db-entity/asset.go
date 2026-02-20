package metadata_db_entity

import (
	"context"
	"database/sql"
	"errors"
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

// GetAssetIdTx gets Asset-ID from db
func GetAssetIdTx(tx *sql.Tx, hash string) int64 {
	var asset = &Asset{Hash: hash}
	err := LoadTx(tx, asset)
	if !errors.Is(err, ErrNotFound) && err != nil {
		return 0
	}
	return asset.Id
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
