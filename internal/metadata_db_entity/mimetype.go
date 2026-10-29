package metadata_db_entity

import (
	"errors"
	"strings"

	"github.com/c8121/go_db_entities"
)

type MimeType struct {
	Id   int64
	Name string `dbe:"Size: 255, Index"`
}

func (m *MimeType) GetEntityName() string {
	return "MimeType"
}

func (m *MimeType) GetId() int64 {
	return m.Id
}

func (m *MimeType) SetId(id int64) {
	m.Id = id
}

var (
	mimeTypeCache map[string]*MimeType
)

func init() {
	mimeTypeCache = make(map[string]*MimeType)
}

func NormalizeName(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	if p := strings.Index(name, ";"); p != -1 {
		name = strings.TrimSpace(name[:p])
	}
	return name
}

func ListMimeTypes() ([]MimeType, error) {

	filter := &go_db_entities.Filter{}
	filter.OrderBy = []string{"Name", "Id"}

	mimeTypeSupplier := func() go_db_entities.Entity {
		return &MimeType{}
	}

	items, err := go_db_entities.Select(DbContext, mimeTypeSupplier, *filter)
	if err != nil {
		return nil, err
	}

	mimeTypes := make([]MimeType, len(items))
	for i, item := range items {
		mimeType := item.(*MimeType)
		mimeTypes[i] = *mimeType
	}

	return mimeTypes, nil
}

func LoadMimeType(name string, createIfNotExists bool) (*MimeType, error) {

	name = NormalizeName(name)
	mimeType, ok := mimeTypeCache[name]
	if ok {
		return mimeType, nil
	}

	mimeType = &MimeType{}
	err := go_db_entities.Load(DbContext, mimeType, go_db_entities.Filter{
		FieldName: "Name",
		Value:     name,
	})
	if errors.Is(err, go_db_entities.ErrNotFound) && createIfNotExists {
		mimeType.Name = name
		if err = go_db_entities.Save(DbContext, mimeType); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	mimeTypeCache[name] = mimeType
	return mimeType, nil
}

func LoadMimeTypeById(id int64) (*MimeType, error) {

	mimeType := &MimeType{}
	err := go_db_entities.Load(DbContext, mimeType, *go_db_entities.NewIdFilter(id))
	if err != nil {
		return nil, err
	}

	return mimeType, nil
}
