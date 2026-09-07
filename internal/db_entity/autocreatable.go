package db_entity

import (
	"github.com/c8121/asset-storage/internal/metadata_db_conn"
	"github.com/c8121/asset-storage/internal/util"
)

type AutoCreatable interface {
	GetCreateQueries() []string
}

// AutoCreate executes DDL to create entity if not exists
func AutoCreate(o AutoCreatable) {
	db := metadata_db_conn.GetDatabase()
	queries := o.GetCreateQueries()
	for _, query := range queries {
		_, err := db.Exec(query)
		util.PanicOnError(err, "Failed to init entity: "+query)
	}
}
