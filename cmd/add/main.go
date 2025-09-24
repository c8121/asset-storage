package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/metadata"
	"github.com/c8121/asset-storage/internal/storage"
)

var (
	currentUser, currentUserErr = user.Current()
	recursive                   = flag.Bool("r", false, "Recursively add files")
)

func main() {
	flag.Parse()
	files := flag.Args()
	if len(files) == 0 {
		fmt.Printf("You must specify at least one file: %s file [file...]\n", filepath.Base(os.Args[0]))
	}

	if currentUserErr != nil {
		panic(currentUserErr)
	}

	config.LoadDefault()
	storage.Init()

	for _, file := range files {
		if file[:1] == "-" {
			continue
		}
		addErr := addFile(file)
		if addErr != nil {
			fmt.Println(addErr)
		}
	}
}

func addFile(path string) error {

	stat, err := os.Stat(path)
	if err != nil {
		fmt.Printf("Cannot get file stat '%s'\n", path)
		return err
	}

	if stat.IsDir() {
		if !*recursive {
			fmt.Printf("'%s' is a directory, omitting (-r not given)\n", path)
			return nil
		}
		files, err := os.ReadDir(path)
		if err != nil {
			return err
		}
		for _, file := range files {
			filePath := filepath.Join(path, file.Name())
			err = addFile(filePath)
			if err != nil {
				fmt.Printf("Error, omitting file: '%s': %s\n", filePath, err)
			}
		}
		return nil
	}

	//Add file to storage
	assetHash, _, mimeType, err := storage.AddFile(path)
	if err != nil {
		fmt.Printf("Error adding '%s': %s\n", path, err)
		return err
	}

	//Create/Update meta-data
	metaDataFile := metadata.GetMetaDataFilePath(assetHash)
	fmt.Printf("MetaDataFile: %s\n", metaDataFile)

	metaData, err := metadata.LoadIfExists(metaDataFile)
	if errors.Is(err, os.ErrNotExist) {
		metaData = metadata.CreateNew(
			mimeType,
			filepath.Base(path),
			filepath.Dir(path),
			currentUser.Username,
			stat.ModTime())
	} else if err != nil {
		fmt.Printf("Error adding meta-data '%s': %s\n", metaDataFile, err)
	} else {
		metaData.AddOrigin(
			filepath.Base(path),
			filepath.Dir(path),
			currentUser.Username,
			stat.ModTime())
	}

	//fmt.Printf("MetaData: %s\n", metaData)
	return metaData.Save(metaDataFile)
}
