package metadata_db

import (
	"time"

	"github.com/c8121/asset-storage/internal/collections"
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

	/*var collection = &metadata_db_entity.Collection{UUID: jsonCollection.UUID}
	err := db_entity.Load(tx, collection)
	if !errors.Is(err, db_entity.ErrNotFound) && err != nil {
		return err
	}

	collection.Name = jsonCollection.Name
	if collection.Created.IsZero() {
		collection.Created = jsonCollection.Created
	}

	err = db_entity.SaveEntity(tx, collection)
	if err != nil {
		return err
	}*/

	return nil
}

// RemoveCollection removes collection-data in database
func RemoveCollection(jsonCollection *collections.JsonAssetCollection) error {

	/*var collection = &metadata_db_entity.Collection{UUID: jsonCollection.UUID}
	err := db_entity.Load(tx, collection)
	if !errors.Is(err, db_entity.ErrNotFound) && err != nil {
		return err
	}

	err = db_entity.DeleteEntity(tx, collection)
	if err != nil {
		return err
	}*/

	return nil
}

func ListCollections(filter *CollectionListFilter) ([]CollectionListItem, error) {

	/*var query = "SELECT id, uuid, name, created" +
		" FROM collection" +
		" ORDER BY created DESC" +
		" LIMIT ? OFFSET ?;"

	var params = make([]any, 0)
	params = append(params, filter.Count)
	params = append(params, filter.Offset)

	return loadCollectionList(query, params...)*/
	return nil, nil
}
