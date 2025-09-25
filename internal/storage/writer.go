package storage

import (
	"os"

	"github.com/c8121/asset-storage/internal/config"
)

type StorageWriter interface {
	Name() string
	Write([]byte) (int, error)
	Close() error
	Move(path string) error
}

type StorageFileWriter struct {
	File *os.File
}

func (writer *StorageFileWriter) Name() string {
	return writer.File.Name()
}

func (writer *StorageFileWriter) Write(b []byte) (int, error) {
	return writer.File.Write(b)
}

func (writer *StorageFileWriter) Close() error {
	return writer.File.Close()
}

func (writer *StorageFileWriter) Move(path string) error {
	return os.Rename(writer.Name(), path)
}

func NewTempFileWriter() (*StorageFileWriter, error) {
	file, err := os.CreateTemp(config.AssetStorageTempDir, "asset-*.tmp")
	if err != nil {
		return nil, err
	}

	return &StorageFileWriter{file}, nil
}
