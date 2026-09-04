package storage

import (
	"errors"
	"fmt"
	"os"

	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/util"
)

// AddFile adds one file to asset-storage.
// Returns content-hash, file-path, mime-type, error as AddedFileInfo (might be more than one if an archive was added)
func AddFile(path string) ([]AddedFileInfo, error) {

	fmt.Println("Add file:", path)

	stat, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Printf("'%s' does not exist\n", path)
		return nil, os.ErrNotExist
	}

	reader, err := os.Open(path)
	if err != nil {
		fmt.Printf("Cannot open '%s': %s\n", path, err)
		return nil, os.ErrNotExist
	}
	defer util.CloseOrLog(reader)

	infos := make([]AddedFileInfo, 0)

	info, err := copyToStorage(reader, stat.Size())
	if err != nil {
		return nil, err
	}

	info.SourcePath = path
	infos = append(infos, *info)

	if !info.IsNewFile {
		fmt.Printf("File already exists: '%s' '%s'\n", info.SourcePath, info.Hash)
	}

	if IsUnpackable(path, info.MimeType) {
		unpacked, err := Unpack(path, info.MimeType)
		if err == nil {
			for _, item := range unpacked {
				item.SourcePath = path + "/" + item.SourcePath
				infos = append(infos, item)
				fmt.Printf(" '--> %s\n", item.SourcePath)
			}
		} else {
			fmt.Printf("Cannot unpack '%s': %s\n", path, err)
		}
	}

	return infos, nil
}

// Open returns a reader to get asset content.
func Open(assetHash string) (StorageReader, error) {
	if path, err := FindByHash(assetHash); err == nil {

		var reader StorageReader
		if config.UseGzip {
			reader, err = NewZipFileReader(path)
		} else {
			reader, err = NewFileReader(path)
		}

		if err == nil {
			if len(config.XorKey) > 0 {
				return NewXorReader(reader), nil
			}
			return reader, nil
		}
	}

	return nil, os.ErrNotExist
}

// Remove deletes the asset (file)
func Remove(assetHash string) error {

	if path, err := FindByHash(assetHash); err == nil {
		return os.Remove(path)
	}

	return os.ErrNotExist
}
