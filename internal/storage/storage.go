package storage

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/c8121/asset-storage/internal/util"
	"github.com/gabriel-vasile/mimetype"
)

const (
	IoBufferSize = 4096
	FilePermissions = 0744
)

// BaseDir Directory for all contents of asset storage.
func BaseDir() string {
	return "/tmp/asset-storage"
}

// AddFile Add one file to asset-storage
// Returns content-hash, file-path, mime-type, error
func AddFile(path string) (assetHash, assetPath, mimeType string, err error) {

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("'%s' does not exist\n", path)
		return "", "", "", os.ErrNotExist
	}
	fmt.Println("Add file:", path)

	tempDest := TempFile()
	fmt.Println("Temp file created:", tempDest.Name())

	buf := make([]byte, IoBufferSize)
	in, err := os.Open(path)
	util.Check(err, "Failed to open file")
	defer func(in *os.File) {
		util.Check(in.Close(), "Failed to close file")
	}(in)

	hash := sha256.New()
	mimetypeName := ""

	for {
		n, err := in.Read(buf)
		if err == io.EOF {
			break
		}
		util.Check(err, "Failed to read file")

		n, err = tempDest.Write(buf[:n])
		util.Check(err, "Failed to write to temp file")

		if len(mimetypeName) == 0 {
			mime := mimetype.Detect(buf[:n])
			mimetypeName = mime.String()
			fmt.Println("MIME type:", mimetypeName)
		}

		hash.Write(buf[:n])
	}

	util.Check(tempDest.Close(), "Failed to close temp file")

	hashHex := fmt.Sprintf("%x", hash.Sum(nil))
	if len(hashHex) < 2 {
		panic("Invalid hash")
	}
	fmt.Println("Hash:", hashHex)

	destPath, err := FindByHash(hashHex)
	if err == nil {
		fmt.Printf("File already exists: '%s'\n", destPath)
		util.Check(os.Remove(tempDest.Name()), "Failed to remove temp file")
		return hashHex, destPath, mimetypeName, nil
	}

	destName := hashHex[2:]
	destDir := fmt.Sprintf("%s/%s/%s",
		BaseDir(),
		TimePeriodName(),
		hashHex[:2])

	util.Check(os.MkdirAll(destDir, FilePermissions), "Failed to create destination directory")

	destPath = fmt.Sprintf("%s/%s",
		destDir,
		destName)

	if _, err := os.Stat(destPath); err == nil || os.IsExist(err) {
		util.Check(os.Remove(tempDest.Name()), "Failed to remove temp file")
		panic("File already exists") //Panic, because check was done above with FindByHash
	}

	fmt.Printf("Adding '%s' to %s\n", path, destPath)
	util.Check(os.Rename(tempDest.Name(), destPath), "Failed to move temp file")
	util.Check(os.Chmod(destPath, FilePermissions), "Failed to set permissions")

	return hashHex, destPath, mimetypeName, nil
}

// FindByHash Check all time-periods if file exists
func FindByHash(hashHex string) (assetPath string, err error) {

	if len(hashHex) < 2 {
		panic("Invalid hash")
	}

	//Fast check first: Check if file exists in current time-period
	destName := hashHex[2:]
	destDir := filepath.Join(
		BaseDir(),
		TimePeriodName(),
		hashHex[:2])

	destPath := filepath.Join(
		destDir,
		destName)
	if _, err := os.Stat(destPath); err == nil || os.IsExist(err) {
		return destPath, nil
	}

	//Evaluate all time-period dirs
	dirs, err := os.ReadDir(BaseDir())
	if errors.Is(err, os.ErrNotExist) {
		return destPath, err
	}
	util.Check(err, "Failed to read directory")
	for _, file := range dirs {
		destDir = filepath.Join(
			BaseDir(),
			file.Name(),
			hashHex[:2])
		destPath = filepath.Join(
			destDir,
			destName)
		if _, err := os.Stat(destPath); err == nil || os.IsExist(err) {
			return destPath, nil
		}
	}

	return "", os.ErrNotExist
}

// TimePeriodName Create a name corresponding to period in time (each 4 hours having same name)
func TimePeriodName() string {
	ts := time.Now().UnixMilli() / 1000 / 60 / 24 / 4
	s := fmt.Sprintf("%x", ts)
	return s
}

// TempDir Temporary dir. Should be on same drive as BaseDir()
func TempDir() string {
	return "/tmp"
}

// TempFile Create temp file of panic
func TempFile() *os.File {
	file, err := os.CreateTemp(TempDir(), "asset-*.tmp")
	if err != nil {
		panic(err)
	}
	return file
}
