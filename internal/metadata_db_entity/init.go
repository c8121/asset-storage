package metadata_db_entity

import "github.com/c8121/asset-storage/internal/db_entity"

func AutoCreateEntities() {

	var autoCreatables = []db_entity.AutoCreatable{
		&MimeType{},
		&FileName{},
		&PathItem{},
		&Owner{},
		&Asset{},
		&Origin{},
		&Collection{},
		&FaceEmbedding{},
	}
	for _, autoCreatable := range autoCreatables {
		db_entity.AutoCreate(autoCreatable)
	}

}
