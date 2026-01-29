package mdsqlite

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/metadata"
	metadata_db "github.com/c8121/asset-storage/internal/metadata-db"
	metadata_db_entity "github.com/c8121/asset-storage/internal/metadata-db-entity"
	"github.com/c8121/asset-storage/internal/util"
	_ "modernc.org/sqlite"
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

	metadata_db.SetDatabase(db)
	metadata_db_entity.AutoCreateEntities()
}

// Close Disconnect from Database
func Close() {
	fmt.Printf("Close DB %s\n", config.AssetMetaDataDb)
	metadata_db.CloseDatabase()
}
