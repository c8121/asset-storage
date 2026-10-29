package metadata_db_entity

import (
	"errors"

	"github.com/c8121/go_db_entities"
)

type Owner struct {
	Id   int64
	Name string `dbe:"Size: 255, Index"`
}

func (o *Owner) GetEntityName() string {
	return "Owner"
}

func (o *Owner) GetId() int64 {
	return o.Id
}

func (o *Owner) SetId(id int64) {
	o.Id = id
}

var (
	OwnerCache map[string]*Owner
)

func init() {
	OwnerCache = make(map[string]*Owner)
}

func LoadOwner(name string, createIfNotExists bool) (*Owner, error) {

	owner, ok := OwnerCache[name]
	if ok {
		return owner, nil
	}

	owner = &Owner{}
	err := go_db_entities.Load(DbContext, owner, *go_db_entities.NewFilter("Name", go_db_entities.CompareEqual, name))
	if !errors.Is(err, go_db_entities.ErrNotFound) && err != nil {
		return nil, err
	}

	OwnerCache[name] = owner
	return owner, err
}
