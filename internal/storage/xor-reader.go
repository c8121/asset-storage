package storage

import (
	"io"

	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/util"
)

// XorReader wraps a StorageReader, xor's bytes on Read(...)
type XorReader struct {
	reader StorageReader
	xor    XorEncoder
}

// Read reads from wrapped reader, xor'ing all bytes
func (r *XorReader) Read(p []byte) (int, error) {
	n, err := r.reader.Read(p)
	if n == 0 && err == io.EOF {
		return n, err
	}
	util.PanicOnIoError(err, "Failed to read bytes")
	r.xor.Encode(p)
	return n, nil
}

// Close closes the wrapped reader
func (r *XorReader) Close() error {
	return r.reader.Close()
}

// NewXorReader creates a new XorReader, wrapping the given StorageReader
func NewXorReader(sr StorageReader) *XorReader {
	return &XorReader{sr, &Xor{config.XorKey, len(config.XorKey), 0}}
}
