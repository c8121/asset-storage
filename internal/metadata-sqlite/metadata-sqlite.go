package mdsqlite

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/c8121/asset-storage/internal/metadata"
	"github.com/c8121/asset-storage/internal/util"
	_ "modernc.org/sqlite"
)

var (
	DBFile = "/home/christianh/asset-storage-metadata.sqlite"
	DB     *sql.DB
)

type (
	AssetListItem struct {
		Hash     string
		FileTime time.Time
	}
)

// Open Open SQLite database file + init
func Open() {
	fmt.Printf("Open DB %s\n", DBFile)
	db, err := sql.Open("sqlite", DBFile)
	util.Check(err, "Failed to open sqlite database: "+DBFile)

	DB = db

	initDatabase()
}

// Close Close Database
func Close() {
	fmt.Printf("Close DB %s\n", DBFile)
	util.LogError(DB.Close())
}

func ListAssets(offset, count int) ([]AssetListItem, error) {

	stmt, err := DB.Prepare("SELECT hash, filetime FROM origin LIMIT ? OFFSET ?;")
	if err != nil {
		return nil, err
	}
	defer util.CloseOrLog(stmt)

	var items []AssetListItem

	if rows, err := stmt.Query(count, offset); err == nil {
		defer util.CloseOrLog(rows)
		for rows.Next() {
			fmt.Println(".")
			var item AssetListItem
			if err := rows.Scan(&item.Hash, &item.FileTime); err != nil {
				return items, err
			}
			items = append(items, item)
		}

	} else {
		return items, err
	}

	return items, nil
}

// AddMetaData Upsert meta-data to database
func AddMetaData(hash string, meta *metadata.AssetMetadata) error {

	stmt, err := DB.Prepare("INSERT INTO asset(hash, mimetype) VALUES(?, ?) " +
		"ON CONFLICT DO UPDATE SET mimetype=excluded.mimetype;")
	if err != nil {
		return err
	}
	defer util.LogError(stmt.Close())

	_, err = stmt.Exec(hash, meta.MimeType)
	if err != nil {
		return err
	}

	for _, origin := range meta.Origins {
		err = addOrigin(hash, &origin)
		if err != nil {
			return err
		}
	}

	return nil

}

// addOrigin Add Origin to database if not exists
func addOrigin(hash string, origin *metadata.Origin) error {

	delErr := removeOrigin(hash, origin)
	if delErr != nil {
		return delErr
	}

	stmt, err := DB.Prepare("INSERT INTO origin(hash, name, path, owner, filetime) VALUES(?, ?, ?, ?, ?);")
	if err != nil {
		return err
	}
	defer util.CloseOrLog(stmt)

	_, err = stmt.Exec(hash, origin.Name, origin.Path, origin.Owner, origin.FileTime)
	return err
}

// removeOrigin Remove Origin from database
func removeOrigin(hash string, origin *metadata.Origin) error {

	stmt, err := DB.Prepare("DELETE FROM origin WHERE " +
		"hash = ? " +
		"AND name = ? " +
		"AND path = ? " +
		"AND owner = ? " +
		"AND filetime = ?;")
	if err != nil {
		return err
	}
	defer util.CloseOrLog(stmt)

	_, err = stmt.Exec(hash, origin.Name, origin.Path, origin.Owner, origin.FileTime)
	return err
}

// initDatabase Create tables and indexes
func initDatabase() {
	fmt.Printf("Init DB %s\n", DBFile)
	dbInitExec("CREATE TABLE IF NOT EXISTS asset(hash TEXT(64) PRIMARY KEY, " +
		"mimetype TEXT(128));")
	dbInitExec("CREATE TABLE IF NOT EXISTS origin(hash TEXT(64), " +
		"name TEXT(1000), path TEXT(4000), owner TEXT(100), filetime DATETIME);")
	dbInitExec("CREATE INDEX IF NOT EXISTS idx_origin_hash on origin(hash);")
	dbInitExec("CREATE INDEX IF NOT EXISTS idx_origin_name on origin(name);")
	dbInitExec("CREATE INDEX IF NOT EXISTS idx_origin_owner on origin(owner);")
	dbInitExec("CREATE INDEX IF NOT EXISTS idx_origin_filetime on origin(filetime);")
}

// dbInitExec Execute DDL
func dbInitExec(ddl string) {
	_, err := DB.Exec(ddl)
	util.Check(err, "Failed to init database")
}
