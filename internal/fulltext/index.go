package fulltext

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/metadata"
	"github.com/c8121/asset-storage/internal/util"
)

type (
	Document struct {
		Name string
		Text string
	}
)

var (
	index bleve.Index  = nil
	batch *bleve.Batch = nil
)

const (
	FilePermissions = 0744
	IndexName       = "index.bleve"

	batchSize = 100
)

func Open() error {

	if index != nil {
		fmt.Println("Index already opened")
		return nil
	}

	indexPath := filepath.Join(config.AssetFulltextIndexBaseDir, IndexName)
	var err error

	index, err = bleve.Open(indexPath)
	if errors.Is(err, bleve.ErrorIndexPathDoesNotExist) {
		fmt.Printf("Creating new index: %s\n", indexPath)

		indexMapping, err := buildIndexMapping()
		if err != nil {
			return err
		}
		index, err = bleve.New(indexPath, indexMapping)
		if err != nil {
			return err
		}

	} else if err != nil {
		return err
	}

	fmt.Printf("Opening index: %s\n", indexPath)

	return nil
}

func Close() {
	if index != nil {
		if batch != nil {
			fmt.Printf("Close index, add batch of %d.\n", batch.Size())
			util.LogError(index.Batch(batch))
			batch = nil
		}

		util.LogError(index.Close())
		index = nil
	}
}

func Add(assetMeta *metadata.JsonAssetMetaData, text string) error {
	if index == nil {
		return fmt.Errorf("index not opened (call fulltext.Open() before)")
	}

	//fmt.Printf("Adding from %s: %s\n", assetHash, text)

	if batch == nil {
		batch = index.NewBatch()
	} else if batch.Size() >= batchSize {
		fmt.Printf("Add batch of %d.\n", batch.Size())
		err := index.Batch(batch)
		if err != nil {
			return err
		}
		batch = index.NewBatch()
	}

	name := assetMeta.Hash //Fallback
	latestOrigin := metadata.GetLatestOrigin(assetMeta)
	if latestOrigin != nil {
		name = latestOrigin.Name
	}

	doc := Document{
		Name: name,
		Text: text,
	}

	return batch.Index(assetMeta.Hash, doc)
}

func Find(query string) ([]string, error) {
	if index == nil {
		return nil, fmt.Errorf("index not opened (call fulltext.Open() before)")
	}

	matchQuery := bleve.NewQueryStringQuery(query)
	//matchQuery := bleve.NewMatchAllQuery()
	searchRequest := bleve.NewSearchRequestOptions(matchQuery, 30, 0, false)
	searchResults, err := index.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	var hashes = make([]string, 0)

	fmt.Printf("Found %d documents\n", searchResults.Total)
	for _, hit := range searchResults.Hits {
		hashes = append(hashes, hit.ID)
	}

	return hashes, nil
}

func buildIndexMapping() (mapping.IndexMapping, error) {

	indexMapping := bleve.NewIndexMapping()
	docMapping := bleve.NewDocumentMapping()

	textMappingNoStore := bleve.NewTextFieldMapping()
	textMappingNoStore.Store = false
	textMappingNoStore.Index = true
	textMappingNoStore.IncludeInAll = true
	docMapping.AddFieldMappingsAt("Text", textMappingNoStore)

	indexMapping.DefaultMapping = docMapping

	return indexMapping, nil
}
