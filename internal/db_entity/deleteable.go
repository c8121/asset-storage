package db_entity

import (
	"database/sql"
	"errors"

	"github.com/c8121/asset-storage/internal/util"
)

type Deletable interface {
	GetDeleteQuery() string
	ExecDelete(stmt *sql.Stmt) (sql.Result, error)
}

// Delete removes an existing record in database
func Delete(prv StatementProvider, o Deletable) error {

	if o == nil {
		return errors.New("object is nil")
	}

	stmt, err := prv.Prepare(o.GetDeleteQuery())
	if err != nil {
		return err
	}
	defer util.CloseOrLog(stmt)

	_, err = o.ExecDelete(stmt)
	if err != nil {
		return err
	}

	return nil
}
