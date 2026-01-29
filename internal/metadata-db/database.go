package metadata_db

import (
	"database/sql"

	metadata_db_entity "github.com/c8121/asset-storage/internal/metadata-db-entity"
	"github.com/c8121/asset-storage/internal/util"
)

var (
	db *sql.DB
)

func SetDatabase(databse *sql.DB) {
	db = databse
	metadata_db_entity.SetDatabase(databse)
}

func CloseDatabase() {
	if db != nil {
		util.LogError(db.Close())
		db = nil
	}
}
