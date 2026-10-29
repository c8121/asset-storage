package metadata_db_entity

import (
	"errors"
	"fmt"
	"time"

	"github.com/c8121/asset-storage/internal/metadata"
	"github.com/c8121/asset-storage/internal/util"
	"github.com/c8121/go_db_entities"
)

type Asset struct {
	Id       int64
	Hash     string    `dbe:"Size: 128, Index"`
	MimeType int64     `dbe:"Index"`
	FileTime time.Time `dbe:"Index"` //Max of all origins
	Name     int64     `dbe:"Index"` //Latest name
}

func (p *Asset) GetEntityName() string {
	return "Asset"
}

func (a *Asset) GetId() int64 {
	return a.Id
}

func (a *Asset) SetId(id int64) {
	a.Id = id
}

func LoadAsset(hash string) (*Asset, error) {
	var asset = &Asset{Hash: hash}
	err := go_db_entities.Load(DbContext, asset, *go_db_entities.NewFilter("Hash", go_db_entities.CompareEqual, hash))
	if !errors.Is(err, go_db_entities.ErrNotFound) && err != nil {
		return nil, err
	}
	return asset, err
}

func LoadAssetById(id int64) (*Asset, error) {
	var asset = &Asset{}
	err := go_db_entities.Load(DbContext, asset, *go_db_entities.NewIdFilter(id))
	if !errors.Is(err, go_db_entities.ErrNotFound) && err != nil {
		return nil, err
	}
	return asset, err
}

// AddAsset adds/updates metadata in database
func AddAsset(jsonMeta *metadata.JsonAssetMetaData) error {

	asset, err := LoadAsset(jsonMeta.Hash)
	if err != nil && !errors.Is(err, go_db_entities.ErrNotFound) {
		return err
	}

	mimeType, err := LoadMimeType(jsonMeta.MimeType, true)
	if err != nil {
		return err
	}

	asset.MimeType = mimeType.Id

	latestOrigin := metadata.GetLatestOrigin(jsonMeta)
	if latestOrigin != nil {
		asset.FileTime = latestOrigin.FileTime
		asset.Name = getFileNameId(latestOrigin.Name)
	}

	err = go_db_entities.Save(DbContext, asset)
	if err != nil {
		return err
	}

	err = RemoveOrigins(asset)
	if err != nil {
		return err
	}

	for _, jsonOrigin := range jsonMeta.Origins {

		var origin = &Origin{
			Asset:    asset.Id,
			Name:     getFileNameId(jsonOrigin.Name),
			Path:     getPathItemId(jsonOrigin.Path),
			Owner:    getOwnerId(jsonOrigin.Owner),
			FileTime: jsonOrigin.FileTime,
		}
		err = go_db_entities.Save(DbContext, origin)
		if err != nil {
			return err
		}
	}

	return nil
}

func getFileNameId(name string) int64 {
	fileName, err := LoadFileName(name, true)
	if err != nil {
		util.LogError(err)
		return 0
	}
	if fileName == nil {
		util.LogError(errors.New("file name not found"))
		return 0
	}
	return fileName.Id
}

func getPathItemId(path string) int64 {
	pathItem, err := LoadPathItem(path, true)
	if err != nil {
		util.LogError(err)
		return 0
	}
	if pathItem == nil {
		util.LogError(errors.New("path not found"))
		return 0
	}
	return pathItem.Id
}

func getOwnerId(name string) int64 {
	owner, err := LoadOwner(name, true)
	if err != nil {
		util.LogError(err)
		return 0
	}
	if owner == nil {
		util.LogError(errors.New("path not found"))
		return 0
	}
	return owner.Id
}

func RemoveAsset(assetId int64, pathId int64) (int, error) {

	if pathId > 0 {
		if err := RemoveOriginsByAssetIdAndPathId(assetId, pathId); err != nil {
			return 9999, err
		}
	} else {
		if err := RemoveOriginsByAssetId(assetId); err != nil {
			return 9999, err
		}
	}

	remainingOrigins, err := LoadOriginsByAssetId(assetId)
	if err != nil {
		return 9999, err
	}
	if len(*remainingOrigins) > 0 {
		fmt.Printf("Keep asset, it has more origins: %d\n", assetId)
		return len(*remainingOrigins), nil
	}

	asset := &Asset{Id: assetId}
	err = go_db_entities.Delete(DbContext, asset)
	if err != nil {
		return 9999, err
	}

	return 0, nil
}
