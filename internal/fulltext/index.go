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
		Type string
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

func Add(assetMeta *metadata.JsonAssetMetaData, text string) error {
	err := initIndex()
	if err != nil {
		return err
	}

	//fmt.Printf("Adding from %s: %s\n", assetHash, text)

	if batch == nil {
		batch = index.NewBatch()
	} else if batch.Size() >= batchSize {
		fmt.Printf("Add batch of %d.\n", batch.Size())
		err = index.Batch(batch)
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
		Type: "document",
	}

	return batch.Index(assetMeta.Hash, doc)
}

func Find(query string) error {
	err := initIndex()
	if err != nil {
		return err
	}

	matchQuery := bleve.NewQueryStringQuery(query)
	searchRequest := bleve.NewSearchRequest(matchQuery)
	searchResults, err := index.Search(searchRequest)
	if err != nil {
		return err
	}

	fmt.Printf("Found %d documents\n", searchResults.Total)
	for _, hit := range searchResults.Hits {
		fmt.Println(hit.ID)
	}

	return nil
}

func initIndex() error {
	if index != nil {
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

func buildIndexMapping() (mapping.IndexMapping, error) {

	// Create a new index mapping
	indexMapping := bleve.NewIndexMapping()

	// Create a field mapping for text where Store is false, but Index is true
	textFieldMapping := bleve.NewTextFieldMapping()
	textFieldMapping.Index = true
	textFieldMapping.Store = false // Do not store the raw text
	textFieldMapping.IncludeInAll = false
	//textFieldMapping.IncludeTermVectors = true // Optional

	// Apply it to a specific field in your document mapping
	docMapping := bleve.NewDocumentMapping()
	docMapping.AddFieldMappingsAt("Text", textFieldMapping)

	indexMapping.AddDocumentMapping("document", docMapping)

	return indexMapping, nil
}

func CloseIndex() {
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
