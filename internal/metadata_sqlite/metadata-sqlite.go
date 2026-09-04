package metadata_sqlite

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/metadata"
	"github.com/c8121/asset-storage/internal/metadata_db_conn"
	"github.com/c8121/asset-storage/internal/metadata_db_entity"
	"github.com/c8121/asset-storage/internal/util"
	_ "gosqlite.org" // Bundles modernc + sqlite-vec
)

// Open Connect to SQLite database file + init
func Open() {

	dbDir := filepath.Dir(config.AssetMetaDataDb)
	util.CreateDirIfNotExists(dbDir, metadata.FilePermissions)

	url := "file:" + config.AssetMetaDataDb +
		"?_pragma=journal_mode(wal)" +
		"&_pragma=busy_timeout(2500)" +
		"&_pragma=synchronous(normal)" +
		"&_txlock=immediate"

	fmt.Printf("Open DB %s\n", config.AssetMetaDataDb)
	db, err := sql.Open("sqlite", url)
	util.PanicOnError(err, "Failed to open sqlite database: "+config.AssetMetaDataDb)

	metadata_db_conn.SetDatabase(db)
	metadata_db_entity.AutoCreateEntities()
}

// Close Disconnect from Database
func Close() {
	fmt.Printf("Close DB %s\n", config.AssetMetaDataDb)
	metadata_db_conn.CloseDatabase()
}
