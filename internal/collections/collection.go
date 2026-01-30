package collections

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/util"
)

type (
	JsonCollection struct {
		Hash        string
		Name        string
		Description string
		Created     time.Time
		Owner       string
		Assets      []string //Asset-Hashes
	}
)

const (
	FilePermissions = 0744
)

// Init creates required directories
func Init() {
	util.CreateDirIfNotExists(config.AssetCollectionsBaseDir, FilePermissions)
}

// AddCollection creates a new collection JSON file
func AddCollection(name string, description string, owner string, assetHashes []string) (*JsonCollection, error) {

	collection := CreateNew(
		name,
		description,
		owner,
		assetHashes)

	collectionDataFile := GetCollectionFilePath(collection.Hash)

	return collection, collection.Save(collectionDataFile)
}

// CreateNew generates collection-hash and creates JsonCollection
func CreateNew(name string, description string, owner string, assetHashes []string) *JsonCollection {

	hashCreator := sha256.New()
	hashCreator.Write([]byte(owner + "/" + name + "/" + description))

	for _, assetHash := range assetHashes {
		hashCreator.Write([]byte("/"))
		hashCreator.Write([]byte(assetHash))
	}

	hash := fmt.Sprintf("%x", hashCreator.Sum(nil))

	collection := &JsonCollection{
		Hash:        hash,
		Name:        name,
		Description: description,
		Created:     time.Now(),
		Owner:       owner,
		Assets:      assetHashes,
	}

	return collection
}

// Save Create dir if not exists and save JSON
func (collection *JsonCollection) Save(path string) error {

	util.PanicOnError(os.MkdirAll(filepath.Dir(path), FilePermissions), "Failed to create destination directory")

	jsonBytes, err := json.Marshal(collection)
	if err != nil {
		return err
	}

	return os.WriteFile(path, jsonBytes, FilePermissions)
}

// LoadIfExists Load JSON-file, if exists.
func LoadIfExists(path string) (*JsonCollection, error) {

	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var collection = &JsonCollection{}
	err = json.Unmarshal(buf, collection)
	if err != nil {
		return nil, err
	}

	return collection, err
}

// LoadByHash returns JsonAssetMetaData loaded from JSON-file
func LoadByHash(collectionHash string) (*JsonCollection, error) {
	path := GetCollectionFilePath(collectionHash)
	meta, err := LoadIfExists(path)
	return meta, err
}

// GetCollectionFilePath returns the path and filename of a collection file.
func GetCollectionFilePath(assetHash string) string {
	name := fmt.Sprintf("%s%s", assetHash[2:], ".json")
	path := filepath.Join(
		config.AssetCollectionsBaseDir,
		assetHash[:2],
		name)
	return path
}
