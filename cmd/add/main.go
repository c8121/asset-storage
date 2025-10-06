package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/metadata"
	metadata_db "github.com/c8121/asset-storage/internal/metadata-db"
	mdsqlite "github.com/c8121/asset-storage/internal/metadata-sqlite"
	"github.com/c8121/asset-storage/internal/storage"
)

const (
	MaxAttemptsPerFile      = 10
	MinWaitSecondsAfterFail = 3
)

var (
	currentUser, currentUserErr = user.Current()
	recursive                   = flag.Bool("r", false, "Recursively add files")
	fileNameFilter              = flag.String("name", "", "File name filter (*.jpg for example")
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
	storage.CreateDirectories()

	mdsqlite.Open()
	defer mdsqlite.Close()

	for _, file := range files {
		if file[:1] == "-" {
			continue
		}
		addErr := addPath(file)
		if addErr != nil {
			fmt.Println(addErr)
		}
	}
}

func addPath(path string) error {

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
			err = addPath(filePath)
			if err != nil {
				fmt.Printf("Error, omitting file: '%s': %s\n", filePath, err)
			}
		}
		return nil
	}

	if fileNameFilter != nil && *fileNameFilter != "" {
		matched, err := filepath.Match(*fileNameFilter, stat.Name())
		if !matched || err != nil {
			return err
		}
	}

	for attempt := 0; attempt < MaxAttemptsPerFile; attempt++ {
		err = addFileAndMetadata(path, stat)
		if err == nil {
			return nil
		}

		wait := time.Duration(attempt*MinWaitSecondsAfterFail) * time.Second
		fmt.Printf("Error, attempt %d/%d, waiting %d: %s\n", attempt+1, MaxAttemptsPerFile, wait, filepath.Base(path))
		time.Sleep(wait)
	}
	return err
}

func addFileAndMetadata(path string, stat os.FileInfo) error {

	//Add file to storage
	infos, err := storage.AddFile(path)
	if err != nil {
		fmt.Printf("Error adding '%s': %s\n", path, err)
		return err
	}

	for _, info := range infos {
		if info.IsNewFile || !config.SkipMetaDataIfExists {
			//Create/Update meta-data
			meta, err := metadata.AddMetaData(
				info.Hash,
				info.MimeType,
				filepath.Base(info.SourcePath),
				filepath.Dir(info.SourcePath),
				currentUser.Username,
				stat.ModTime())
			if err != nil {
				fmt.Printf("Error adding meta-data '%s': %s\n", path, err)
				return err
			}

			//Create/Update meta-data-database
			err = metadata_db.AddMetaData(meta)
			if err != nil {
				fmt.Printf("Error adding meta-data to database '%s': %s\n", path, err)
			}
		}
	}
	return err
}
