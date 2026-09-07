package db_entity

import (
	"database/sql"
	"errors"

	"github.com/c8121/asset-storage/internal/util"
)

type Updatable interface {
	GetUpdateQuery() string
	GetUpdateQueryArgs() []any
	Exec(stmt *sql.Stmt) (sql.Result, error)
}

// Update updates an existing record in database with contents of given Insertable
func Update(tx *sql.Tx, o Updatable) error {

	if o == nil {
		return errors.New("object is nil")
	}

	stmt, err := tx.Prepare(o.GetUpdateQuery())
	if err != nil {
		return err
	}
	defer util.CloseOrLog(stmt)

	_, err = o.Exec(stmt)
	if err != nil {
		return err
	}

	return nil
}
