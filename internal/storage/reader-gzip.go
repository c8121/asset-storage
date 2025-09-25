package storage

import (
	"compress/gzip"
	"os"
)

type StorageZipFileReader struct {
	File      *os.File
	zipReader *gzip.Reader
}

func (reader *StorageZipFileReader) Read(b []byte) (int, error) {
	return reader.zipReader.Read(b)
}

func (reader *StorageZipFileReader) Close() error {
	if err := reader.zipReader.Close(); err != nil {
		return err
	}
	return reader.File.Close()
}

func NewZipFileReader(path string) (*StorageZipFileReader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	zip := &StorageZipFileReader{}
	zip.File = file
	zip.zipReader, err = gzip.NewReader(file)
	if err != nil {
		return nil, err
	}

	return zip, nil
}
