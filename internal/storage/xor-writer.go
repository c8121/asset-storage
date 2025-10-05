package storage

import (
	"github.com/c8121/asset-storage/internal/config"
)

// XorWriter wraps a StorageWriter, xor's bytes on Write(...)
type XorWriter struct {
	writer StorageWriter
	xor    XorEncoder
}

// StorageWriter implementation:

func (w *XorWriter) Name() string {
	return w.writer.Name()
}

// Write writes to wrapped writer, xor.ing all bytes
func (w *XorWriter) Write(p []byte) (int, error) {
	w.xor.Encode(p)
	return w.writer.Write(p)
}

func (w *XorWriter) Close() error {
	return w.writer.Close()
}

func (w *XorWriter) Move(path string) error {
	return w.writer.Move(path)
}

func (w *XorWriter) Remove() error {
	return w.writer.Remove()
}

// NewXorWriter creates a new XorWriter, wrapping the given StorageWriter
func NewXorWriter(sw StorageWriter) *XorWriter {
	return &XorWriter{sw, &Xor{config.XorKey, len(config.XorKey), 0}}
}
