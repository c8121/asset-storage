package metadata_db_entity

type FileName struct {
	Id   int64
	Name string `dbe:"Size: 1034, Index"`
}

func (f *FileName) GetEntityName() string {
	return "FileName"
}

func (f *FileName) GetId() int64 {
	return f.Id
}

func (f *FileName) SetId(id int64) {
	f.Id = id
}

// LoadFileName gets FileName from db, created it required
func LoadFileName(name string, createIfNotExists bool) (*FileName, error) {

	/*var fileName = &FileName{Name: name}
	err := db_entity.LoadEntity(createIfNotExists, fileName)
	if err != nil {
		return nil, err
	}

	return fileName, nil*/
	panic("implement me")
}

func LoadFileNameById(id int64) (*FileName, error) {

	/*var fileName = &FileName{Id: id}
	err := db_entity.LoadEntity(false, fileName)
	if err != nil {
		return nil, err
	}

	return fileName, nil*/
	panic("implement me")
}
