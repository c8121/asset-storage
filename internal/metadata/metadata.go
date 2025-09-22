package metadata

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

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

// BaseDir Directory for all meta-data of assets.
func BaseDir() string {
	return "/tmp/asset-metadata"
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

	util.Check(os.MkdirAll(filepath.Dir(path), FilePermissions), "Failed to create destination directory")

	jsonBytes, err := json.Marshal(assetMetaData)
	if err != nil {
		return err
	}

	return os.WriteFile(path, jsonBytes, FilePermissions)
}

// LoadIfExists Load JSON-file, if exists.
func LoadIfExists(path string) (AssetMetadata, error) {

	var assetMetadata AssetMetadata

	buf, err := os.ReadFile(path)
	if err != nil {
		return assetMetadata, err
	}

	err = json.Unmarshal(buf, &assetMetadata)
	return assetMetadata, err
}

func GetMetaDataFilePath(assetHash string) string {
	name := fmt.Sprintf("%s%s", assetHash[2:], ".json")
	path := filepath.Join(
		BaseDir(),
		assetHash[:2],
		name)
	return path
}
