package metadata_db_entity

import (
	"database/sql"
	"fmt"

	"github.com/c8121/asset-storage/internal/db_entity"
	"github.com/c8121/asset-storage/internal/util"
)

type FileName struct {
	Id   int64
	Name string
}

// LoadFileName gets FileName from db, created it required
func LoadFileName(tx *sql.Tx, name string, createIfNotExists bool) (*FileName, error) {

	var fileName = &FileName{Name: name}
	err := db_entity.LoadEntity(tx, createIfNotExists, fileName)
	if err != nil {
		return nil, err
	}

	return fileName, nil
}

func (n *FileName) GetId() int64 {
	return n.Id
}

func (n *FileName) GetSelectQuery() string {
	query := "SELECT id, name FROM fileName"
	where := ""
	if n.Name != "" {
		where = util.JoinStrings(" AND ", where, "name = ?")
	}
	if n.Id != 0 {
		where = util.JoinStrings(" AND ", where, "id = ?")
	}
	if where != "" {
		query = query + " WHERE " + where
	} else {
		fmt.Printf("Warn: No filter defined for filename\n")
	}
	return query
}

func (n *FileName) GetSelectQueryArgs() []any {
	args := make([]any, 0)
	if n.Name != "" {
		args = append(args, n.Name)
	}
	if n.Id != 0 {
		args = append(args, n.Id)
	}
	return args
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
