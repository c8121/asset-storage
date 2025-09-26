package storage

import (
	"github.com/c8121/asset-storage/internal/config"
)

// XorReader wraps a StorageReader, xor's bytes on Read(...)
type XorReader struct {
	reader StorageReader
	buf    []byte
}

// XorReader wraps a StorageWriter, xor's bytes on Write(...)
type XorWriter struct {
	writer StorageWriter
	buf    []byte
}

// Read reads from wrapped reader, xor'ing all bytes
func (r *XorReader) Read(p []byte) (int, error) {

	if len(config.XorKey) == 0 {
		return r.reader.Read(p)
	}

	min := min(len(r.buf), len(p))

	n, err := r.reader.Read(r.buf[:min])
	if err != nil {
		return n, err
	}
	xor(r.buf[:n])
	copy(p, r.buf[:n])
	return n, err
}

// Close closes the wrapped reader
func (r *XorReader) Close() error {
	return r.reader.Close()
}

// NewXorReader creates a new XorReader, wrapping the given StorageReader
func NewXorReader(sr StorageReader) *XorReader {
	return &XorReader{sr, make([]byte, len(config.XorKey))}
}

// Write writes to wrapped writer, xor.ing all bytes
func (w *XorWriter) Write(p []byte) (int, error) {

	if len(config.XorKey) == 0 {
		return w.writer.Write(p)
	}

	n := len(p)
	bn := len(w.buf)
	for off := 0; off < n; {
		i := 0
		for ; i < bn && off < n; i++ {
			w.buf[i] = p[off]
			off++
		}
		xor(w.buf[:i])
		_, err := w.writer.Write(w.buf[:i])
		if err != nil {
			return 0, err
		}
	}
	return len(p), nil
}

// NewXorWriter creates a new XorWriter, wrapping the given StorageWriter
func NewXorWriter(sw StorageWriter) *XorWriter {
	return &XorWriter{sw, make([]byte, len(config.XorKey))}
}

// xor if config.XorKey length greater 0
func xor(b []byte) {

	kl := len(config.XorKey)
	if kl == 0 {
		return
	}

	ki := 0

	for i := range b {

		b[i] ^= config.XorKey[ki]

		if ki++; ki == kl {
			ki = 0
		}
	}
}
