package metadata_db_entity

import (
	"fmt"

	"github.com/c8121/go_db_entities"
)

func CheckCreateTables() {

	entities := []go_db_entities.Entity{
		&Asset{},
		&Collection{},
		&FaceEmbedding{},
		&FileName{},
		&MimeType{},
		&Origin{},
		&Owner{},
		&PathItem{},
	}

	for _, entity := range entities {
		err := go_db_entities.Create(DbContext, entity)
		if err != nil {
			fmt.Println(err)
		}
	}
}
