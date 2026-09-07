package db_entity

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/c8121/asset-storage/internal/util"
)

var (
	ErrNotFound = errors.New("not found")
)

type StatementProvider interface {
	Prepare(query string) (*sql.Stmt, error)
}

type Loadable interface {
	GetSelectQuery() string
	GetSelectQueryArgs() []any
	Scan(rows *sql.Rows) error
}

// Load selects all data from database and applies to given Loadable
func Load(tx StatementProvider, o Loadable) error {

	if o == nil {
		return errors.New("object is nil")
	}

	stmt, err := tx.Prepare(o.GetSelectQuery())
	if err != nil {
		return err
	}
	defer util.CloseOrLog(stmt)

	if rows, err := stmt.Query(o.GetSelectQueryArgs()...); err == nil {
		defer util.CloseOrLog(rows)
		if rows.Next() {
			if err := o.Scan(rows); err != nil {
				fmt.Printf("Error scanning rows for %T: %s\n", o, err)
				return err
			}
		} else {
			return ErrNotFound
		}

	} else {
		return err
	}

	return nil
}
