package metadata_db

import (
	"database/sql"
	"fmt"

	"github.com/c8121/asset-storage/internal/util"
)

var (
	db *sql.DB
)

func SetDatabase(databse *sql.DB) {
	db = databse
	dbInitMimeType()
	dbInitPathItem()
	dbInitOwner()
	dbInitAsset()
	dbInitOrigin()
}

func CloseDatabase() {
	if db != nil {
		util.LogError(db.Close())
		db = nil
	}
}

func commitOrLog(tx *sql.Tx) error {
	err := tx.Commit()
	if err != nil {
		fmt.Println(fmt.Errorf("commit failed: %v", err))
	}
	return err
}

func rollbackOrLog(tx *sql.Tx) {
	err := tx.Rollback()
	if err != nil {
		//fmt.Println(fmt.Errorf("Rollback failed: %T, %v", err, err))
	}
}

// dbInitExec Execute DDL
func dbInitExec(ddl string) {
	_, err := db.Exec(ddl)
	util.PanicOnError(err, "Failed to init database")
}
