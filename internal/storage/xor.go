package storage

import (
	"io"

	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/util"
)

type XorEncoder interface {
	Encode(b []byte)
}

type Xor struct {
	Key []byte
	kl  int //key length
	ki  int //key index
}

// XorReader wraps a StorageReader, xor's bytes on Read(...)
type XorReader struct {
	reader StorageReader
	xor    XorEncoder
}

// XorWriter wraps a StorageWriter, xor's bytes on Write(...)
type XorWriter struct {
	writer StorageWriter
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

// Write writes to wrapped writer, xor.ing all bytes
func (w *XorWriter) Write(p []byte) (int, error) {
	w.xor.Encode(p)
	return w.writer.Write(p)
}

// NewXorWriter creates a new XorWriter, wrapping the given StorageWriter
func NewXorWriter(sw StorageWriter) *XorWriter {
	return &XorWriter{sw, &Xor{config.XorKey, len(config.XorKey), 0}}
}

// Encode does b ^ e.Key if key length greater 0
func (e *Xor) Encode(b []byte) {

	if e.kl == 0 {
		return
	}

	for i := range b {

		b[i] ^= e.Key[e.ki]

		if e.ki++; e.ki == e.kl {
			e.ki = 0
		}
	}
}
