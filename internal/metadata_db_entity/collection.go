package metadata_db_entity

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/c8121/asset-storage/internal/util"
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

func (c *Collection) GetSelectQuery() string {
	query := "SELECT id, hash, name, created FROM collection WHERE"
	where := ""
	if c.Hash != "" {
		where = util.JoinStrings(" AND ", where, "hash = ?")
	}
	if c.Id != 0 {
		where = util.JoinStrings(" AND ", where, "id = ?")
	}
	if where != "" {
		query = query + " WHERE " + where
	} else {
		fmt.Printf("Warn: No filter defined for collection\n")
	}
	return query
}

func (c *Collection) GetSelectQueryArgs() []any {
	args := make([]any, 0)
	if c.Hash != "" {
		args = append(args, c.Hash)
	}
	if c.Id != 0 {
		args = append(args, c.Id)
	}
	return args
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
