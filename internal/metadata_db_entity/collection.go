package metadata_db_entity

import (
	"database/sql"
	"time"
)

type Collection struct {
	Id      int64
	Hash    string
	Name    string
	Created time.Time
}

func (c *Collection) GetId() int64 {
	return c.Id
}

func (c *Collection) Load() error {
	return Load(c)
}

func (c *Collection) Save() error {
	return Save(c)
}

func (c *Collection) Get(insertIfNotExists bool) error {
	return Get(insertIfNotExists, c)
}

func (c *Collection) GetSelectQuery() string {
	return "SELECT id, hash, name, created FROM collection WHERE hash = ?;"
}

func (c *Collection) GetSelectQueryArgs() []any {
	return []any{c.Hash}
}

func (c *Collection) Scan(rows *sql.Rows) error {
	return rows.Scan(&c.Id, &c.Hash, &c.Name, &c.Created)
}

func (c *Collection) GetInsertQuery() string {
	return "INSERT INTO collection(hash, name, created) VALUES(?,?,?);"
}

func (c *Collection) GetUpdateQuery() string {
	return "UPDATE asset SET hash=?, name=?, created=? WHERE id = ?;"
}

func (c *Collection) GetUpdateQueryArgs() []any {
	return []any{&c.Hash, &c.Name, &c.Created}
}

func (c *Collection) Exec(stmt *sql.Stmt) (sql.Result, error) {
	return stmt.Exec(&c.Hash, &c.Name, &c.Created, &c.Id)
}

func (c *Collection) SetId(id int64) {
	c.Id = id
}

func (c *Collection) GetCreateQueries() []string {
	return []string{
		"CREATE TABLE IF NOT EXISTS collection(id integer PRIMARY KEY, hash TEXT(64), name TEXT(1024), created DATETIME);",
		"CREATE INDEX IF NOT EXISTS idx_collection_hash on collection(hash);",
		"CREATE INDEX IF NOT EXISTS idx_collection_name on collection(name);",
		"CREATE INDEX IF NOT EXISTS idx_collection_created on collection(created);",
	}
}
