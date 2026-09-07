package metadata_db_entity

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/c8121/asset-storage/internal/db_entity"
	"github.com/c8121/asset-storage/internal/metadata_db_conn"
	"github.com/c8121/asset-storage/internal/util"
)

type MimeType struct {
	Id   int64
	Name string
}

var (
	mimeTypeCache map[string]*MimeType
)

func init() {
	mimeTypeCache = make(map[string]*MimeType)
}

func NormalizeName(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	if p := strings.Index(name, ";"); p != -1 {
		name = strings.TrimSpace(name[:p])
	}
	return name
}

func ListMimeTypes() ([]MimeType, error) {

	stmt, err := metadata_db_conn.GetDatabase().Prepare("SELECT id, name FROM mimeType ORDER BY name, id;")
	if err != nil {
		return nil, err
	}
	defer util.CloseOrLog(stmt)

	var items []MimeType

	if rows, err := stmt.Query(); err == nil {
		defer util.CloseOrLog(rows)
		for rows.Next() {
			var item MimeType
			if err := rows.Scan(&item.Id, &item.Name); err != nil {
				return items, err
			}
			items = append(items, item)
		}

	} else {
		return items, err
	}

	return items, nil
}

func LoadMimeType(tx *sql.Tx, name string, createIfNotExists bool) (*MimeType, error) {
	name = NormalizeName(name)
	mimeType, ok := mimeTypeCache[name]
	if ok {
		return mimeType, nil
	}

	mimeType = &MimeType{Name: name}
	err := db_entity.LoadEntity(tx, createIfNotExists, mimeType)
	if err == nil {
		mimeTypeCache[name] = mimeType
	}

	return mimeType, err
}

func (m *MimeType) GetId() int64 {
	return m.Id
}

func (m *MimeType) GetSelectQuery() string {
	query := "SELECT id, name FROM mimeType"
	where := ""
	if m.Name != "" {
		where = util.JoinStrings(" AND ", where, "name = ?")
	}
	if m.Id != 0 {
		where = util.JoinStrings(" AND ", where, "id = ?")
	}
	if where != "" {
		query = query + " WHERE " + where
	} else {
		fmt.Printf("Warn: No filter defined for mimetype\n")
	}
	return query
}

func (m *MimeType) GetSelectQueryArgs() []any {
	args := make([]any, 0)
	if m.Name != "" {
		args = append(args, m.Name)
	}
	if m.Id != 0 {
		args = append(args, m.Id)
	}
	return args
}

func (m *MimeType) Scan(rows *sql.Rows) error {
	return rows.Scan(&m.Id, &m.Name)
}

func (m *MimeType) GetInsertQuery() string {
	return "INSERT INTO mimeType(name) VALUES(?);"
}

func (m *MimeType) GetUpdateQuery() string {
	return "UPDATE mimeType SET name=? WHERE id = ?;"
}

func (m *MimeType) GetUpdateQueryArgs() []any {
	return []any{&m.Name, &m.Id}
}

func (m *MimeType) Exec(stmt *sql.Stmt) (sql.Result, error) {
	return stmt.Exec(&m.Name, &m.Id)
}

func (m *MimeType) SetId(id int64) {
	m.Id = id
}

func (a *MimeType) GetCreateQueries() []string {
	return []string{
		"CREATE TABLE IF NOT EXISTS mimeType(id integer PRIMARY KEY, name TEXT(32));",
	}
}
