package metadata_db

import (
	"errors"
	"time"

	"github.com/c8121/asset-storage/internal/collections"
	"github.com/c8121/asset-storage/internal/metadata_db/entity"
	"github.com/c8121/go_db_entities"
)

type CollectionListItem struct {
	Id      int64
	UUID    string
	Name    string
	Created time.Time
}

type CollectionListFilter struct {
	Offset int
	Count  int
}

// AddCollection adds/updates collection-data in database
func AddCollection(jsonCollection *collections.JsonAssetCollection) error {

	collection := &metadata_db_entity.Collection{}
	err := go_db_entities.Load(metadata_db_entity.DbContext, collection, *go_db_entities.NewFilter("UUID", go_db_entities.CompareEqual, jsonCollection.UUID))
	if err != nil && !errors.Is(go_db_entities.ErrNotFound, err) {
		return err
	}

	collection.UUID = jsonCollection.UUID
	collection.Name = jsonCollection.Name
	if collection.Created.IsZero() {
		collection.Created = jsonCollection.Created
	}

	err = go_db_entities.Save(metadata_db_entity.DbContext, collection)
	if err != nil {
		return err
	}

	return nil
}

// RemoveCollection removes collection-data in database
func RemoveCollection(jsonCollection *collections.JsonAssetCollection) error {

	collection := &metadata_db_entity.Collection{}
	err := go_db_entities.Load(metadata_db_entity.DbContext, collection, *go_db_entities.NewFilter("UUID", go_db_entities.CompareEqual, jsonCollection.UUID))
	if err != nil {
		return err
	}

	err = go_db_entities.Delete(metadata_db_entity.DbContext, collection)
	if err != nil {
		return err
	}

	return nil
}

func ListCollections(filter *CollectionListFilter) ([]CollectionListItem, error) {

	f := go_db_entities.Filter{
		Offset:  filter.Offset,
		Count:   filter.Count,
		OrderBy: []string{"Created DESC"},
	}

	mimeTypeSupplier := func() go_db_entities.Entity {
		return &metadata_db_entity.Collection{}
	}

	items, err := go_db_entities.Select(metadata_db_entity.DbContext, mimeTypeSupplier, f)
	if err != nil && !errors.Is(go_db_entities.ErrNotFound, err) {
		return nil, err
	}

	collections := make([]CollectionListItem, len(items))
	for i, item := range items {
		collection := item.(*metadata_db_entity.Collection)
		listItem := &CollectionListItem{
			Id:      collection.Id,
			UUID:    collection.UUID,
			Name:    collection.Name,
			Created: collection.Created,
		}
		collections[i] = *listItem
	}

	return collections, nil
}
