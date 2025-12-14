package metadata_db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/c8121/asset-storage/internal/util"
)

var (
	ErrNotIdentifyable = errors.New("not implementing WithId")
	ErrNotFound        = errors.New("not found")
	ErrNotSelectable   = errors.New("not a Selectable")
	ErrNotInsertable   = errors.New("not a Insertable")
	ErrNotUpdateable   = errors.New("not a Updateable")
)

type WithId interface {
	GetId() int64
}

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

// Get first tries to Load(...), then Insert(...) if insertIfNotExists = true
func Get(insertIfNotExists bool, o Selectable) error {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer rollbackOrLog(tx)

	err = GetTx(tx, insertIfNotExists, o)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetTx first tries to Load(...), then Insert(...) if insertIfNotExists = true
func GetTx(tx *sql.Tx, insertIfNotExists bool, o any) error {

	scanable, ok := o.(Selectable)
	if !ok {
		return ErrNotSelectable
	}

	err := LoadTx(tx, scanable)
	if errors.Is(err, ErrNotFound) {
		if insertIfNotExists {

			insertable, ok := o.(Insertable)
			if !ok {
				return ErrNotInsertable
			}

			err = InsertTx(tx, insertable)
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

// Load selects data from database and applies to given Selectable
func Load(o Selectable) error {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer commitOrLog(tx)

	return LoadTx(tx, o)
}

// LoadTx selects data from database and applies to given Selectable
func LoadTx(tx *sql.Tx, o Selectable) error {

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

// Save check if object exists in database (GetId() != 0) and then does Insert or Update
func Save(o any) error {

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer rollbackOrLog(tx)

	err = SaveTx(tx, o)
	if err != nil {
		return err
	}

	return commitOrLog(tx)
}

// Save checks it object exists in database (GetId() != 0) and then does Insert or Update
func SaveTx(tx *sql.Tx, o any) error {

	withId, ok := o.(WithId)
	if !ok {
		return ErrNotIdentifyable
	}

	var err error
	if withId.GetId() == 0 {
		insertable, ok := o.(Insertable)
		if !ok {
			return ErrNotInsertable
		}
		err = InsertTx(tx, insertable)
	} else {
		updateable, ok := o.(Updateable)
		if !ok {
			return ErrNotUpdateable
		}
		err = UpdateTx(tx, updateable)
	}

	return err
}

// InsertTx created new record in database with contents of given Insertable
func InsertTx(tx *sql.Tx, o Insertable) error {

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

// Insert updates an existing record in database with contents of given Insertable
func UpdateTx(tx *sql.Tx, o Updateable) error {

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
