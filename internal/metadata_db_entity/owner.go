package metadata_db_entity

import (
	"database/sql"
	"fmt"

	"github.com/c8121/asset-storage/internal/db_entity"
	"github.com/c8121/asset-storage/internal/util"
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
	err := db_entity.LoadEntity(tx, createIfNotExists, owner)
	if err == nil {
		OwnerCache[name] = owner
	}

	return owner, err
}

func (o *Owner) GetId() int64 {
	return o.Id
}

func (o *Owner) GetSelectQuery() string {
	query := "SELECT id, name FROM owner"
	where := ""
	if o.Name != "" {
		where = util.JoinStrings(" AND ", where, "name = ?")
	}
	if o.Id != 0 {
		where = util.JoinStrings(" AND ", where, "id = ?")
	}
	if where != "" {
		query = query + " WHERE " + where
	} else {
		fmt.Printf("Warn: No filter defined for owner\n")
	}
	return query
}

func (o *Owner) GetSelectQueryArgs() []any {
	args := make([]any, 0)
	if o.Name != "" {
		args = append(args, o.Name)
	}
	if o.Id != 0 {
		args = append(args, o.Id)
	}
	return args
}

func (o *Owner) Scan(rows *sql.Rows) error {
	return rows.Scan(&o.Id, &o.Name)
}

func (o *Owner) GetInsertQuery() string {
	return "INSERT INTO owner(name) VALUES(?);"
}

func (o *Owner) ExecInsert(stmt *sql.Stmt) (sql.Result, error) {
	return stmt.Exec(&o.Name)
}

func (o *Owner) GetUpdateQuery() string {
	return "UPDATE owner SET name=? WHERE id = ?;"
}

func (o *Owner) ExecUpdate(stmt *sql.Stmt) (sql.Result, error) {
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
