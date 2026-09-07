package db_entity

import (
	"database/sql"
	"errors"

	"github.com/c8121/asset-storage/internal/util"
)

type Insertable interface {
	GetInsertQuery() string
	Exec(stmt *sql.Stmt) (sql.Result, error)
	SetId(int64)
}

// Insert creates new record in database with contents of given Insertable
func Insert(tx *sql.Tx, o Insertable) error {

	if o == nil {
		return errors.New("object is nil")
	}

	stmt, err := tx.Prepare(o.GetInsertQuery())
	if err != nil {
		return err
	}
	defer util.CloseOrLog(stmt)

	r, err := o.Exec(stmt)
	if err != nil {
		return err
	}

	id, err := r.LastInsertId()
	if err != nil {
		return err
	}
	o.SetId(id)

	return nil
}
