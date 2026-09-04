package metadata_db_entity

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/c8121/asset-storage/internal/metadata_db_conn"
	"github.com/c8121/asset-storage/internal/util"
)

var (
	ErrNotIdentifiable = errors.New("not implementing WithId")
	ErrNotFound        = errors.New("not found")
	ErrNotSelectable   = errors.New("not a Selectable")
	ErrNotInsertable   = errors.New("not a Insertable")
	ErrNotUpdatable    = errors.New("not a Updatable")
)

type StatementProvider interface {
	Prepare(query string) (*sql.Stmt, error)
}

type WithId interface {
	GetId() int64
}

type AutoCreatable interface {
	GetCreateQueries() []string
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

type Updatable interface {
	GetUpdateQuery() string
	GetUpdateQueryArgs() []any
	Exec(stmt *sql.Stmt) (sql.Result, error)
}

// AutoCreate executed DDL to create entity if not exists
func AutoCreate(o AutoCreatable) {
	db := metadata_db_conn.GetDatabase()
	queries := o.GetCreateQueries()
	for _, query := range queries {
		_, err := db.Exec(query)
		util.PanicOnError(err, "Failed to init entity: "+query)
	}
}

// Get first tries to Load(...), then Insert(...) if insertIfNotExists = true
func Get(tx *sql.Tx, insertIfNotExists bool, o any) error {

	scanable, ok := o.(Selectable)
	if !ok {
		return ErrNotSelectable
	}

	err := Load(tx, scanable)
	if errors.Is(err, ErrNotFound) {
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

// Load selects data from database and applies to given Selectable
func Load(tx StatementProvider, o Selectable) error {

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

// Save checks if object exists in database (GetId() != 0) and then does Insert or Update
func Save(tx *sql.Tx, o any) error {

	withId, ok := o.(WithId)
	if !ok {
		return ErrNotIdentifiable
	}

	var err error
	if withId.GetId() == 0 {
		insertable, ok := o.(Insertable)
		if !ok {
			return ErrNotInsertable
		}
		err = Insert(tx, insertable)
	} else {
		updatable, ok := o.(Updatable)
		if !ok {
			return ErrNotUpdatable
		}
		err = Update(tx, updatable)
	}

	return err
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

// Insert updates an existing record in database with contents of given Insertable
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
