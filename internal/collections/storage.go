package collections

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/c8121/asset-storage/internal/config"
)

type (
	JsonAssetCollection struct {
		UUID        string
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

// Save Create dir if not exists and save JSON
func (collection *JsonAssetCollection) Save(path string) error {

	err := os.MkdirAll(filepath.Dir(path), FilePermissions)
	if err != nil {
		fmt.Printf("Error creating directory %s: %v\n", path, err)
		return err
	}

	jsonBytes, err := json.Marshal(collection)
	if err != nil {
		return err
	}

	return os.WriteFile(path, jsonBytes, FilePermissions)
}

// LoadIfExists Load JSON-file, if exists.
func LoadIfExists(path string) (*JsonAssetCollection, error) {

	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var collection = &JsonAssetCollection{}
	err = json.Unmarshal(buf, collection)
	if err != nil {
		return nil, err
	}

	return collection, err
}

// LoadById returns JsonAssetMetaData loaded from JSON-file
func LoadById(uuid string) (*JsonAssetCollection, error) {
	path := GetCollectionFilePath(uuid)
	meta, err := LoadIfExists(path)
	return meta, err
}

// GetCollectionFilePath returns the path and filename of a collection file.
func GetCollectionFilePath(uuid string) string {
	name := fmt.Sprintf("%s%s", uuid[2:], ".json")
	path := filepath.Join(
		config.AssetCollectionsBaseDir,
		uuid[:2],
		name)
	return path
}
