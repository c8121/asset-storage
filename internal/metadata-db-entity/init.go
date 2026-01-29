package metadata_db_entity

func AutoCreateEntities() {

	var autoCreateables = []AutoCreatable{
		&MimeType{},
		&FileName{},
		&PathItem{},
		&Owner{},
		&Asset{},
		&Origin{},
	}
	for _, autoCreateable := range autoCreateables {
		AutoCreate(autoCreateable)
	}

}
