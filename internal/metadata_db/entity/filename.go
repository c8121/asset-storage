package metadata_db_entity

import (
	"errors"

	"github.com/c8121/go_db_entities"
)

type FileName struct {
	Id   int64
	Name string `dbe:"Size: 1024, Index"`
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

	fileName := &FileName{}
	err := go_db_entities.Load(DbContext, fileName, go_db_entities.Filter{
		FieldName: "Name",
		Value:     name,
	})
	if errors.Is(err, go_db_entities.ErrNotFound) && createIfNotExists {
		fileName.Name = name
		if err = go_db_entities.Save(DbContext, fileName); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return fileName, nil
}

func LoadFileNameById(id int64) (*FileName, error) {

	fileName := &FileName{}
	err := go_db_entities.Load(DbContext, fileName, *go_db_entities.NewIdFilter(id))
	if err != nil {
		return nil, err
	}

	return fileName, nil
}
