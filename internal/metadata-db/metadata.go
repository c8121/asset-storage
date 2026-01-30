package metadata_db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/c8121/asset-storage/internal/metadata"
	metadata_db_entity "github.com/c8121/asset-storage/internal/metadata-db-entity"
	"github.com/c8121/asset-storage/internal/util"
)

// AddMetaData adds/updates meta-data in database
func AddMetaData(jsonMeta *metadata.JsonAssetMetaData) error {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer util.RollbackOrLog(tx)

	err = AddMetaDataTx(tx, jsonMeta)
	if err != nil {
		return err
	}

	return util.CommitOrLog(tx)
}

// AddMetaDataTx adds/updates meta-data in database
func AddMetaDataTx(tx *sql.Tx, jsonMeta *metadata.JsonAssetMetaData) error {

	var asset = &metadata_db_entity.Asset{Hash: jsonMeta.Hash}
	err := metadata_db_entity.LoadTx(tx, asset)
	if !errors.Is(err, metadata_db_entity.ErrNotFound) && err != nil {
		return err
	}

	mimeType, err := metadata_db_entity.GetMimeTypeTx(tx, jsonMeta.MimeType, true)
	if err != nil {
		return err
	}

	asset.MimeType = mimeType.Id

	latestOrigin := metadata.GetLatestOrigin(jsonMeta)
	if latestOrigin != nil {
		asset.FileTime = latestOrigin.FileTime
		asset.Name = metadata_db_entity.GetFileNameIdTx(tx, latestOrigin.Name, true)
	}

	err = metadata_db_entity.SaveTx(tx, asset)
	if err != nil {
		return err
	}

	err = removeOriginsTx(tx, asset)
	if err != nil {
		return err
	}

	for _, jsonOrigin := range jsonMeta.Origins {

		var origin = &metadata_db_entity.Origin{
			Asset:    asset.Id,
			Name:     metadata_db_entity.GetFileNameIdTx(tx, jsonOrigin.Name, true),
			Path:     metadata_db_entity.GetPathItemIdTx(tx, jsonOrigin.Path, true),
			Owner:    metadata_db_entity.GetOwnerIdTx(tx, jsonOrigin.Owner, true),
			FileTime: jsonOrigin.FileTime,
		}
		err = metadata_db_entity.SaveTx(tx, origin)
		if err != nil {
			return err
		}
	}

	return nil
}

func removeOriginsTx(tx *sql.Tx, asset *metadata_db_entity.Asset) error {

	stmt, err := tx.Prepare("DELETE FROM origin WHERE asset = ?;")
	if err != nil {
		return err
	}
	defer util.CloseOrLog(stmt)

	_, err = stmt.Exec(asset.Id)
	return err
}
