package mdsqlite

import (
	"fmt"

	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/util"
)

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

	dbInitExec("PRAGMA journal_mode = WAL")
}

// dbInitExec Execute DDL
func dbInitExec(ddl string) {
	_, err := DB.Exec(ddl)
	util.PanicOnError(err, "Failed to init database")
}
