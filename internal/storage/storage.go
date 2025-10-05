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

type (
	AddedFileInfo struct {
		Hash        string
		StoragePath string
		MimeType    string
		IsNewFile   bool
		Size        int64
	}
)

// CreateDirectories creates required directories
func CreateDirectories() {
	util.CreateDirIfNotExists(config.AssetStorageBaseDir, FilePermissions)
	util.CreateDirIfNotExists(config.AssetStorageTempDir, FilePermissions)
}

// AddFile adds one file to asset-storage.
// Returns content-hash, file-path, mime-type, error
func AddFile(path string) (*AddedFileInfo, error) {

	writer, err := createTempWriter(path)
	if err != nil {
		return nil, err
	}

	reader, err := os.Open(path)
	if err != nil {
		fmt.Printf("Cannot open '%s': %s\n", path, err)
		return nil, os.ErrNotExist
	}
	defer util.CloseOrLog(reader)

	info, err := moveToStorage(reader, writer)
	if err != nil {
		return nil, err
	}

	return info, nil
}

// createTempWriter creates a new StorageWriter to a temp file
func createTempWriter(path string) (StorageWriter, error) {

	stat, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Printf("'%s' does not exist\n", path)
		return nil, os.ErrNotExist
	}
	fmt.Println("Add file:", path)

	tempDest, err := newTempWriter(stat.Size())
	if err != nil {
		return nil, fmt.Errorf("failed to create temp-writer: %w", err)
	}

	var outWriter StorageWriter
	if len(config.XorKey) > 0 {
		outWriter = NewXorWriter(tempDest)
	} else {
		outWriter = tempDest
	}

	return outWriter, nil
}

func moveToStorage(reader io.Reader, writer StorageWriter) (*AddedFileInfo, error) {

	var info = &AddedFileInfo{IsNewFile: false}

	buf := make([]byte, IoBufferSize)

	hash := sha256.New()

	for {
		n, err := reader.Read(buf)
		if n == 0 && err == io.EOF {
			break
		} else if err != nil {
			return info, fmt.Errorf("failed to read: %w", err)
		}

		if len(info.MimeType) == 0 { //must be before outWriter.Write, because buf might get xor'ed
			mime := mimetype.Detect(buf[:n])
			info.MimeType = mime.String()
		}

		hash.Write(buf[:n]) //must be before outWriter.Write

		n, err = writer.Write(buf[:n])
		if err != nil {
			return info, fmt.Errorf("failed to write: %w", err)
		}
		info.Size += int64(n)

	}

	util.CloseOrLog(writer)

	info.Hash = fmt.Sprintf("%x", hash.Sum(nil))
	if len(info.Hash) < 2 {
		return info, fmt.Errorf("invalid hash length: %d", len(info.Hash))
	}

	var err error
	info.StoragePath, err = FindByHash(info.Hash)
	if err == nil {
		info.IsNewFile = false
		fmt.Printf("File already exists: '%s'\n", info.StoragePath)
		util.LogError(writer.Remove())
		return info, nil
	}

	destName := info.Hash[2:]
	destDir := fmt.Sprintf("%s/%s/%s",
		config.AssetStorageBaseDir,
		TimePeriodName(),
		info.Hash[:2])

	err = os.MkdirAll(destDir, FilePermissions)
	if err != nil {
		return info, fmt.Errorf("failed to create directory: %w", err)
	}

	info.StoragePath = fmt.Sprintf("%s/%s",
		destDir,
		destName)

	if _, err := os.Stat(info.StoragePath); err == nil || os.IsExist(err) {
		util.PanicOnError(os.Remove(writer.Name()), "Failed to remove temp file")
		panic("File already exists") //Panic, because check was done above with FindByHash
	}

	err = writer.Move(info.StoragePath)
	if err != nil {
		return info, fmt.Errorf("failed to move file: %w", err)
	}
	util.LogError(os.Chmod(info.StoragePath, FilePermissions))

	return info, nil
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
//   - NewTempZipFileWriter or NewMemZipFileWriter if config.UseGzip is true
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
