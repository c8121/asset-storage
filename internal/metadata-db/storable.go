package metadata_db

import (
	"database/sql"
	"errors"

	"github.com/c8121/asset-storage/internal/util"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrNotSelectable = errors.New("not a scanable")
	ErrNotInsertable = errors.New("not a insertable")
)

type Selectable interface {
	GetSelectQuery() string
	GetSelectQueryArgs() []any
	Scan(rows *sql.Rows) error
}

type Insertable interface {
	GetInsertQuery() string
	Exec(stmt *sql.Stmt) (sql.Result, error)
	SetId(int64)
}

type Updateable interface {
	GetUpdateQuery() string
	GetUpdateQueryArgs() []any
	Exec(stmt *sql.Stmt) (sql.Result, error)
}

func Create(tx *sql.Tx, insertIfNotExists bool, o any) error {

	scanable, ok := o.(Selectable)
	if !ok {
		return ErrNotSelectable
	}

	err := Load(tx, scanable)
	if err == ErrNotFound {
		if insertIfNotExists {

			insertable, ok := o.(Insertable)
			if !ok {
				return ErrNotInsertable
			}

			err = Insert(tx, insertable)
			if err != nil {
				return err
			}

		} else {
			return ErrNotFound
		}
	} else if err != nil {
		return err
	}

	return nil
}

func Load(tx *sql.Tx, o Selectable) error {

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

func Update(tx *sql.Tx, o Updateable) error {

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
