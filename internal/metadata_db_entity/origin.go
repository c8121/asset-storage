package metadata_db_entity

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/c8121/asset-storage/internal/util"
)

type Origin struct {
	Id       int64
	Asset    int64
	Name     int64
	Path     int64
	Owner    int64
	FileTime time.Time
}

func RemoveOriginsTx(tx *sql.Tx, asset *Asset) error {

	stmt, err := tx.Prepare("DELETE FROM origin WHERE asset = ?;")
	if err != nil {
		return err
	}
	defer util.CloseOrLog(stmt)

	_, err = stmt.Exec(asset.Id)
	return err
}

func RemoveOriginsByAssetIdTx(tx *sql.Tx, assetId int64) error {

	stmt, err := tx.Prepare("DELETE FROM origin WHERE asset = ?;")
	if err != nil {
		return err
	}
	defer util.CloseOrLog(stmt)

	_, err = stmt.Exec(assetId)
	return err
}

func RemoveOriginsByAssetIdAndPathIdTx(tx *sql.Tx, assetId int64, pathId int64) error {

	stmt, err := tx.Prepare("DELETE FROM origin WHERE asset = ? AND path = ?;")
	if err != nil {
		return err
	}
	defer util.CloseOrLog(stmt)

	_, err = stmt.Exec(assetId, pathId)
	return err
}

func GetOriginsTx(tx *sql.Tx, assetId int64) (*[]Origin, error) {

	stmt, err := tx.Prepare("SELECT id, asset, name, path, owner, fileTime FROM origin WHERE asset = ?;")
	if err != nil {
		return nil, err
	}
	defer util.CloseOrLog(stmt)

	list := &[]Origin{}

	if rows, err := stmt.Query(assetId); err == nil {
		defer util.CloseOrLog(rows)
		if rows.Next() {
			var origin Origin
			if err := origin.Scan(rows); err != nil {
				fmt.Printf("Error scanning rows: %s\n", err)
				return nil, err
			}
			*list = append(*list, origin)
		}

	} else {
		return nil, err
	}

	return list, nil
}

func (o *Origin) GetId() int64 {
	return o.Id
}

func (o *Origin) Load() error {
	return Load(o)
}

func (o *Origin) Save() error {
	return Save(o)
}

func (o *Origin) Get(insertIfNotExists bool) error {
	return Get(insertIfNotExists, o)
}

func (o *Origin) GetSelectQuery() string {
	return "SELECT id, asset, name, path, owner, fileTime FROM origin WHERE id = ?;"
}

func (o *Origin) GetSelectQueryArgs() []any {
	return []any{o.Id}
}

func (o *Origin) Scan(rows *sql.Rows) error {
	return rows.Scan(&o.Id, &o.Asset, &o.Name, &o.Path, &o.Owner, &o.FileTime)
}

func (o *Origin) GetInsertQuery() string {
	return "INSERT INTO origin(asset, name, path, owner, fileTime) VALUES(?,?,?,?,?);"
}

func (o *Origin) GetUpdateQuery() string {
	return "UPDATE origin SET asset=?, name=?, path=?, owner=?, fileTime=? WHERE id = ?;"
}

func (o *Origin) GetUpdateQueryArgs() []any {
	return []any{&o.Asset, &o.Name, &o.Path, &o.Owner, &o.FileTime, &o.Id}
}

func (o *Origin) Exec(stmt *sql.Stmt) (sql.Result, error) {
	return stmt.Exec(&o.Asset, &o.Name, &o.Path, &o.Owner, &o.FileTime, &o.Id)
}

func (o *Origin) SetId(id int64) {
	o.Id = id
}

func (a *Origin) GetCreateQueries() []string {
	return []string{
		"CREATE TABLE IF NOT EXISTS origin(id integer PRIMARY KEY, asset INTEGER, name INTEGER, path INTEGER, owner INTEGER, fileTime DATETIME);",
		"CREATE INDEX IF NOT EXISTS idx_origin_asset on origin(asset);",
		"CREATE INDEX IF NOT EXISTS idx_origin_name on origin(name);",
		"CREATE INDEX IF NOT EXISTS idx_origin_path on origin(path);",
		"CREATE INDEX IF NOT EXISTS idx_origin_owner on origin(owner);",
		"CREATE INDEX IF NOT EXISTS idx_origin_fileTime on origin(fileTime);",
	}
}
