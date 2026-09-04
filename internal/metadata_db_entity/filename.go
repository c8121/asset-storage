package metadata_db_entity

import (
	"database/sql"
	"errors"
	"fmt"
)

type FileName struct {
	Id   int64
	Name string
}

func GetFileNameId(tx *sql.Tx, name string, createIfNotExists bool) int64 {
	fileName, err := GetFileName(tx, name, createIfNotExists)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return fileName.Id
}

// GetFileName gets FileName from db
func GetFileName(tx *sql.Tx, name string, createIfNotExists bool) (*FileName, error) {

	var fileName = &FileName{Name: name}
	err := Get(tx, createIfNotExists, fileName)
	if errors.Is(err, ErrNotFound) {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return fileName, nil
}

func (n *FileName) GetId() int64 {
	return n.Id
}

func (n *FileName) Save(tx *sql.Tx) error {
	return Save(tx, n)
}

func (n *FileName) GetSelectQuery() string {
	return "SELECT id, name FROM fileName WHERE name = ?;"
}

func (n *FileName) GetSelectQueryArgs() []any {
	return []any{n.Name}
}

func (n *FileName) Scan(rows *sql.Rows) error {
	return rows.Scan(&n.Id, &n.Name)
}

func (n *FileName) GetInsertQuery() string {
	return "INSERT INTO fileName(name) VALUES(?);"
}

func (n *FileName) GetUpdateQuery() string {
	return "UPDATE fileName SET name=? WHERE id = ?;"
}

func (n *FileName) GetUpdateQueryArgs() []any {
	return []any{&n.Name, &n.Id}
}

func (n *FileName) Exec(stmt *sql.Stmt) (sql.Result, error) {
	return stmt.Exec(&n.Name, &n.Id)
}

func (n *FileName) SetId(id int64) {
	n.Id = id
}

func (a *FileName) GetCreateQueries() []string {
	return []string{
		"CREATE TABLE IF NOT EXISTS fileName(id integer PRIMARY KEY, name TEXT(1024));",
		"CREATE INDEX IF NOT EXISTS idx_fileName_name on fileName(name);",
	}
}
