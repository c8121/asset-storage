package metadata_db_entity

import (
	"database/sql"

	"github.com/c8121/asset-storage/internal/db_entity"
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

func LoadOwner(tx *sql.Tx, name string, createIfNotExists bool) (*Owner, error) {

	owner, ok := OwnerCache[name]
	if ok {
		return owner, nil
	}

	owner = &Owner{Name: name}
	err := db_entity.LoadEntity(tx, createIfNotExists, owner, "name", name)
	if err == nil {
		OwnerCache[name] = owner
	}

	return owner, err
}

func (o *Owner) GetId() int64 {
	return o.Id
}

func (o *Owner) GetSelectQuery(filterColumn string) string {
	return "SELECT id, name FROM owner WHERE " + filterColumn + " = ?;"
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

func (a *Owner) GetCreateQueries() []string {
	return []string{
		"CREATE TABLE IF NOT EXISTS owner(id integer PRIMARY KEY, name TEXT(64));",
	}
}
