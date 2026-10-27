package metadata_db_entity

import (
	"time"
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

	/*stmt, err := prv.Prepare("SELECT id, asset, name, path, owner, fileTime FROM origin WHERE asset = ?;")
	if err != nil {
		return nil, err
	}
	defer util.CloseOrLog(stmt)

	list := &[]Origin{}

	if rows, err := stmt.Query(assetId); err == nil {
		defer util.CloseOrLog(rows)
		if rows.Next() {
			var origin Origin
			if err := origin.Scan(rows); err != nil {
				fmt.Printf("Error scanning rows: %s\n", err)
				return nil, err
			}
			*list = append(*list, origin)
		}

	} else {
		return nil, err
	}

	return list, nil*/
	panic("implement me")
}

func RemoveOrigins(asset *Asset) error {

	/*stmt, err := prv.Prepare("DELETE FROM origin WHERE asset = ?;")
	if err != nil {
		return err
	}
	defer util.CloseOrLog(stmt)

	_, err = stmt.Exec(asset.Id)
	return err*/
	panic("implement me")
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

	/*stmt, err := prv.Prepare("DELETE FROM origin WHERE asset = ? AND path = ?;")
	if err != nil {
		return err
	}
	defer util.CloseOrLog(stmt)

	_, err = stmt.Exec(assetId, pathId)
	return err*/
	panic("implement me")
}
