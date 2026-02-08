package ssh_server

import (
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	"github.com/c8121/asset-storage/internal/util"
	"github.com/pkg/sftp"
)

var (
	ErrorInvalidPathRequest = errors.New("Invalid path request")
)

func NewVirtualSftpHandler(rootDirectory string, username string) *VirtualSftpHandler {
	handler := &VirtualSftpHandler{
		rootDirectory: rootDirectory,
		username:      username,
		permissions:   0700, //Create permissions
		newFiles:      make([]string, 0),
	}
	return handler
}

type VirtualSftpHandler struct {
	rootDirectory string
	username      string
	permissions   os.FileMode

	newFiles []string
}

func (h *VirtualSftpHandler) GetHandlers() sftp.Handlers {
	return sftp.Handlers{
		FileGet:  h,
		FilePut:  h,
		FileCmd:  h,
		FileList: h,
	}
}

func (h *VirtualSftpHandler) GetUsername() string {
	return h.username
}

func (h *VirtualSftpHandler) GetNewFiles() []SshFileInfo {
	existingNewFiles := make([]SshFileInfo, 0)

	for _, path := range h.newFiles {
		if stat, err := os.Stat(path); err == nil {
			if stat.Mode().IsRegular() {
				info := SshFileInfo{
					LocalPath: path,
					UserPath:  path[len(h.rootDirectory):],
				}
				existingNewFiles = append(existingNewFiles, info)
			}
		}
	}

	return existingNewFiles
}

func (h *VirtualSftpHandler) Fileread(r *sftp.Request) (io.ReaderAt, error) {
	path, err := h.resolve(r.Filepath)
	if err != nil {
		fmt.Printf("Fileread: Resolve failed: %s (%s)\n", path, err)
		return nil, err
	}
	//fmt.Printf("Fileread %s (%s)\n", path, err)

	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening file for read: %v", err)
		return nil, err
	}
	return file, nil
}

func (h *VirtualSftpHandler) Filewrite(r *sftp.Request) (io.WriterAt, error) {
	path, err := h.resolve(r.Filepath)
	if err != nil {
		fmt.Printf("Filewrite: Resolve failed: %s (%s)\n", path, err)
		return nil, err
	}
	//fmt.Printf("Filewrite %s (%s)\n", path, err)

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, h.permissions)
	if err != nil {
		fmt.Printf("Error opening file for write: %v", err)
		return nil, err
	}

	h.newFiles = append(h.newFiles, path)

	return file, nil
}

func (h *VirtualSftpHandler) Filecmd(r *sftp.Request) error {
	path, err := h.resolve(r.Filepath)
	if err != nil {
		fmt.Printf("Filecmd: Resolve failed: %s (%s)\n", path, err)
		return err
	}
	//fmt.Printf("Filecmd %s %s (%s)\n", r.Method, path, err)

	switch r.Method {
	case "Mkdir":
		if err = os.MkdirAll(path, h.permissions); err != nil {
			fmt.Printf("Failed to create directory %s: %s\n", path, err)
			return err
		}
		return nil
	case "Remove":
		if err := os.Remove(path); err != nil {
			fmt.Printf("Failed to remove file %s: %s\n", path, err)
			return err
		}
		return nil
	case "Rmdir":
		if err := os.RemoveAll(path); err != nil {
			fmt.Printf("Failed to remove directory %s: %s\n", path, err)
			return err
		}
		return nil
	case "Rename":
		newPath, err := h.resolve(r.Target)
		if err != nil {
			fmt.Printf("Failed to rename file %s to %s: %s\n", path, newPath, err)
			return err
		}
		if err := os.Rename(path, newPath); err != nil {
			fmt.Printf("Failed to rename file %s to %s: %s\n", path, newPath, err)
			return err
		}

		h.newFiles = append(h.newFiles, newPath)

		return nil
	case "Setstat":
		attr := r.Attributes()
		if attr == nil {
			fmt.Printf("Setstat failed, no attributed provided for %s: %s\n", path, err)
			return err
		}
		if attr.Mode != 0 {
			fileMode := os.FileMode(attr.Mode & 0o777)
			if err := os.Chmod(path, fileMode); err != nil {
				fmt.Printf("Setstat: Chmod failed for %s: %s\n", path, err)
			}
		}
		if attr.Atime != 0 {
			if err := os.Chtimes(path, time.Unix(int64(attr.Atime), 0), time.Time{}); err != nil {
				fmt.Printf("Setstat: Chtimes (Atime) failed for %s: %s\n", path, err)
			}
		}
		if attr.Mtime != 0 {
			if err := os.Chtimes(path, time.Time{}, time.Unix(int64(attr.Mtime), 0)); err != nil {
				fmt.Printf("Setstat: Chtimes (Mtime) failed for %s: %s\n", path, err)
			}
		}
		if attr.UID != 0 || attr.GID != 0 {
			if runtime.GOOS != "windows" {
				if err := os.Chown(path, int(attr.UID), int(attr.GID)); err != nil {
					fmt.Printf("Setstat: Chown failed for %s: %s\n", path, err)
				}
			} else {
				fmt.Printf("Setstat: Chown not available for %s: %s\n", path, err)
			}
		}
		if attr.Size > 0 {
			if stat, err := os.Stat(path); err == nil {
				if stat.Mode().IsRegular() {
					if err := os.Truncate(path, int64(attr.Size)); err != nil {
						fmt.Printf("Setstat: Set size failed for %s: %s\n", path, err)
					}
				} else {
					fmt.Printf("Setstat: Set size not possible for non regular %s: %s\n", path, err)
				}
			} else {
				fmt.Printf("Setstat: Set size not possible for %s: %s\n", path, err)
			}
		}
		return nil
	}

	fmt.Printf("Ignored: Filecmd %s\n", r.Method)
	return os.ErrInvalid
}

func (h *VirtualSftpHandler) Filelist(r *sftp.Request) (sftp.ListerAt, error) {

	var err error

	path, err := h.resolve(r.Filepath)
	if err != nil {
		fmt.Printf("Filelist: Resolve failed: %s (%s)\n", path, err)
		return nil, err
	}
	//fmt.Printf("Filelist %s %s (%s)\n", path, r.Method, err)

	var lister sftp.ListerAt

	switch r.Method {
	case "List":
		lister, err = NewFileLister(path)

	case "Stat", "Lstat", "Readlink":
		lister, err = NewStatFileLister(path)

	default:
		fmt.Printf("Unsupported Filelist Method: Path=%s, Method=%s (%s)\n", path, r.Method, err)
		return nil, errors.New("unsupported")
	}

	if err != nil {
		fmt.Printf("Filelist error: %s (%s)\n", path, err)
		return nil, err
	}
	return lister, nil

}

func (h *VirtualSftpHandler) resolve(path string) (string, error) {
	return resolve(h.rootDirectory, path)
}

type FileLister struct {
	list []os.FileInfo
}

func NewFileLister(path string) (*FileLister, error) {

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil, os.ErrNotExist
	}

	lister := &FileLister{}

	fis, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer util.CloseOrLog(fis)

	lister.list, err = fis.Readdir(-1)
	if err != nil {
		fmt.Printf("Error reading directory %s: %s\n", path, err)
		return nil, err
	}

	return lister, nil
}

func NewStatFileLister(path string) (*FileLister, error) {

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil, os.ErrNotExist
	}

	lister := &FileLister{}

	f, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}

	lister.list = []os.FileInfo{f}
	return lister, nil
}

func (l *FileLister) ListAt(ls []os.FileInfo, offset int64) (int, error) {
	// For empty directories, return (0, nil) on the first call so some clients
	// (e.g., WinSCP) don't interpret immediate EOF as an error.
	//fmt.Printf("ListAt: offset=%d length=%d\n", offset, len(l.list))
	//if len(l.list) == 0 && offset == 0 {
	//	return 0, nil
	//}
	if offset >= int64(len(l.list)) {
		return 0, io.EOF
	}
	n := copy(ls, l.list[offset:])
	if n < len(ls) {
		return n, io.EOF
	}
	return n, nil
}
