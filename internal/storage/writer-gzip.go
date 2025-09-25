package storage

import (
	"compress/gzip"
	"os"

	"github.com/c8121/asset-storage/internal/config"
)

type StorageZipFileWriter struct {
	File      *os.File
	zipWriter *gzip.Writer
}

func (writer *StorageZipFileWriter) Name() string {
	return writer.File.Name()
}

func (writer *StorageZipFileWriter) Write(b []byte) (int, error) {
	return writer.zipWriter.Write(b)
}

func (writer *StorageZipFileWriter) Close() error {
	if err := writer.zipWriter.Flush(); err != nil {
		return err
	}
	if err := writer.zipWriter.Close(); err != nil {
		return err
	}
	return writer.File.Close()
}

func (writer *StorageZipFileWriter) Move(path string) error {
	return os.Rename(writer.Name(), path)
}

func (writer *StorageZipFileWriter) Remove() error {
	return os.Remove(writer.Name())
}

func NewTempZipFileWriter() (*StorageZipFileWriter, error) {
	file, err := os.CreateTemp(config.AssetStorageTempDir, "asset-*.tmp")
	if err != nil {
		return nil, err
	}

	zip := &StorageZipFileWriter{}
	zip.File = file
	zip.zipWriter = gzip.NewWriter(file)

	return zip, nil
}
