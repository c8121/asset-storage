package storage

import (
	"bytes"
	"os"
)

type StorageMemFileWriter struct {
	buf *bytes.Buffer
}

func (writer *StorageMemFileWriter) Name() string {
	return "InMemory"
}

func (writer *StorageMemFileWriter) Write(b []byte) (int, error) {
	return writer.buf.Write(b)
}

func (writer *StorageMemFileWriter) Close() error {
	return nil
}

func (writer *StorageMemFileWriter) Move(path string) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = out.Write(writer.buf.Bytes())
	return err
}

func (writer *StorageMemFileWriter) Remove() error {
	return nil
}

func NewMemFileWriter(size int64) (*StorageMemFileWriter, error) {
	return &StorageMemFileWriter{bytes.NewBuffer(make([]byte, 0, size))}, nil
}
