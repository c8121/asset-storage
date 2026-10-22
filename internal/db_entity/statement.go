package db_entity

import "database/sql"

type StatementProvider interface {
	Prepare(query string) (*sql.Stmt, error)
}
