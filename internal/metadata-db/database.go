package metadata_db

import (
	"database/sql"

	"github.com/c8121/asset-storage/internal/util"
)

var (
	db *sql.DB
)

func SetDatabase(databse *sql.DB) {
	db = databse
	dbInitMimeType()
	dbInitAsset()
}

// dbInitExec Execute DDL
func dbInitExec(ddl string) {
	_, err := db.Exec(ddl)
	util.PanicOnError(err, "Failed to init database")
}
