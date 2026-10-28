package metadata_db_entity

import (
	"time"

	"github.com/c8121/go_db_entities"
)

type Origin struct {
	Id       int64
	Asset    int64     `dbe:"Index"`
	Name     int64     `dbe:"Index"`
	Path     int64     `dbe:"Index"`
	Owner    int64     `dbe:"Index"`
	FileTime time.Time `dbe:"Index"`
}

func (o *Origin) GetEntityName() string {
	return "Origin"
}

func (o *Origin) GetId() int64 {
	return o.Id
}

func (o *Origin) SetId(id int64) {
	o.Id = id
}

func LoadOriginsByAssetId(assetId int64) (*[]Origin, error) {

	originSupplier := func() go_db_entities.Entity {
		return &Origin{}
	}

	items, err := go_db_entities.Select(DbContext, originSupplier, *go_db_entities.NewFilter("Asset", go_db_entities.CompareEqual, assetId))
	if err != nil {
		return nil, err
	}

	origins := make([]Origin, len(items))
	for i, item := range items {
		origin := item.(*Origin)
		origins[i] = *origin
	}

	return &origins, nil
}

func RemoveOrigins(asset *Asset) error {

	origin := &Origin{}
	err := go_db_entities.DeleteBatch(DbContext, origin, go_db_entities.Filter{FieldName: "Asset", Comparison: go_db_entities.CompareEqual, Value: asset.Id})
	if err != nil {
		return err
	}

	return nil
}

func RemoveOriginsByAssetId(assetId int64) error {

	/*stmt, err := prv.Prepare("DELETE FROM origin WHERE asset = ?;")
	if err != nil {
		return err
	}
	defer util.CloseOrLog(stmt)

	_, err = stmt.Exec(assetId)
	return err*/
	panic("implement me")
}

func RemoveOriginsByAssetIdAndPathId(assetId int64, pathId int64) error {

	filter := go_db_entities.Filter{
		FieldName: "Asset",
		Value:     assetId,
		Next: &go_db_entities.Filter{
			FieldName: "Path",
			Value:     pathId,
		},
	}

	origin := &Origin{}
	err := go_db_entities.DeleteBatch(DbContext, origin, filter)
	if err != nil {
		return err
	}
	return nil
}
