package ssh_server

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/c8121/asset-storage/internal/util"
)

func NewVirtualRsyncHandler(rootDirectory string, username string) *VirtualRsyncHandler {
	handler := &VirtualRsyncHandler{
		rootDirectory: rootDirectory,
		username:      username,
	}
	return handler
}

type VirtualRsyncHandler struct {
	rootDirectory string
	username      string
}

func (h *VirtualRsyncHandler) GetUsername() string {
	return h.username
}

func (h *VirtualRsyncHandler) GetNewFiles() []string {
	existingNewFiles := make([]string, 0)

	readFiles(h.rootDirectory, &existingNewFiles)

	return existingNewFiles
}

func readFiles(dirPath string, readInto *[]string) {
	stat, err := os.Stat(dirPath)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("Directory does in exists: " + dirPath)
		return
	}
	if !stat.IsDir() {
		fmt.Println("Not a directory: " + dirPath)
		return
	}

	fis, err := os.Open(dirPath)
	if err != nil {
		fmt.Printf("Failed to open directory %s: %v\n", dirPath, err)
		return
	}
	defer util.CloseOrLog(fis)

	list, err := fis.Readdir(-1)
	if err != nil {
		fmt.Printf("Failed to read directory %s: %v\n", dirPath, err)
		return
	}

	for _, f := range list {
		path := filepath.Join(dirPath, f.Name())
		if f.IsDir() {
			readFiles(path, readInto)
		} else if f.Mode().IsRegular() {
			*readInto = append(*readInto, path)
		}
	}
}

func (h *VirtualRsyncHandler) parseCommand(payload string) (string, []string, error) {

	if parts := strings.Split(payload, " "); len(parts) < 2 {
		return "", nil, errors.New("Invalid payload for rsync: " + payload)
	} else {
		cmd := parts[0]
		args := parts[1:]

		destPath := args[len(args)-1]
		destPath, err := resolve(h.rootDirectory, destPath)
		if err != nil {
			return "", nil, errors.New("Failed to resolve path for rsync: " + payload)
		}
		args[len(args)-1] = destPath

		return cmd, args, nil
	}
}
