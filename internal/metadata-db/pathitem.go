package metadata_db

import (
	"context"
	"database/sql"
	"strings"
)

type PathItem struct {
	Id     int64
	Parent int64
	Name   string
}

func SplitPath(path string) []string {

	l := len(path)
	if l == 0 {
		return []string{""}
	}

	names := make([]string, 0)

	s := 0
	i := 0
	if strings.HasPrefix(strings.ToLower(path), "file:") {
		s = 5
		i = 5
	}

	for ; i < l; i++ {
		c := path[i]
		if c == '/' || c == '\\' {
			if i-s > 0 {
				names = append(names, path[s:i])
			}
			s = i + 1
		}
	}
	if i > s {
		names = append(names, path[s:])
	}

	return names
}

func GetPathItem(path string, createIfNotExists bool) (*PathItem, error) {

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	pathItem, err := GetPathItemTx(tx, path, createIfNotExists)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return pathItem, nil
}

func GetPathItemTx(tx *sql.Tx, path string, createIfNotExists bool) (*PathItem, error) {

	names := SplitPath(path)
	var parent int64 = 0

	var pathItem *PathItem
	for _, name := range names {
		pathItem = &PathItem{Parent: parent, Name: name}
		err := GetTx(tx, createIfNotExists, pathItem)
		if err == ErrNotFound {
			return nil, err
		}
		if err != nil {
			return nil, err
		}
		parent = pathItem.Id
	}

	return pathItem, nil
}

func (p *PathItem) GetId() int64 {
	return p.Id
}

func (p *PathItem) Save() error {
	return Save(p)
}

func (p *PathItem) GetSelectQuery() string {
	return "SELECT id, parent, name FROM pathItem WHERE parent = ? and name = ?;"
}

func (p *PathItem) GetSelectQueryArgs() []any {
	return []any{p.Parent, p.Name}
}

func (p *PathItem) Scan(rows *sql.Rows) error {
	return rows.Scan(&p.Id, &p.Parent, &p.Name)
}

func (p *PathItem) GetInsertQuery() string {
	return "INSERT INTO pathItem(parent, name) VALUES(?,?);"
}

func (p *PathItem) GetUpdateQuery() string {
	return "UPDATE pathItem SET parent=?, name=? WHERE id = ?;"
}

func (p *PathItem) GetUpdateQueryArgs() []any {
	return []any{&p.Parent, &p.Name, &p.Id}
}

func (p *PathItem) Exec(stmt *sql.Stmt) (sql.Result, error) {
	return stmt.Exec(&p.Parent, &p.Name, &p.Id)
}

func (p *PathItem) SetId(id int64) {
	p.Id = id
}

func dbInitPathItem() {
	dbInitExec("CREATE TABLE IF NOT EXISTS pathItem(id integer PRIMARY KEY, parent integer, name TEXT(1024));")
	dbInitExec("CREATE INDEX IF NOT EXISTS idx_pathItem_parent on pathItem(parent);")
	dbInitExec("CREATE INDEX IF NOT EXISTS idx_pathItem_name on pathItem(name);")
}
