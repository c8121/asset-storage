package metadata_db

import (
	"context"
	"database/sql"
	"time"
)

type Asset struct {
	Id       int64
	Hash     string
	MimeType int64
	FileTime time.Time //Max of all origins
}

func (a *Asset) Load() error {

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = Load(tx, a)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (a *Asset) Save() error {

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if a.Id == 0 {
		err = Insert(tx, a)
	} else {
		err = Update(tx, a)
	}
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (a *Asset) Create(insertIfNotExists bool) error {

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = Create(tx, insertIfNotExists, a)
	if err != nil {
		return err
	}

	return tx.Commit()
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
	return "UPDATE asset SET hash=?, mimeType=?, fileTime=? WHERE ID = ?;"
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
