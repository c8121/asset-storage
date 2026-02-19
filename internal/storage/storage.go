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
		SourcePath  string
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
// Returns content-hash, file-path, mime-type, error as AddedFileInfo (might be more than one if an archive was added)
func AddFile(path string) ([]AddedFileInfo, error) {

	fmt.Println("Add file:", path)

	stat, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Printf("'%s' does not exist\n", path)
		return nil, os.ErrNotExist
	}

	reader, err := os.Open(path)
	if err != nil {
		fmt.Printf("Cannot open '%s': %s\n", path, err)
		return nil, os.ErrNotExist
	}
	defer util.CloseOrLog(reader)

	infos := make([]AddedFileInfo, 0)

	info, err := copyToStorage(reader, stat.Size())
	if err != nil {
		return nil, err
	}

	info.SourcePath = path
	infos = append(infos, *info)

	if !info.IsNewFile {
		fmt.Printf("File already exists: '%s' '%s'\n", info.SourcePath, info.Hash)
	}

	if IsUnpackable(path, info.MimeType) {
		unpacked, err := Unpack(path, info.MimeType)
		if err == nil {
			for _, item := range unpacked {
				item.SourcePath = path + "/" + item.SourcePath
				infos = append(infos, item)
				fmt.Printf(" '--> %s\n", item.SourcePath)
			}
		} else {
			fmt.Printf("Cannot unpack '%s': %s\n", path, err)
		}
	}

	return infos, nil
}

func copyToStorage(reader io.Reader, size int64) (*AddedFileInfo, error) {

	var info = &AddedFileInfo{IsNewFile: false}

	writer, err := newTempWriter(size)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp-writer: %w", err)
	}
	defer util.CloseOrLog(writer)

	var outWriter StorageWriter
	if len(config.XorKey) > 0 {
		outWriter = NewXorWriter(writer)
	} else {
		outWriter = writer
	}
	defer util.CloseOrLog(outWriter)

	buf := make([]byte, IoBufferSize)
	hash := sha256.New()

	for {
		n, err := reader.Read(buf)
		if n > 0 {
			if len(info.MimeType) == 0 { //must be before outWriter.Write, because buf might get xor'ed
				mime := mimetype.Detect(buf[:n])
				info.MimeType = mime.String()
			}

			hash.Write(buf[:n]) //must be before outWriter.Write

			n, err = outWriter.Write(buf[:n])
			if err != nil {
				return info, fmt.Errorf("failed to write: %w", err)
			}
			info.Size += int64(n)
		}
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return info, fmt.Errorf("failed to read: %w (%d bytes read)", err, info.Size)
			}
		}
	}

	util.CloseOrLog(outWriter)

	info.Hash = fmt.Sprintf("%x", hash.Sum(nil))
	if len(info.Hash) < 2 {
		return info, fmt.Errorf("invalid hash length: %d", len(info.Hash))
	}

	info.StoragePath, err = FindByHash(info.Hash)
	if err == nil {
		info.IsNewFile = false
		util.LogError(outWriter.Remove())
		return info, nil
	}

	info.IsNewFile = true

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
		util.PanicOnError(os.Remove(outWriter.Name()), "Failed to remove temp file")
		panic("File already exists") //Panic, because check was done above with FindByHash
	}

	err = outWriter.Move(info.StoragePath)
	if err != nil {
		return info, fmt.Errorf("failed to move file: %w", err)
	}
	util.LogError(os.Chmod(info.StoragePath, FilePermissions))

	return info, nil
}

func Walk(handler func(path string)) {

	timePeriodDirs, err := os.ReadDir(config.AssetStorageBaseDir)
	if errors.Is(err, os.ErrNotExist) {
		return
	}
	util.PanicOnError(err, "Failed to read directory")

	for _, timePeriodEntry := range timePeriodDirs {

		path := filepath.Join(config.AssetStorageBaseDir, timePeriodEntry.Name())
		children, err := os.ReadDir(path)
		util.PanicOnError(err, "Failed to read directory")

		for _, child := range children {
			dir := filepath.Join(path, child.Name())
			files, err := os.ReadDir(dir)
			util.PanicOnError(err, "Failed to read directory")

			for _, file := range files {
				filePath := filepath.Join(dir, file.Name())
				handler(filePath)
			}
		}

	}
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

// HashFromStoragePath Extract full hash from path (.../hash[:2]/hash[2:]...)
func HashFromStoragePath(path string) string {
	dir, name := filepath.Split(path)
	_, dir2 := filepath.Split(dir[:len(dir)-1])
	hash := dir2 + name
	p := strings.Index(hash, ".")
	if p > -1 {
		return hash[:p]
	}

	return hash
}

// HashFromContent calculates the content hash
func HashFromContent(path string) (string, error) {

	reader, err := os.Open(path)
	if err != nil {
		fmt.Printf("Cannot open '%s': %s\n", path, err)
		return "", os.ErrNotExist
	}
	defer util.CloseOrLog(reader)

	buf := make([]byte, IoBufferSize)
	hash := sha256.New()

	for {
		n, err := reader.Read(buf)
		if n > 0 {
			hash.Write(buf[:n])
		}
		if err != nil && err == io.EOF {
			break
		}
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// Open returns a reader to get asset content.
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
			}
			return reader, nil
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
