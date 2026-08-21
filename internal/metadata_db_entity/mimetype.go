package metadata_db_entity

import (
	"context"
	"database/sql"
	"strings"

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

	var query = "SELECT id, name FROM mimeType ORDER BY name, id;"

	stmt, err := db.Prepare(query)
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

func GetMimeType(name string, createIfNotExists bool) (*MimeType, error) {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer util.RollbackOrLog(tx)

	mimeType, err := GetMimeTypeTx(tx, name, createIfNotExists)
	if err != nil {
		return nil, err
	}

	if err = util.CommitOrLog(tx); err != nil {
		return nil, err
	}

	return mimeType, nil
}

func GetMimeTypeTx(tx *sql.Tx, name string, createIfNotExists bool) (*MimeType, error) {
	name = NormalizeName(name)
	mimeType, ok := mimeTypeCache[name]
	if ok {
		return mimeType, nil
	}

	mimeType = &MimeType{Name: name}
	err := GetTx(tx, createIfNotExists, mimeType)
	if err == nil {
		mimeTypeCache[name] = mimeType
	}

	return mimeType, err
}

func (m *MimeType) GetId() int64 {
	return m.Id
}

func (m *MimeType) Save() error {
	return Save(m)
}

func (m *MimeType) GetSelectQuery() string {
	return "SELECT id, name FROM mimeType WHERE name = ?;"
}

func (m *MimeType) GetSelectQueryArgs() []any {
	return []any{m.Name}
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
