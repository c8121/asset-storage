package metadata_db

import (
	"database/sql"
	"time"
)

type Asset struct {
	Id       int64
	Hash     string
	MimeType int64
	FileTime time.Time //Max of all origins
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
	return "SELECT id, hash, mimeType, fileTime FROM asset WHERE hash = ?;"
}

func (a *Asset) GetSelectQueryArgs() []any {
	return []any{a.Hash}
}

func (a *Asset) Scan(rows *sql.Rows) error {
	return rows.Scan(&a.Id, &a.Hash, &a.MimeType, &a.FileTime)
}

func (a *Asset) GetInsertQuery() string {
	return "INSERT INTO asset(hash, mimeType, fileTime) VALUES(?,?,?);"
}

func (a *Asset) GetUpdateQuery() string {
	return "UPDATE asset SET hash=?, mimeType=?, fileTime=? WHERE id = ?;"
}

func (a *Asset) GetUpdateQueryArgs() []any {
	return []any{&a.Hash, &a.MimeType, &a.FileTime, &a.Id}
}

func (a *Asset) Exec(stmt *sql.Stmt) (sql.Result, error) {
	return stmt.Exec(&a.Hash, &a.MimeType, &a.FileTime, &a.Id)
}

func (a *Asset) SetId(id int64) {
	a.Id = id
}

func dbInitAsset() {
	dbInitExec("CREATE TABLE IF NOT EXISTS asset(id integer PRIMARY KEY, hash TEXT(64), mimeType integer, fileTime DATETIME);")
	dbInitExec("CREATE INDEX IF NOT EXISTS idx_asset_hash on asset(hash);")
	dbInitExec("CREATE INDEX IF NOT EXISTS idx_asset_mimeType on asset(mimeType);")
	dbInitExec("CREATE INDEX IF NOT EXISTS idx_asset_fileTime on asset(fileTime);")
}
