package storage

import (
	"os"
)

type StorageReader interface {
	Read([]byte) (int, error)
	Close() error
}

type StorageFileReader struct {
	File *os.File
}

func (reader *StorageFileReader) Read(b []byte) (int, error) {
	return reader.File.Read(b)
}

func (reader *StorageFileReader) Close() error {
	return reader.File.Close()
}

func NewFileReader(path string) (*StorageFileReader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return &StorageFileReader{file}, nil
}
