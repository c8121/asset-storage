package metadata_db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/c8121/asset-storage/internal/metadata"
)

func AddMetaData(jsonMeta *metadata.JsonAssetMetaData) error {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer rollbackOrLog(tx)

	err = AddMetaDataTx(tx, jsonMeta)
	if err != nil {
		return err
	}

	return commitOrLog(tx)
}

func AddMetaDataTx(tx *sql.Tx, jsonMeta *metadata.JsonAssetMetaData) error {

	var asset = &Asset{Hash: jsonMeta.Hash}
	err := LoadTx(tx, asset)
	if !errors.Is(err, ErrNotFound) && err != nil {
		return err
	}

	mimeType, err := GetMimeTypeTx(tx, jsonMeta.MimeType, true)
	if err != nil {
		return err
	}

	asset.MimeType = mimeType.Id

	latestOrigin := metadata.GetLatestOrigin(jsonMeta)
	if latestOrigin != nil {
		asset.FileTime = latestOrigin.FileTime
	}

	err = SaveTx(tx, asset)
	if err != nil {
		return err
	}

	for _, jsonOrigin := range jsonMeta.Origins {

		var origin = &Origin{
			Asset:    asset.Id,
			Name:     jsonOrigin.Name,
			Path:     GetPathItemIdTx(tx, jsonOrigin.Path, true),
			Owner:    GetOwnerIdTx(tx, jsonOrigin.Owner, true),
			FileTime: jsonOrigin.FileTime,
		}
		err = SaveTx(tx, origin)
		if err != nil {
			return err
		}
	}

	return nil
}
