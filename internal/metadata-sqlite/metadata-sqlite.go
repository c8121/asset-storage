package mdsqlite

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/metadata"
	"github.com/c8121/asset-storage/internal/util"
	_ "modernc.org/sqlite"
)

var (
	DB *sql.DB
)

// Open Connect to SQLite database file + init
func Open() {

	dbDir := filepath.Dir(config.AssetMetaDataDb)
	util.CreateDirIfNotExists(dbDir, metadata.FilePermissions)

	url := "file:" + config.AssetMetaDataDb +
		"?_pragma=journal_mode(wal)" +
		"&_pragma=busy_timeout(500)" +
		"&_pragma=synchronous(normal)" +
		"&_txlock=immediate"

	fmt.Printf("Open DB %s\n", config.AssetMetaDataDb)
	db, err := sql.Open("sqlite", url)
	util.PanicOnError(err, "Failed to open sqlite database: "+config.AssetMetaDataDb)

	DB = db

	initDatabase()
}

// Close Disconnect from Database
func Close() {
	fmt.Printf("Close DB %s\n", config.AssetMetaDataDb)
	util.LogError(DB.Close())
}

// AddMetaData Upsert meta-data to database
func AddMetaData(hash string, meta *metadata.AssetMetadata) error {

	ctx := context.Background()
	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = AddMetaDataTx(tx, hash, meta)
	if err != nil {
		return err
	}

	return tx.Commit()

}

// AddMetaDataTx Upsert meta-data to database within a transaction
func AddMetaDataTx(tx *sql.Tx, hash string, meta *metadata.AssetMetadata) error {

	stmt, err := tx.Prepare("INSERT INTO asset(hash, mimetype) VALUES(?, ?) " +
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
		err = AddOriginTx(tx, hash, &origin)
		if err != nil {
			return err
		}
	}

	return nil

}

// AddOrigin Add Origin to database if not exists
func AddOrigin(hash string, origin *metadata.Origin) error {

	ctx := context.Background()
	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = AddOriginTx(tx, hash, origin)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// AddOriginTx Add Origin to database if not exists within a transaction
func AddOriginTx(tx *sql.Tx, hash string, origin *metadata.Origin) error {

	delErr := removeOrigin(tx, hash, origin)
	if delErr != nil {
		return delErr
	}

	stmt, err := tx.Prepare("INSERT INTO origin(hash, name, path, owner, filetime) VALUES(?, ?, ?, ?, ?);")
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
func removeOrigin(tx *sql.Tx, hash string, origin *metadata.Origin) error {

	stmt, err := tx.Prepare("DELETE FROM origin WHERE " +
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
