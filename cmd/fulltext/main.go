package main

import (
	"fmt"
	"strings"

	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/filter"
	"github.com/c8121/asset-storage/internal/fulltext"
	"github.com/c8121/asset-storage/internal/metadata"
	"github.com/c8121/asset-storage/internal/storage"
	"github.com/c8121/asset-storage/internal/util"
)

func main() {

	config.LoadDefault()

	fulltext.CreateDirectories()
	defer fulltext.CloseIndex()

	handler := func(path string) {
		hash := storage.HashFromStoragePath(path)

		assetMeta, err := metadata.LoadByHash(hash)
		if err != nil {
			util.LogError(err)
			return
		}

		if assetMeta == nil {
			util.LogError(fmt.Errorf("asset metadata not found: %s", hash))
			return
		}

		if strings.Contains(assetMeta.MimeType, "image/") {
			//fmt.Printf("Ignore %s, mime-type: %s\n", hash, assetMeta.MimeType)
			return
		}

		textFilter := filter.GetFirstFilterByNameAndMimeType("Text", assetMeta.MimeType)
		if textFilter == nil {
			//fmt.Printf("No Text filter found: %s, %s\n", hash, assetMeta.MimeType)
			return
		}

		//fmt.Printf("%s: %v\t", assetMeta.Hash, textFilter)

		params := make(map[string]string)
		params["lang"] = "deu" //TODO

		bytes, _, err := textFilter.Apply(hash, assetMeta, params)
		if err != nil {
			util.LogError(err)
			return
		}

		err = fulltext.AddText(hash, string(bytes))
		if err != nil {
			util.LogError(err)
		}
	}

	storage.Walk(handler)
}
