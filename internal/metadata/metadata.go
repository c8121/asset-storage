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
	AssetMetadata struct {
		MimeType string
		Origins  []Origin
	}

	Origin struct {
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
func AddMetaData(assetHash string, mimeType string, name string, path string, owner string, fileTime time.Time) (*AssetMetadata, error) {

	metaDataFile := GetMetaDataFilePath(assetHash)
	fmt.Printf("MetaDataFile: %s\n", metaDataFile)

	_, metaData, err := LoadIfExists(metaDataFile)
	if errors.Is(err, os.ErrNotExist) {
		metaData = CreateNew(
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
	return &metaData, metaData.Save(metaDataFile)

}

// CreateNew Create new AssetMetadata struct, filled with given data
func CreateNew(mimeType string, name string, path string, owner string, fileTime time.Time) AssetMetadata {
	assetMetadata := AssetMetadata{
		MimeType: mimeType,
		Origins: []Origin{
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
func (assetMetaData *AssetMetadata) AddOrigin(name string, path string, owner string, time time.Time) {

	for _, origin := range assetMetaData.Origins {
		if origin.Name == name &&
			origin.Path == path &&
			origin.Owner == owner &&
			origin.FileTime == time {
			return
		}
	}

	assetMetaData.Origins = append(assetMetaData.Origins, Origin{
		Name:     name,
		Path:     path,
		Owner:    owner,
		FileTime: time,
	})
}

// Save Create dir if not exists and save JSON
func (assetMetaData *AssetMetadata) Save(path string) error {

	util.PanicOnError(os.MkdirAll(filepath.Dir(path), FilePermissions), "Failed to create destination directory")

	jsonBytes, err := json.Marshal(assetMetaData)
	if err != nil {
		return err
	}

	return os.WriteFile(path, jsonBytes, FilePermissions)
}

// LoadIfExists Load JSON-file, if exists. Returns hash, meta-data, error
func LoadIfExists(path string) (string, AssetMetadata, error) {

	var assetMetadata AssetMetadata

	buf, err := os.ReadFile(path)
	if err != nil {
		return "", assetMetadata, err
	}

	err = json.Unmarshal(buf, &assetMetadata)
	return storage.HashFromPath(path), assetMetadata, err
}

// LoadByHash returns AssetMetadata loaded from JSON-file
func LoadByHash(assetHash string) (AssetMetadata, error) {
	path := GetMetaDataFilePath(assetHash)
	_, meta, err := LoadIfExists(path)
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
