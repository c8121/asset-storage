package metadata

import (
	"errors"
	"os"
	"time"
)

// AddMetaData creates or updates metadata JSON file
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
