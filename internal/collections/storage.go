package collections

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/c8121/asset-storage/shared/config"
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

// Save Create dir if not exists and save JSON
func (collection *JsonAssetCollection) Save(path string) error {

	err := os.MkdirAll(filepath.Dir(path), config.AppConfig.FilePermissions)
	if err != nil {
		fmt.Printf("Error creating directory %s: %v\n", path, err)
		return err
	}

	jsonBytes, err := json.Marshal(collection)
	if err != nil {
		return err
	}

	return os.WriteFile(path, jsonBytes, config.AppConfig.FilePermissions)
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
		config.AppConfig.AssetCollectionsBaseDir,
		uuid[:2],
		name)
	return path
}

// RemoveCollection deletes the collection file
func RemoveCollection(uuid string) error {
	collectionFile := GetCollectionFilePath(uuid)
	if _, err := os.Stat(collectionFile); err != nil {
		return err
	}

	return os.Remove(collectionFile)
}

// Walk visits every collection (file) in storage and calls given handler function on it.
func Walk(handler func(path string)) error {

	children, err := os.ReadDir(config.AppConfig.AssetCollectionsBaseDir)
	if err != nil {
		return err
	}

	for _, child := range children {
		dir := filepath.Join(config.AppConfig.AssetCollectionsBaseDir, child.Name())
		files, err := os.ReadDir(dir)
		if err != nil {
			return err
		}

		for _, file := range files {
			filePath := filepath.Join(dir, file.Name())
			handler(filePath)
		}
	}
	return nil
}
