package collections

import (
	"errors"
	"os"
	"time"

	"github.com/google/uuid"
)

// AddCollection adds a new collection or updates a existing.
// Parameter uuid can be empty to create a new collection
func AddCollection(uuid string, name string, description string, owner string, assetHashes []string) (*JsonAssetCollection, error) {

	var collectionFile string
	var collection *JsonAssetCollection
	var err error

	if uuid != "" {
		collectionFile = GetCollectionFilePath(uuid)
		collection, err = LoadIfExists(collectionFile)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
	}

	if collection == nil {
		collection = CreateNew(name, description, owner, assetHashes)
		collectionFile = GetCollectionFilePath(collection.UUID)
	} else {
		collection.Name = name
		collection.Description = description
		collection.Owner = owner
		collection.Assets = assetHashes
	}

	return collection, collection.Save(collectionFile)
}

// CreateNew generates collection-hash and creates JsonAssetCollection
func CreateNew(name string, description string, owner string, assetHashes []string) *JsonAssetCollection {

	collection := &JsonAssetCollection{
		UUID:        uuid.NewString(),
		Name:        name,
		Description: description,
		Created:     time.Now(),
		Owner:       owner,
		Assets:      assetHashes,
	}

	return collection
}
