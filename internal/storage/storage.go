package storage

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/util"
	"github.com/gabriel-vasile/mimetype"
)

const (
	IoBufferSize    = 4096
	FilePermissions = 0744
)

// Init creates required directories
func Init() {
	util.CreateDirIfNotExists(config.AssetStorageBaseDir, FilePermissions)
	util.CreateDirIfNotExists(config.AssetStorageTempDir, FilePermissions)
}

// AddFile adds one or more file to asset-storage.
// Returns content-hash, file-path, mime-type, error
func AddFile(path string) (assetHash, assetPath, mimeType string, err error) {

	stat, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Printf("'%s' does not exist\n", path)
		return "", "", "", os.ErrNotExist
	}
	fmt.Println("Add file:", path)

	tempDest, err := getTempWriter(stat.Size())
	util.PanicOnError(err, "Failed to create temp file")
	fmt.Println("Temp file created:", tempDest.Name())

	var outWriter io.Writer
	if len(config.XorKey) > 0 {
		outWriter = NewXorWriter(tempDest)
	} else {
		outWriter = tempDest
	}

	buf := make([]byte, IoBufferSize)
	in, err := os.Open(path)
	if err != nil {
		fmt.Printf("Cannot open '%s': %s\n", path, err)
		return "", "", "", os.ErrNotExist
	}
	util.PanicOnError(err, "Failed to open file")
	defer func(in *os.File) {
		util.PanicOnError(in.Close(), "Failed to close file")
	}(in)

	hash := sha256.New()
	mimetypeName := ""
	size := int64(0)

	for {
		n, err := in.Read(buf)
		if err == io.EOF {
			break
		}
		util.PanicOnError(err, "Failed to read file")

		if len(mimetypeName) == 0 { //must be before outWriter.Write, because buf might get xor'ed
			mime := mimetype.Detect(buf[:n])
			mimetypeName = mime.String()
		}

		hash.Write(buf[:n]) //must be before outWriter.Write

		n, err = outWriter.Write(buf[:n])
		util.PanicOnError(err, "Failed to write to temp file")
		size += int64(n)

	}

	util.PanicOnError(tempDest.Close(), "Failed to close temp file")

	fmt.Printf("MIME type: %s, Size: %d\n", mimetypeName, size)

	hashHex := fmt.Sprintf("%x", hash.Sum(nil))
	if len(hashHex) < 2 {
		panic("Invalid hash")
	}

	destPath, err := FindByHash(hashHex)
	if err == nil {
		fmt.Printf("File already exists: '%s'\n", destPath)
		util.PanicOnError(tempDest.Remove(), "Failed to remove temp file")
		return hashHex, destPath, mimetypeName, nil
	}

	destName := hashHex[2:]
	destDir := fmt.Sprintf("%s/%s/%s",
		config.AssetStorageBaseDir,
		TimePeriodName(),
		hashHex[:2])

	util.PanicOnError(os.MkdirAll(destDir, FilePermissions), "Failed to create destination directory")

	destPath = fmt.Sprintf("%s/%s",
		destDir,
		destName)

	if _, err := os.Stat(destPath); err == nil || os.IsExist(err) {
		util.PanicOnError(os.Remove(tempDest.Name()), "Failed to remove temp file")
		panic("File already exists") //Panic, because check was done above with FindByHash
	}

	fmt.Printf("Adding '%s' to %s\n", path, destPath)
	util.PanicOnError(tempDest.Move(destPath), "Failed to move temp file")
	util.PanicOnError(os.Chmod(destPath, FilePermissions), "Failed to set permissions")

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
		config.AssetStorageBaseDir,
		TimePeriodName(),
		hashHex[:2])

	destPath := filepath.Join(
		destDir,
		destName)
	if _, err := os.Stat(destPath); err == nil || os.IsExist(err) {
		return destPath, nil
	}

	//Evaluate all time-period dirs
	dirs, err := os.ReadDir(config.AssetStorageBaseDir)
	if errors.Is(err, os.ErrNotExist) {
		return destPath, err
	}
	util.PanicOnError(err, "Failed to read directory")
	for _, file := range dirs {
		destDir = filepath.Join(
			config.AssetStorageBaseDir,
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

// HashFromPath Extract full hash from path (.../hash[:2]/hash[2:]...)
func HashFromPath(path string) string {
	dir, name := filepath.Split(path)
	_, dir2 := filepath.Split(dir[:len(dir)-1])
	hash := dir2 + name
	p := strings.Index(hash, ".")
	if p > -1 {
		return hash[:p]
	} else {
		return hash
	}

}

// Open returns asset content
func Open(assetHash string) (StorageReader, error) {
	if path, err := FindByHash(assetHash); err == nil {

		var reader StorageReader
		if config.UseGzip {
			reader, err = NewZipFileReader(path)
		} else {
			reader, err = NewFileReader(path)
		}

		if err == nil {
			if len(config.XorKey) > 0 {
				return NewXorReader(reader), nil
			} else {
				return reader, nil
			}
		}
	}

	return nil, os.ErrNotExist
}

// TimePeriodName Create a name corresponding to period in time (each 4 hours having same name)
func TimePeriodName() string {
	ts := time.Now().UnixMilli() / 1000 / 60 / 24 / 4
	s := fmt.Sprintf("%x", ts)
	return s
}

func getTempWriter(size int64) (StorageWriter, error) {

	var writer StorageWriter
	var err error

	if config.UseGzip {
		if size <= config.MaxMemFileSize {
			writer, err = NewMemZipFileWriter()
		} else {
			writer, err = NewTempZipFileWriter()
		}
	} else if size <= config.MaxMemFileSize {
		writer, err = NewMemFileWriter()
	} else {
		writer, err = NewTempFileWriter()
	}

	return writer, err
}
