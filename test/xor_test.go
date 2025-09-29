package test

import (
	"fmt"
	"io"
	"testing"

	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/storage"
	"github.com/c8121/asset-storage/internal/util"
)

type TestReaderWriter struct {
	Data []byte
}

func (w *TestReaderWriter) Read(p []byte) (int, error) {
	if len(w.Data) == 0 {
		return 0, io.EOF
	}
	i := 0
	for ; i < len(p) && i < len(w.Data); i++ {
		p[i] = w.Data[i]
	}
	w.Data = w.Data[i:]
	return i, nil
}

func (w *TestReaderWriter) Move(path string) error {
	panic("unimplemented")
}

func (w *TestReaderWriter) Name() string {
	panic("unimplemented")
}

func (w *TestReaderWriter) Remove() error {
	panic("unimplemented")
}

func (w *TestReaderWriter) Write(p []byte) (int, error) {
	w.Data = append(w.Data, p...)
	return len(p), nil
}

func (w *TestReaderWriter) Close() error {
	panic("unimplemented")
}

func TestXor(t *testing.T) {

	readerTest(t)

	encodeTest(t, "123", "abc")
	encodeTest(t, "1234", "abcde")
	encodeTest(t, "12345", "abcdefghi")
	encodeTest(t, "1234", "a")
	encodeTest(t, "01234567890123456789012345678901234567890123456789012345678901234567890123456789", "AbCs")
	encodeTest(t, "AbCs", "01234567890123456789012345678901234567890123456789012345678901234567890123456789")

}

func readerTest(t *testing.T) {

	s := "Hello World"

	trw := &TestReaderWriter{make([]byte, 0)}
	trw.Write([]byte(s))
	if string(trw.Data) != s {
		t.Errorf("Data not written correctly")
	}

	out := make([]byte, 20)
	os := make([]byte, 0)
	for {
		n, err := trw.Read(out)
		if n == 0 && err == io.EOF {
			break
		}
		util.PanicOnIoError(err, "Failed to read bytes")
		os = append(os, out[:n]...)
	}

	ts := string(os)
	if s != ts {
		t.Errorf("Data not read correctly: '%s'\n", ts)
	}

}

func encodeTest(t *testing.T, s string, k string) {

	config.XorKey = []byte(k)

	in := []byte(s)

	trw := &TestReaderWriter{make([]byte, 0)}
	w := storage.NewXorWriter(trw)

	n, _ := w.Write(in)
	fmt.Printf("Written %d bytes\n", n)
	if n != len(s) {
		t.Errorf("%d bytes written, but data length is %d", n, len(s))
	}

	ts := string(trw.Data)
	if s == ts {
		t.Errorf("Data not xor'ed, still equal")
	}

	r := storage.NewXorReader(trw)

	ts = ""
	buf := make([]byte, 3)
	for {
		n, err := r.Read(buf)
		if n == 0 && err == io.EOF {
			break
		}
		util.PanicOnIoError(err, "Failed to read file")
		ts += string(buf[:n])
	}

	if s != ts {
		t.Errorf("Data not restored: '%s'\n", ts)
	}
}
