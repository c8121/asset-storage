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

	tempDest, err := newTempWriter(stat.Size())
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
	defer util.CloseOrLog(in)

	hash := sha256.New()
	mimetypeName := ""
	size := int64(0)

	for {
		n, err := in.Read(buf)
		if n == 0 && err == io.EOF {
			break
		} else if err != nil {
			return "", "", "", fmt.Errorf("failed to read: %w", err)
		}

		if len(mimetypeName) == 0 { //must be before outWriter.Write, because buf might get xor'ed
			mime := mimetype.Detect(buf[:n])
			mimetypeName = mime.String()
		}

		hash.Write(buf[:n]) //must be before outWriter.Write

		n, err = outWriter.Write(buf[:n])
		if err != nil {
			return "", "", "", fmt.Errorf("failed to write: %w", err)
		}
		size += int64(n)

	}

	util.CloseOrLog(tempDest)

	fmt.Printf("MIME type: %s, Size: %d\n", mimetypeName, size)

	hashHex := fmt.Sprintf("%x", hash.Sum(nil))
	if len(hashHex) < 2 {
		return "", "", "", fmt.Errorf("invalid hash length: %d", len(hashHex))
	}

	destPath, err := FindByHash(hashHex)
	if err == nil {
		fmt.Printf("File already exists: '%s'\n", destPath)
		util.LogError(tempDest.Remove())
		return hashHex, destPath, mimetypeName, nil
	}

	destName := hashHex[2:]
	destDir := fmt.Sprintf("%s/%s/%s",
		config.AssetStorageBaseDir,
		TimePeriodName(),
		hashHex[:2])

	err = os.MkdirAll(destDir, FilePermissions)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to create directory: %w", err)
	}

	destPath = fmt.Sprintf("%s/%s",
		destDir,
		destName)

	if _, err := os.Stat(destPath); err == nil || os.IsExist(err) {
		util.PanicOnError(os.Remove(tempDest.Name()), "Failed to remove temp file")
		panic("File already exists") //Panic, because check was done above with FindByHash
	}

	fmt.Printf("Adding '%s' to %s\n", path, destPath)
	err = tempDest.Move(destPath)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to move file: %w", err)
	}
	util.LogError(os.Chmod(destPath, FilePermissions))

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

// newTempWriter creates either
//   - NewTempZipFileWriter or NewMemZipFileWriter if congig.UseGzip is true
//   - NewTempFileWriter or NewMemFileWriter
func newTempWriter(size int64) (StorageWriter, error) {

	var writer StorageWriter
	var err error

	if config.UseGzip {
		if size <= config.MaxMemFileSize {
			writer, err = NewMemZipFileWriter(size)
		} else {
			writer, err = NewTempZipFileWriter()
		}
	} else if size <= config.MaxMemFileSize {
		writer, err = NewMemFileWriter(size)
	} else {
		writer, err = NewTempFileWriter()
	}

	return writer, err
}
