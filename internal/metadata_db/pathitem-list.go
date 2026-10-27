package metadata_db

import (
	"github.com/c8121/asset-storage/internal/metadata_db_entity"
	"github.com/c8121/asset-storage/internal/util"
)

func FindPathItem(path string) (*metadata_db_entity.PathItem, error) {

	names := util.SplitPath(path)

	var item *metadata_db_entity.PathItem = nil
	var parentId = 0

	for _, name := range names {

		var err error
		item, err = metadata_db_entity.GetPathItemByName(name, parentId)
		if err != nil {
			return nil, err
		}
		if item == nil {
			return nil, nil
		}
		parentId = int(item.Id)
	}

	return item, nil
}
