package metadata_db_entity

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/c8121/asset-storage/internal/db_entity"
	"github.com/c8121/asset-storage/internal/metadata_db_conn"
	"github.com/c8121/asset-storage/internal/util"
)

type PathItem struct {
	Id     int64
	Parent int64
	Name   string
}

var (
	pathItemCache     map[string]*PathItem
	loadPathChunkSize = 100
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

// LoadPathItem gets PathItem from db, splits path and searches
func LoadPathItem(prv db_entity.StatementProvider, path string, createIfNotExists bool) (*PathItem, error) {

	names := SplitPath(path)
	var parent int64 = 0

	var pathItem *PathItem
	for _, name := range names {

		cacheKey := fmt.Sprintf("%d/%s", parent, name)
		cachedItem, ok := pathItemCache[cacheKey]
		if !ok {
			pathItem = &PathItem{Parent: parent, Name: name}
			err := db_entity.LoadEntity(prv, createIfNotExists, pathItem)
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

// FindPathIds looks for origin-paths in database and
// returns it as map: asset.Id -> [path.Id,...]
func FindPathIds(prv db_entity.StatementProvider, assetIds []int64) (map[int64][]int64, error) {

	var idMap = make(map[int64][]int64, len(assetIds))
	var chunk = make([]int64, 0)

	for i, assetId := range assetIds {

		chunk = append(chunk, assetId)
		if len(chunk) == loadPathChunkSize || (i+1) == len(assetIds) {

			var query = "SELECT a.id, o.path" +
				" FROM origin o " +
				" INNER JOIN asset a ON o.asset = a.id "

			query += " WHERE a.id in(" +
				strings.Repeat("?,", len(chunk)-1) + "?" +
				");"

			var params = make([]any, 0)
			for _, id := range chunk {
				params = append(params, id)
			}

			stmt, err := prv.Prepare(query)
			if err != nil {
				return nil, err
			}
			defer util.CloseOrLog(stmt)

			if rows, err := stmt.Query(params...); err == nil {
				defer util.CloseOrLog(rows)
				for rows.Next() {
					var assetId int64
					var pathId int64
					if err := rows.Scan(&assetId, &pathId); err != nil {
						return nil, err
					}

					list, exists := idMap[assetId]
					if exists {
						idMap[assetId] = append(list, pathId)
					} else {
						list = make([]int64, 0)
						idMap[assetId] = append(list, pathId)
					}
				}

			} else {
				return nil, err
			}

			chunk = make([]int64, 0)
		}
	}

	return idMap, nil
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

func (p *PathItem) GetSelectQuery() string {
	query := "SELECT id, parent, name FROM pathItem"
	where := ""
	if p.Name != "" {
		where = util.JoinStrings(" AND ", where, "name = ?")
	}
	if p.Parent != 0 {
		where = util.JoinStrings(" AND ", where, "parent = ?")
	}
	if p.Id != 0 {
		where = util.JoinStrings(" AND ", where, "id = ?")
	}
	if where != "" {
		query = query + " WHERE " + where
	} else {
		fmt.Printf("Warn: No filter defined for path-item\n")
	}
	return query
}

func (p *PathItem) GetSelectQueryArgs() []any {
	args := make([]any, 0)
	if p.Name != "" {
		args = append(args, p.Name)
	}
	if p.Parent != 0 {
		args = append(args, p.Parent)
	}
	if p.Id != 0 {
		args = append(args, p.Id)
	}
	return args
}

func (p *PathItem) Scan(rows *sql.Rows) error {
	return rows.Scan(&p.Id, &p.Parent, &p.Name)
}

func (p *PathItem) GetInsertQuery() string {
	return "INSERT INTO pathItem(parent, name) VALUES(?,?);"
}

func (p *PathItem) ExecInsert(stmt *sql.Stmt) (sql.Result, error) {
	return stmt.Exec(&p.Parent, &p.Name)
}

func (p *PathItem) GetUpdateQuery() string {
	return "UPDATE pathItem SET parent=?, name=? WHERE id = ?;"
}

func (p *PathItem) ExecUpdate(stmt *sql.Stmt) (sql.Result, error) {
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
