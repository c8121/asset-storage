package metadata_db_entity

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/c8121/asset-storage/internal/metadata_db_conn"
	"github.com/c8121/asset-storage/internal/util"
)

type PathItem struct {
	Id     int64
	Parent int64
	Name   string
}

var (
	pathItemCache map[string]*PathItem
)

func init() {
	pathItemCache = make(map[string]*PathItem)
}

// SplitPath split at \ and /, remove 'file:'
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

	if len(names) == 0 {
		return []string{""}
	}

	return names
}

// GetPathItem gets PathItem from db, splits path and searches
func GetPathItem(path string, createIfNotExists bool) (*PathItem, error) {

	tx, err := metadata_db_conn.BeginTransaction()
	if err != nil {
		return nil, err
	}

	pathItem, err := GetPathItemTx(tx, path, createIfNotExists)
	if err != nil {
		return nil, err
	}

	if err = metadata_db_conn.CommitOrLog(tx); err != nil {
		return nil, err
	}

	return pathItem, nil
}

// GetPathItemIdTx gets PathItem-ID from db, splits path and searches
func GetPathItemIdTx(tx *sql.Tx, path string, createIfNotExists bool) int64 {
	pathItem, err := GetPathItemTx(tx, path, createIfNotExists)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return pathItem.Id
}

// GetPathItemTx gets PathItem from db, splits path and searches
func GetPathItemTx(tx *sql.Tx, path string, createIfNotExists bool) (*PathItem, error) {

	names := SplitPath(path)
	var parent int64 = 0

	var pathItem *PathItem
	for _, name := range names {

		cacheKey := fmt.Sprintf("%d/%s", parent, name)
		cachedItem, ok := pathItemCache[cacheKey]
		if !ok {
			pathItem = &PathItem{Parent: parent, Name: name}
			err := Get(tx, createIfNotExists, pathItem)
			if errors.Is(err, ErrNotFound) {
				return nil, err
			}
			if err != nil {
				return nil, err
			}
			pathItemCache[cacheKey] = pathItem
		} else {
			pathItem = cachedItem
			//fmt.Printf("Cached: %s = %v\n", cacheKey, pathItem)
		}

		parent = pathItem.Id
	}

	return pathItem, nil
}

func GetAssetIdsFromPath(pathId int64) ([]int64, error) {

	list := make([]int64, 0)

	stmt, err := metadata_db_conn.GetDatabase().Prepare("SELECT DISTINCT asset FROM origin where path = ?;")
	if err != nil {
		return nil, err
	}
	defer util.CloseOrLog(stmt)

	if rows, err := stmt.Query(pathId); err == nil {
		defer util.CloseOrLog(rows)
		for rows.Next() {
			var assetId int64
			if err := rows.Scan(&assetId); err != nil {
				fmt.Printf("Error scanning rows: %s\n", err)
				return nil, err
			} else {
				list = append(list, assetId)
			}
		}

	} else {
		return nil, err
	}

	return list, err

}

func RemovePathIfEmpty(pathId int64) error {

	tx, err := metadata_db_conn.BeginTransaction()
	if err != nil {
		return err
	}

	cnt := countPathRefs(tx, pathId, "SELECT COUNT(*) FROM pathItem where parent = ?;")
	if cnt > 0 {
		fmt.Printf("Will not remove path %d, it has %d children\n", pathId, cnt)
		return nil
	}

	cnt = countPathRefs(tx, pathId, "SELECT COUNT(*) FROM origin where path = ?;")
	if cnt > 0 {
		fmt.Printf("Will not remove path %d, it is referenced in %d origins\n", pathId, cnt)
		return nil
	}

	fmt.Printf("Removing path %d\n", pathId)
	stmt, err := tx.Prepare("DELETE FROM pathItem where id = ?;")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(pathId)
	if err != nil {
		return err
	}

	return metadata_db_conn.CommitOrLog(tx)
}

func countPathRefs(tx *sql.Tx, pathId int64, sql string) int64 {

	stmt, err := tx.Prepare(sql)
	if err != nil {
		return 99999
	}
	defer util.CloseOrLog(stmt)

	if rows, err := stmt.Query(pathId); err == nil {
		defer util.CloseOrLog(rows)
		if rows.Next() {
			var count int64
			if err := rows.Scan(&count); err != nil {
				fmt.Printf("Error scanning rows: %s\n", err)
				return 99999
			}

			return count
		}
	}

	return 99999
}

func (p *PathItem) GetId() int64 {
	return p.Id
}

func (p *PathItem) Save(tx *sql.Tx) error {
	return Save(tx, p)
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

func (a *PathItem) GetCreateQueries() []string {
	return []string{
		"CREATE TABLE IF NOT EXISTS pathItem(id integer PRIMARY KEY, parent integer, name TEXT(1024));",
		"CREATE INDEX IF NOT EXISTS idx_pathItem_parent on pathItem(parent);",
		"CREATE INDEX IF NOT EXISTS idx_pathItem_name on pathItem(name);",
	}
}
