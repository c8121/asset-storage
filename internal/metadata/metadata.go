package metadata

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/storage"
	"github.com/c8121/asset-storage/internal/util"
)

type (
	JsonAssetMetaData struct {
		Hash     string
		MimeType string
		Origins  []JsonAssetOrigin
	}

	JsonAssetOrigin struct {
		Name     string
		Path     string
		Owner    string
		FileTime time.Time
	}
)

const (
	FilePermissions = 0744
)

// Save Create dir if not exists and save JSON
func (assetMetaData *JsonAssetMetaData) Save(path string) error {

	util.PanicOnError(os.MkdirAll(filepath.Dir(path), FilePermissions), "Failed to create destination directory")

	jsonBytes, err := json.Marshal(assetMetaData)
	if err != nil {
		return err
	}

	return os.WriteFile(path, jsonBytes, FilePermissions)
}

// LoadIfExists Load JSON-file, if exists.
func LoadIfExists(path string) (*JsonAssetMetaData, error) {

	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var assetMetadata = &JsonAssetMetaData{}
	err = json.Unmarshal(buf, assetMetadata)
	if err != nil {
		return nil, err
	}

	if assetMetadata.Hash == "" {
		assetMetadata.Hash = storage.HashFromStoragePath(path)
	}
	return assetMetadata, err
}

// LoadByHash returns JsonAssetMetaData loaded from JSON-file
func LoadByHash(assetHash string) (*JsonAssetMetaData, error) {
	path := GetMetaDataFilePath(assetHash)
	meta, err := LoadIfExists(path)
	return meta, err
}

// GetMetaDataFilePath returns the path and filename of a metadata file.
func GetMetaDataFilePath(assetHash string) string {
	name := fmt.Sprintf("%s%s", assetHash[2:], ".json")
	path := filepath.Join(
		config.AssetMetaDataBaseDir,
		assetHash[:2],
		name)
	return path
}

// Walk visits every asset (file) in storage and calls given handler function on it.
func Walk(handler func(path string)) {

	children, err := os.ReadDir(config.AssetMetaDataBaseDir)
	util.PanicOnError(err, "Failed to read directory")

	for _, child := range children {
		dir := filepath.Join(config.AssetMetaDataBaseDir, child.Name())
		files, err := os.ReadDir(dir)
		util.PanicOnError(err, "Failed to read directory")

		for _, file := range files {
			filePath := filepath.Join(dir, file.Name())
			handler(filePath)
		}
	}
}

// RemoveMetaData deletes the metadata file
func RemoveMetaData(assetHash string) error {
	metaDataFile := GetMetaDataFilePath(assetHash)
	if _, err := os.Stat(metaDataFile); err != nil {
		return err
	}

	return os.Remove(metaDataFile)
}
