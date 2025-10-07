package metadata_db

import (
	"context"
	"database/sql"
	"fmt"
)

type Owner struct {
	Id   int64
	Name string
}

var (
	OwnerCache map[string]*Owner
)

func init() {
	OwnerCache = make(map[string]*Owner)
}

func GetOwnerIdTx(tx *sql.Tx, name string, createIfNotExists bool) int64 {
	owner, err := GetOwnerTx(tx, name, createIfNotExists)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return owner.Id
}

func GetOwner(name string, createIfNotExists bool) (*Owner, error) {

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer rollbackOrLog(tx)

	owner, err := GetOwnerTx(tx, name, createIfNotExists)
	if err != nil {
		return nil, err
	}

	if err = commitOrLog(tx); err != nil {
		return nil, err
	}

	return owner, nil
}

func GetOwnerTx(tx *sql.Tx, name string, createIfNotExists bool) (*Owner, error) {

	owner, ok := OwnerCache[name]
	if ok {
		return owner, nil
	}

	owner = &Owner{Name: name}
	err := GetTx(tx, createIfNotExists, owner)
	if err == nil {
		OwnerCache[name] = owner
	}

	return owner, err
}

func (o *Owner) GetId() int64 {
	return o.Id
}

func (o *Owner) Save() error {
	return Save(o)
}

func (o *Owner) GetSelectQuery() string {
	return "SELECT id, name FROM owner WHERE name = ?;"
}

func (o *Owner) GetSelectQueryArgs() []any {
	return []any{o.Name}
}

func (o *Owner) Scan(rows *sql.Rows) error {
	return rows.Scan(&o.Id, &o.Name)
}

func (o *Owner) GetInsertQuery() string {
	return "INSERT INTO owner(name) VALUES(?);"
}

func (o *Owner) GetUpdateQuery() string {
	return "UPDATE owner SET name=? WHERE id = ?;"
}

func (o *Owner) GetUpdateQueryArgs() []any {
	return []any{&o.Name, &o.Id}
}

func (o *Owner) Exec(stmt *sql.Stmt) (sql.Result, error) {
	return stmt.Exec(&o.Name, &o.Id)
}

func (o *Owner) SetId(id int64) {
	o.Id = id
}

func dbInitOwner() {
	dbInitExec("CREATE TABLE IF NOT EXISTS owner(id integer PRIMARY KEY, name TEXT(64));")
}
