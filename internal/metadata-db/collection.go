package metadata_db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/c8121/asset-storage/internal/collections"
	metadata_db_entity "github.com/c8121/asset-storage/internal/metadata-db-entity"
	"github.com/c8121/asset-storage/internal/util"
)

type CollectionListItem struct {
	Id      int64
	Hash    string
	Name    string
	Created time.Time
}

type CollectionListFilter struct {
	Offset int
	Count  int
}

// AddCollection adds/updates collection-data in database
func AddCollection(jsonCollection *collections.JsonCollection) error {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer util.RollbackOrLog(tx)

	err = AddCollectionTx(tx, jsonCollection)
	if err != nil {
		return err
	}

	return util.CommitOrLog(tx)
}

// AddCollectionTx adds/updates collection-data in database
func AddCollectionTx(tx *sql.Tx, jsonCollection *collections.JsonCollection) error {

	var collection = &metadata_db_entity.Collection{Hash: jsonCollection.Hash}
	err := metadata_db_entity.LoadTx(tx, collection)
	if !errors.Is(err, metadata_db_entity.ErrNotFound) && err != nil {
		return err
	}

	collection.Name = jsonCollection.Name
	if collection.Created.IsZero() {
		collection.Created = jsonCollection.Created
	}

	err = metadata_db_entity.SaveTx(tx, collection)
	if err != nil {
		return err
	}

	return nil
}

func ListCollections(filter *CollectionListFilter) ([]CollectionListItem, error) {

	var query = "SELECT id, hash, name, created" +
		" FROM collection" +
		" ORDER BY created DESC" +
		" LIMIT ? OFFSET ?;"

	var params = make([]any, 0)
	params = append(params, filter.Count)
	params = append(params, filter.Offset)

	return loadCollectionList(query, params...)
}

// loadCollectionList queries the database
func loadCollectionList(query string, params ...any) ([]CollectionListItem, error) {

	fmt.Printf("Query: %s\n", query)
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer util.CloseOrLog(stmt)

	var items []CollectionListItem

	if rows, err := stmt.Query(params...); err == nil {
		defer util.CloseOrLog(rows)
		for rows.Next() {
			var item CollectionListItem
			if err := rows.Scan(&item.Id, &item.Hash, &item.Name, &item.Created); err != nil {
				return nil, err
			}
			items = append(items, item)
		}

	} else {
		return nil, err
	}

	return items, nil
}
