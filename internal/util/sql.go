package util

import (
	"database/sql"
	"fmt"
)

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
		//fmt.Println(fmt.Errorf("Rollback failed: %T, %v", err, err))
	}
}
