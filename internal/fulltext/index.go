package fulltext

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/v2/analysis/lang/en"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/util"
)

type (
	JsonAssetText struct {
		name string
		text string
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

func AddText(assetHash string, text string) error {
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

	doc := &JsonAssetText{
		name: assetHash,
		text: text,
	}
	return batch.Index(assetHash, doc)
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
	} else {
		fmt.Printf("Opening index: %s\n", indexPath)
	}

	return nil
}

// See https://github.com/blevesearch/beer-search/blob/master/mapping.go
func buildIndexMapping() (mapping.IndexMapping, error) {
	// a generic reusable mapping for english text
	englishTextFieldMapping := bleve.NewTextFieldMapping()
	englishTextFieldMapping.Analyzer = en.AnalyzerName

	// a generic reusable mapping for keyword text
	keywordFieldMapping := bleve.NewTextFieldMapping()
	keywordFieldMapping.Analyzer = keyword.Name

	documentMapping := bleve.NewDocumentMapping()

	// name
	documentMapping.AddFieldMappingsAt("name", englishTextFieldMapping)

	// content
	documentMapping.AddFieldMappingsAt("text", englishTextFieldMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("assets", documentMapping)

	indexMapping.TypeField = "type"
	indexMapping.DefaultAnalyzer = "en"

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
