package metadata_db_entity

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/c8121/asset-storage/internal/util"
	"github.com/c8121/go_db_entities"
)

type PathItem struct {
	Id     int64
	Parent int64  `dbe:"Index"`
	Name   string `dbe:"Size: 1024, Index"`
}

func (p *PathItem) GetEntityName() string {
	return "PathItem"
}

func (p *PathItem) GetId() int64 {
	return p.Id
}

func (p *PathItem) SetId(id int64) {
	p.Id = id
}

var (
	pathItemCache     map[string]*PathItem
	loadPathChunkSize = 100
)

func init() {
	pathItemCache = make(map[string]*PathItem)
}

// SplitPath split at \ and /, remove 'file:'
func SplitPath(path string) []string {

	l := len(path)
	if l == 0 {
		return []string{""}
	}

	names := make([]string, 0)

	s := 0
	i := 0
	if strings.HasPrefix(strings.ToLower(path), "file:") {
		s = 5
		i = 5
	}

	for ; i < l; i++ {
		c := path[i]
		if c == '/' || c == '\\' {
			if i-s > 0 {
				names = append(names, path[s:i])
			}
			s = i + 1
		}
	}
	if i > s {
		names = append(names, path[s:])
	}

	if len(names) == 0 {
		return []string{""}
	}

	return names
}

// LoadPathItem gets PathItem from db, splits path and searches
func LoadPathItem(path string, createIfNotExists bool) (*PathItem, error) {

	names := SplitPath(path)
	var parent int64 = 0

	var pathItem *PathItem
	for _, name := range names {

		cacheKey := fmt.Sprintf("%d/%s", parent, name)
		cachedItem, ok := pathItemCache[cacheKey]
		if !ok {
			var err error
			pathItem, err = GetPathItemByName(name, parent)
			if errors.Is(err, go_db_entities.ErrNotFound) && createIfNotExists {
				pathItem = &PathItem{Parent: parent, Name: name}
				if err = go_db_entities.Save(DbContext, pathItem); err != nil {
					return nil, err
				}
			} else if err != nil {
				return nil, err
			}
			pathItemCache[cacheKey] = pathItem
		} else {
			pathItem = cachedItem
			//fmt.Printf("Cached: %s = %v\n", cacheKey, pathItem)
		}

		parent = pathItem.Id
	}

	return pathItem, nil
}

func ListPathItems(parentId int64) ([]PathItem, error) {

	filter := go_db_entities.NewFilter("Parent", go_db_entities.CompareEqual, parentId)
	filter.OrderBy = []string{"Name"}

	pathItemSupplier := func() go_db_entities.Entity {
		return &PathItem{}
	}

	items, err := go_db_entities.Select(DbContext, pathItemSupplier, *filter)
	if err != nil && !errors.Is(go_db_entities.ErrNotFound, err) {
		return nil, err
	}

	pathItems := make([]PathItem, len(items))
	for i, item := range items {
		pathItem := item.(*PathItem)
		pathItems[i] = *pathItem
	}

	return pathItems, nil
}

func GetPathItemByName(name string, parentId int64) (*PathItem, error) {

	filter := go_db_entities.NewFilter("Name", go_db_entities.CompareEqual, name)
	filter.Next = go_db_entities.NewFilter("Parent", go_db_entities.CompareEqual, parentId)

	pathItem := &PathItem{}
	err := go_db_entities.Load(DbContext, pathItem, *filter)
	if err != nil {
		return nil, err
	}

	return pathItem, nil
}

func FindPathItem(path string) (*PathItem, error) {

	names := util.SplitPath(path)

	var item *PathItem = nil
	var parentId int64 = 0

	for _, name := range names {

		var err error
		item, err = GetPathItemByName(name, parentId)
		if err != nil {
			return nil, err
		}
		if item == nil {
			return nil, nil
		}
		parentId = item.Id
	}

	return item, nil
}

// FindAssetPathIds looks for origin-paths in database and
// returns it as map: asset.Id -> [path.Id,...]
func FindAssetPathIds(assetIds []int64) (map[int64][]int64, error) {

	var idMap = make(map[int64][]int64, len(assetIds))
	var chunk = make([]int64, 0)

	var query = go_db_entities.TableQuery{
		FieldNames: []string{"Path"},
		Entity:     &Origin{},
	}
	query.AddJoin(go_db_entities.TableQueryJoin{
		TableQuery: &go_db_entities.TableQuery{
			FieldNames: []string{"Id"},
			Entity:     &Asset{},
		},
		OnFrom: "Asset",
		OnTo:   "Id",
	})

	for i, assetId := range assetIds {

		chunk = append(chunk, assetId)
		if len(chunk) == loadPathChunkSize || (i+1) == len(assetIds) {

			query.Filter = &go_db_entities.Filter{
				FieldName: "Asset.Id",
				Value:     chunk,
			}

			rows, err := go_db_entities.Query(DbContext, query)
			if err != nil {
				return nil, err
			}

			for _, row := range rows {
				pathId := row[0].(int64)
				assetId := row[1].(int64)

				list, exists := idMap[assetId]
				if exists {
					idMap[assetId] = append(list, pathId)
				} else {
					list = make([]int64, 0)
					idMap[assetId] = append(list, pathId)
				}
			}

			chunk = make([]int64, 0)
		}
	}

	return idMap, nil
}

func GetAssetIdsFromPath(pathId int64) ([]int64, error) {

	query := go_db_entities.TableQuery{
		Entity:     &Origin{},
		FieldNames: []string{"Asset"},
		Filter: &go_db_entities.Filter{
			FieldName: "Path",
			Value:     pathId,
			Distinct:  true,
		},
	}

	rows, err := go_db_entities.Query(DbContext, query)
	if err != nil && !errors.Is(go_db_entities.ErrNotFound, err) {
		return nil, err
	}

	list := make([]int64, len(rows))
	for i, row := range rows {
		list[i] = row[0].(int64)
	}

	return list, err

}

func RemovePathIfEmpty(pathId int64) error {

	cnt := countPathRefs(go_db_entities.TableQuery{
		FieldNames: []string{"Id"},
		Entity:     &PathItem{},
		Filter:     &go_db_entities.Filter{FieldName: "Parent", Value: pathId},
	})
	if cnt > 0 {
		fmt.Printf("Will not remove path %d, it has %d children\n", pathId, cnt)
		return nil
	}

	cnt = countPathRefs(go_db_entities.TableQuery{
		FieldNames: []string{"Id"},
		Entity:     &Origin{},
		Filter:     &go_db_entities.Filter{FieldName: "Path", Value: pathId},
	})
	if cnt > 0 {
		fmt.Printf("Will not remove path %d, it is referenced in %d origins\n", pathId, cnt)
		return nil
	}

	fmt.Printf("Removing path %d\n", pathId)
	err := go_db_entities.DeleteBatch(DbContext, &PathItem{}, *go_db_entities.NewIdFilter(pathId))
	if err != nil {
		return err
	}

	return nil
}

func countPathRefs(query go_db_entities.TableQuery) int {

	rows, err := go_db_entities.Query(DbContext, query)
	if err != nil && !errors.Is(go_db_entities.ErrNotFound, err) {
		log.Println(err)
		return 9999
	}

	return len(rows)
}
