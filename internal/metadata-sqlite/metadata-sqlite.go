package mdsqlite

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"time"

	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/metadata"
	"github.com/c8121/asset-storage/internal/util"
	_ "modernc.org/sqlite"
)

var (
	DB *sql.DB
)

type (
	AssetListItem struct {
		Hash     string
		Name     string
		FileTime time.Time
	}

	AssetListFilter struct {
		MimeType string
	}
)

// Open Connect to SQLite database file + init
func Open() {

	dbDir := filepath.Dir(config.AssetMetaDataDb)
	util.CreateDirIfNotExists(dbDir, metadata.FilePermissions)

	fmt.Printf("Open DB %s\n", config.AssetMetaDataDb)
	db, err := sql.Open("sqlite", config.AssetMetaDataDb)
	util.PanicOnError(err, "Failed to open sqlite database: "+config.AssetMetaDataDb)

	DB = db

	initDatabase()
}

// Close Disconnect from Database
func Close() {
	fmt.Printf("Close DB %s\n", config.AssetMetaDataDb)
	util.LogError(DB.Close())
}

// ListAssets returns an array of AssetListItem, sorted by date desc
func ListAssets(offset, count int, filter *AssetListFilter) ([]AssetListItem, error) {

	var query = "SELECT o.hash, name, filetime FROM origin o INNER JOIN asset a ON o.hash = a.hash "
	var where = ""
	var limit = "ORDER BY filetime DESC, o.hash ASC LIMIT ? OFFSET ?;"

	var params = make([]any, 0)

	if filter != nil && filter.MimeType != "" {
		params = append(params, filter.MimeType)
		where += "(mimetype LIKE ?)"
	}

	params = append(params, count)
	params = append(params, offset)

	if where != "" {
		query += "WHERE " + where + " "
	}
	query += limit
	//fmt.Printf("Query: %s %v\n", query, params)

	stmt, err := DB.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer util.CloseOrLog(stmt)

	var items []AssetListItem

	if rows, err := stmt.Query(params...); err == nil {
		defer util.CloseOrLog(rows)
		for rows.Next() {
			var item AssetListItem
			if err := rows.Scan(&item.Hash, &item.Name, &item.FileTime); err != nil {
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
		return fmt.Errorf("failed to prepare: %w", err)
	}
	defer util.CloseOrLog(stmt)

	_, err = stmt.Exec(hash, meta.MimeType)
	if err != nil {
		return fmt.Errorf("failed to execute: %w", err)
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
		return fmt.Errorf("failed to prepare: %w", err)
	}
	defer util.CloseOrLog(stmt)

	_, err = stmt.Exec(hash, origin.Name, origin.Path, origin.Owner, origin.FileTime)
	if err != nil {
		return fmt.Errorf("failed to execute: %w", err)
	}

	return nil
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
	fmt.Printf("Init DB %s\n", config.AssetMetaDataDb)
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
	util.PanicOnError(err, "Failed to init database")
}
