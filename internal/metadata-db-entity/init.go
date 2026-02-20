package metadata_db_entity

func AutoCreateEntities() {

	var autoCreateables = []AutoCreatable{
		&MimeType{},
		&FileName{},
		&PathItem{},
		&Owner{},
		&Asset{},
		&Origin{},
		&Collection{},
		&FaceSimilarity{},
	}
	for _, autoCreateable := range autoCreateables {
		AutoCreate(autoCreateable)
	}

}
