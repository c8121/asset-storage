package metadata_db

import (
	"context"
	"database/sql"
	"strings"

	"github.com/c8121/asset-storage/internal/util"
)

type MimeType struct {
	Id   int64
	Name string
}

var (
	mimeTypeCache map[string]*MimeType
)

func init() {
	mimeTypeCache = make(map[string]*MimeType)
}

func DbGetMimeType(name string, createIfNotExists bool) (*MimeType, error) {

	ctx := context.Background()
	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	mimeType, err := DbGetMimeTypeTx(tx, name, createIfNotExists)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return mimeType, nil
}

func DbGetMimeTypeTx(tx *sql.Tx, name string, createIfNotExists bool) (*MimeType, error) {

	name = strings.ToLower(strings.TrimSpace(name))
	if p := strings.Index(name, ";"); p != -1 {
		name = strings.TrimSpace(name[:p])
	}

	mimeType, ok := mimeTypeCache[name]
	if ok {
		return mimeType, nil
	}

	stmt, err := tx.Prepare("SELECT id, name FROM mime_type WHERE name = ?;")
	if err != nil {
		return nil, err
	}
	defer util.CloseOrLog(stmt)

	if rows, err := stmt.Query(name); err == nil {
		defer util.CloseOrLog(rows)
		if rows.Next() {
			if err := rows.Scan(&mimeType.Id, &mimeType.Name); err != nil {
				return nil, err
			}
		} else {
			if !createIfNotExists {
				return nil, nil
			}
			mimeType, err = dbCreateMimeTypeTx(tx, name)
			if err != nil {
				return nil, err
			}
		}

	} else {
		return nil, err
	}

	mimeTypeCache[name] = mimeType

	return mimeType, nil
}

func dbCreateMimeTypeTx(tx *sql.Tx, name string) (*MimeType, error) {

	stmt, err := tx.Prepare("INSERT INTO mime_type(name) VALUES(?);")
	if err != nil {
		return nil, err
	}
	defer util.CloseOrLog(stmt)

	r, err := stmt.Exec(name)
	if err != nil {
		return nil, err
	}

	var mimeType MimeType

	mimeType.Id, err = r.LastInsertId()
	if err != nil {
		return nil, err
	}

	mimeType.Name = name

	return &mimeType, nil
}

func DbInitMimeType() {
	_, err := DB.Exec("CREATE TABLE IF NOT EXISTS mime_type(id integer PRIMARY KEY, name TEXT(32));")
	util.PanicOnError(err, "failed to create table: mime_type")
}
