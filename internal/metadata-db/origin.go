package metadata_db

import (
	"database/sql"
	"time"
)

type Origin struct {
	Id       int64
	Asset    int64
	Name     string
	Path     int64
	Owner    int64
	FileTime time.Time
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

func dbInitOrigin() {
	dbInitExec("CREATE TABLE IF NOT EXISTS origin(id integer PRIMARY KEY, asset INTEGER, name TEXT(1024), path INTEGER, owner INTEGER, fileTime DATETIME);")
	dbInitExec("CREATE INDEX IF NOT EXISTS idx_origin_asset on origin(asset);")
	dbInitExec("CREATE INDEX IF NOT EXISTS idx_origin_name on origin(name);")
	dbInitExec("CREATE INDEX IF NOT EXISTS idx_origin_path on origin(path);")
	dbInitExec("CREATE INDEX IF NOT EXISTS idx_origin_owner on origin(owner);")
	dbInitExec("CREATE INDEX IF NOT EXISTS idx_origin_fileTime on origin(fileTime);")
}
