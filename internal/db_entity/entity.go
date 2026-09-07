package db_entity

import (
	"database/sql"
	"errors"
)

var (
	ErrNotIdentifiable = errors.New("not implementing WithId")
	ErrNotLoadable     = errors.New("not a Loadable")
	ErrNotInsertable   = errors.New("not a Insertable")
	ErrNotUpdatable    = errors.New("not a Updatable")
)

type WithId interface {
	GetId() int64
}

// LoadEntity first tries to Load(...), then Insert(...) if insertIfNotExists = true
func LoadEntity(tx *sql.Tx, insertIfNotExists bool, o any) error {

	loadable, ok := o.(Loadable)
	if !ok {
		return ErrNotLoadable
	}

	err := Load(tx, loadable)
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

// SaveEntity checks if object exists in database (GetId() != 0) and then does Insert or Update
func SaveEntity(tx *sql.Tx, o any) error {

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
