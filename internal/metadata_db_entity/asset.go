package metadata_db_entity

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/c8121/asset-storage/internal/db_entity"
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

func LoadAsset(prv db_entity.StatementProvider, hash string) (*Asset, error) {
	var asset = &Asset{Hash: hash}
	err := db_entity.Load(prv, asset)
	if !errors.Is(err, db_entity.ErrNotFound) && err != nil {
		return nil, err
	}
	return asset, err
}

func LoadAssetById(prv db_entity.StatementProvider, id int64) (*Asset, error) {
	var asset = &Asset{Id: id}
	err := db_entity.Load(prv, asset)
	if !errors.Is(err, db_entity.ErrNotFound) && err != nil {
		return nil, err
	}
	return asset, err
}

// AddAsset adds/updates meta-data in database
func AddAsset(prv db_entity.StatementProvider, jsonMeta *metadata.JsonAssetMetaData) error {

	asset, err := LoadAsset(prv, jsonMeta.Hash)
	if err != nil && !errors.Is(err, db_entity.ErrNotFound) {
		return err
	}

	mimeType, err := LoadMimeType(prv, jsonMeta.MimeType, true)
	if err != nil {
		return err
	}

	asset.MimeType = mimeType.Id

	latestOrigin := metadata.GetLatestOrigin(jsonMeta)
	if latestOrigin != nil {
		asset.FileTime = latestOrigin.FileTime
		asset.Name = getFileNameId(prv, latestOrigin.Name)
	}

	err = db_entity.SaveEntity(prv, asset)
	if err != nil {
		return err
	}

	err = RemoveOrigins(prv, asset)
	if err != nil {
		return err
	}

	for _, jsonOrigin := range jsonMeta.Origins {

		var origin = &Origin{
			Asset:    asset.Id,
			Name:     getFileNameId(prv, jsonOrigin.Name),
			Path:     getPathItemId(prv, jsonOrigin.Path),
			Owner:    getOwnerId(prv, jsonOrigin.Owner),
			FileTime: jsonOrigin.FileTime,
		}
		err = db_entity.SaveEntity(prv, origin)
		if err != nil {
			return err
		}
	}

	return nil
}

func getFileNameId(prv db_entity.StatementProvider, name string) int64 {
	fileName, err := LoadFileName(prv, name, true)
	if err != nil {
		util.LogError(err)
		return 0
	}
	if fileName == nil {
		util.LogError(errors.New("file name not found"))
		return 0
	}
	return fileName.Id
}

func getPathItemId(prv db_entity.StatementProvider, path string) int64 {
	pathItem, err := LoadPathItem(prv, path, true)
	if err != nil {
		util.LogError(err)
		return 0
	}
	if pathItem == nil {
		util.LogError(errors.New("path not found"))
		return 0
	}
	return pathItem.Id
}

func getOwnerId(prv db_entity.StatementProvider, name string) int64 {
	owner, err := LoadOwner(prv, name, true)
	if err != nil {
		util.LogError(err)
		return 0
	}
	if owner == nil {
		util.LogError(errors.New("path not found"))
		return 0
	}
	return owner.Id
}

func RemoveAsset(prv db_entity.StatementProvider, assetId int64, pathId int64) (int, error) {

	if pathId > 0 {
		if err := RemoveOriginsByAssetIdAndPathId(prv, assetId, pathId); err != nil {
			return 9999, err
		}
	} else {
		if err := RemoveOriginsByAssetId(prv, assetId); err != nil {
			return 9999, err
		}
	}

	remainingOrigins, err := LoadOriginsByAssetId(prv, assetId)
	if err != nil {
		return 9999, err
	}
	if len(*remainingOrigins) > 0 {
		fmt.Printf("Keep asset, it has more origins: %d\n", assetId)
		return len(*remainingOrigins), nil
	}

	stmt, err := prv.Prepare("DELETE FROM asset WHERE id=?;")
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

func (a *Asset) GetSelectQuery() string {
	query := "SELECT id, hash, mimeType, fileTime, name FROM asset"
	where := ""
	if a.Hash != "" {
		where = util.JoinStrings(" AND ", where, "hash = ?")
	}
	if a.Id != 0 {
		where = util.JoinStrings(" AND ", where, "id = ?")
	}
	if where != "" {
		query = query + " WHERE " + where
	} else {
		fmt.Printf("Warn: No filter defined for asset\n")
	}
	return query
}

func (a *Asset) GetSelectQueryArgs() []any {
	args := make([]any, 0)
	if a.Hash != "" {
		args = append(args, a.Hash)
	}
	if a.Id != 0 {
		args = append(args, a.Id)
	}
	return args
}

func (a *Asset) Scan(rows *sql.Rows) error {
	return rows.Scan(&a.Id, &a.Hash, &a.MimeType, &a.FileTime, &a.Name)
}

func (a *Asset) GetInsertQuery() string {
	return "INSERT INTO asset(hash, mimeType, fileTime, name) VALUES(?,?,?,?);"
}

func (a *Asset) ExecInsert(stmt *sql.Stmt) (sql.Result, error) {
	return stmt.Exec(&a.Hash, &a.MimeType, &a.FileTime, &a.Name)
}

func (a *Asset) GetUpdateQuery() string {
	return "UPDATE asset SET hash=?, mimeType=?, fileTime=?, name=? WHERE id = ?;"
}

func (a *Asset) ExecUpdate(stmt *sql.Stmt) (sql.Result, error) {
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
