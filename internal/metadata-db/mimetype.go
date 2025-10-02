package metadata_db

import (
	"database/sql"
	"fmt"
	"strings"
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

func GetMimeType(name string, createIfNotExists bool) (*MimeType, error) {
	name = NormalizeName(name)
	if len(name) == 0 {
		return nil, fmt.Errorf("invalid Mime-Type: Empty name")
	}

	mimeType, ok := mimeTypeCache[name]
	if ok {
		return mimeType, nil
	}

	mimeType = &MimeType{Name: name}
	err := Get(createIfNotExists, mimeType)
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

func dbInitMimeType() {
	dbInitExec("CREATE TABLE IF NOT EXISTS mimeType(id integer PRIMARY KEY, name TEXT(32));")
}
