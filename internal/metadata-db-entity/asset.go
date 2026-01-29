package metadata_db_entity

import (
	"database/sql"
	"time"
)

type Asset struct {
	Id       int64
	Hash     string
	MimeType int64
	FileTime time.Time //Max of all origins
	Name     int64     //Latest name
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
