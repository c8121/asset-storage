package metadata_db_conn

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/c8121/asset-storage/internal/util"
)

var (
	db *sql.DB
)

func SetDatabase(database *sql.DB) {
	db = database
}

func GetDatabase() *sql.DB {
	return db
}

func CloseDatabase() {
	if db != nil {
		util.LogError(db.Close())
		db = nil
	}
}

func BeginTransaction() (*sql.Tx, error) {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func CommitOrLog(tx *sql.Tx) error {
	err := tx.Commit()
	if err != nil {
		fmt.Println(fmt.Errorf("commit failed: %v", err))
	}
	return err
}

func RollbackOrLog(tx *sql.Tx) {
	err := tx.Rollback()
	if err != nil {
		//fmt.Println(fmt.Errorf("rollback failed: %T, %v", err, err))
	}
}
