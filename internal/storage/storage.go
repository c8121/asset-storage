package storage

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/c8121/asset-storage/internal/util"
)

const (
	IoBufferSize = 4096
)

// BaseDir Directory for all contents of asset storage.
func BaseDir() string {
	return "/tmp/asset-storage"
}

// AddFile Add one file to asset-storage
func AddFile(path string) {

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("%s does not exist\n", path)
		return
	}

	tempDest := TempFile()
	fmt.Println("Temp file created:", tempDest.Name())

	buf := make([]byte, IoBufferSize)
	in, err := os.Open(path)
	util.Check(err, "Failed to open file")
	defer func(in *os.File) {
		util.Check(in.Close(), "Failed to close file")
	}(in)

	hash := sha256.New()

	for {
		n, err := in.Read(buf)
		if err == io.EOF {
			break
		}
		util.Check(err, "Failed to read file")

		n, err = tempDest.Write(buf[:n])
		util.Check(err, "Failed to write to temp file")

		hash.Write(buf[:n])
	}

	util.Check(tempDest.Close(), "Failed to close temp file")

	destName := fmt.Sprintf("%x", hash.Sum(nil))
	destDir := fmt.Sprintf("%s/%s/%s",
		BaseDir(),
		TimePeriodName(),
		destName[:2])

	util.Check(os.MkdirAll(destDir, os.ModeDir), "Failed to create destination directory")

	destPath := fmt.Sprintf("%s/%s",
		destDir,
		destName[2:])

	if _, err := os.Stat(destPath); err == nil || os.IsExist(err) {
		fmt.Printf("%s already exists\n", destPath)
		return
	}

	fmt.Printf("Adding %s to %s\n", path, destPath)
	util.Check(os.Rename(tempDest.Name(), destPath), "Failed to move temp file")
}

// TimePeriodName Create a name corresponding to period in time (each 4 hours having same name)
func TimePeriodName() string {
	ts := time.Now().UnixMilli() / 1000 / 60 / 24 / 4
	s := fmt.Sprintf("%x", ts)
	return s
}

// TempDir Temporary dir
func TempDir() string {
	return "/tmp"
}

// TempFile Create temp file of panic
func TempFile() *os.File {
	file, err := os.CreateTemp(TempDir(), "*.tmp")
	if err != nil {
		panic(err)
	}
	return file
}
