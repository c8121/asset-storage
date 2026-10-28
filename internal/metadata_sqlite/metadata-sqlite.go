package metadata_sqlite

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/metadata"
	"github.com/c8121/asset-storage/internal/metadata_db_entity"
	"github.com/c8121/asset-storage/internal/util"
	"github.com/c8121/go_db_entities"
	"github.com/c8121/go_db_entities/sqlite"
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

	metadata_db_entity.DbContext = go_db_entities.NewContext(db, sqlite.NewQueryBuilder())

	metadata_db_entity.CheckCreateTables()

	fmt.Println("DB Ready")
}

// Close Disconnect from Database
func Close() {
	fmt.Printf("Close DB %s\n", config.AssetMetaDataDb)
	util.CloseOrLog(metadata_db_entity.DbContext)
}
