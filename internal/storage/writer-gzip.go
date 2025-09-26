package storage

import (
	"bytes"
	"compress/gzip"
	"io"
	"os"

	"github.com/c8121/asset-storage/internal/config"
)

type StorageZipFileWriter struct {
	File      *os.File
	zipWriter *gzip.Writer
	outWriter io.Writer
}

func (writer *StorageZipFileWriter) Name() string {
	if writer.File == nil {
		return "InMemoryZip"
	}
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
	if writer.File != nil {
		return writer.File.Close()
	}
	return nil
}

func (writer *StorageZipFileWriter) Move(path string) error {
	if writer.File != nil {
		return os.Rename(writer.Name(), path)
	} else {
		out, err := os.Create(path)
		if err != nil {
			return err
		}
		defer out.Close()

		buf := writer.outWriter.(*bytes.Buffer)
		_, err = out.Write(buf.Bytes())
		return err
	}
}

func (writer *StorageZipFileWriter) Remove() error {
	if writer.File != nil {
		return os.Remove(writer.Name())
	}
	return nil
}

func NewTempZipFileWriter() (*StorageZipFileWriter, error) {
	file, err := os.CreateTemp(config.AssetStorageTempDir, "asset-*.tmp")
	if err != nil {
		return nil, err
	}

	zip := &StorageZipFileWriter{}
	zip.File = file
	zip.outWriter = file
	zip.zipWriter = gzip.NewWriter(zip.outWriter)

	return zip, nil
}

func NewMemZipFileWriter() (*StorageZipFileWriter, error) {

	buf := bytes.NewBuffer(make([]byte, 0, config.MaxMemFileSize))

	zip := &StorageZipFileWriter{}
	zip.File = nil
	zip.outWriter = buf
	zip.zipWriter = gzip.NewWriter(zip.outWriter)

	return zip, nil
}
