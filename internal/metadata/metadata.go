package metadata

import (
	"encoding/json"
	"errors"
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

// Init creates required directories
func Init() {
	util.CreateDirIfNotExists(config.AssetMetaDataBaseDir, FilePermissions)
}

// AddMetaData creates or updates meta-data JSON file
func AddMetaData(hash string, mimeType string, name string, path string, owner string, fileTime time.Time) (*JsonAssetMetaData, error) {

	metaDataFile := GetMetaDataFilePath(hash)

	metaData, err := LoadIfExists(metaDataFile)
	if errors.Is(err, os.ErrNotExist) {
		metaData = CreateNew(
			hash,
			mimeType,
			name,
			path,
			owner,
			fileTime)
	} else if err != nil {
		return nil, err
	} else {
		metaData.AddOrigin(
			name,
			path,
			owner,
			fileTime)
	}

	//fmt.Printf("MetaData: %s\n", metaData)
	return metaData, metaData.Save(metaDataFile)

}

// CreateNew Create new JsonAssetMetaData struct, filled with given data
func CreateNew(hash string, mimeType string, name string, path string, owner string, fileTime time.Time) *JsonAssetMetaData {
	assetMetadata := &JsonAssetMetaData{
		Hash:     hash,
		MimeType: mimeType,
		Origins: []JsonAssetOrigin{
			{
				Name:     name,
				Path:     path,
				Owner:    owner,
				FileTime: fileTime,
			},
		},
	}

	return assetMetadata
}

// AddOrigin Add origin data if not exists
func (assetMetaData *JsonAssetMetaData) AddOrigin(name string, path string, owner string, time time.Time) {

	for _, origin := range assetMetaData.Origins {
		if origin.Name == name &&
			origin.Path == path &&
			origin.Owner == owner &&
			origin.FileTime == time {
			return
		}
	}

	assetMetaData.Origins = append(assetMetaData.Origins, JsonAssetOrigin{
		Name:     name,
		Path:     path,
		Owner:    owner,
		FileTime: time,
	})
}

// GetLatestOrigin finds the newest origin within given meta-data
func GetLatestOrigin(assetMetaData *JsonAssetMetaData) *JsonAssetOrigin {
	var latest *JsonAssetOrigin = nil
	for _, origin := range assetMetaData.Origins {
		if latest == nil || latest.FileTime.Before(origin.FileTime) {
			latest = &origin
		}
	}
	return latest
}

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
		assetMetadata.Hash = storage.HashFromPath(path)
	}
	return assetMetadata, err
}

// LoadByHash returns JsonAssetMetaData loaded from JSON-file
func LoadByHash(assetHash string) (*JsonAssetMetaData, error) {
	path := GetMetaDataFilePath(assetHash)
	meta, err := LoadIfExists(path)
	return meta, err
}

// GetMetaDataFilePath returns the path and filename of a meta-data file.
func GetMetaDataFilePath(assetHash string) string {
	name := fmt.Sprintf("%s%s", assetHash[2:], ".json")
	path := filepath.Join(
		config.AssetMetaDataBaseDir,
		assetHash[:2],
		name)
	return path
}
