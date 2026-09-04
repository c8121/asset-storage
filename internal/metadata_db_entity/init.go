package metadata_db_entity

func AutoCreateEntities() {

	var autoCreatables = []AutoCreatable{
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
		AutoCreate(autoCreatable)
	}

}
