package sftp_server

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/c8121/asset-storage/internal/util"
	"github.com/pkg/sftp"
)

var (
	ErrorInvalidPathRequest = errors.New("Invalid path request")
)

func NewVirtualSftpHandler(rootDirectory string) sftp.Handlers {
	handler := &VirtualSftpHandler{
		rootDirectory: rootDirectory,
		permissions:   0700,
	}
	return sftp.Handlers{
		FileGet:  handler,
		FilePut:  handler,
		FileCmd:  handler,
		FileList: handler,
	}
}

type VirtualSftpHandler struct {
	rootDirectory string
	permissions   os.FileMode
}

func (h *VirtualSftpHandler) Fileread(r *sftp.Request) (io.ReaderAt, error) {
	path, err := h.resolve(r.Filepath)
	fmt.Printf("Fileread %s (%s)\n", path, err)

	return nil, os.ErrPermission
}

func (h *VirtualSftpHandler) Filewrite(r *sftp.Request) (io.WriterAt, error) {
	path, err := h.resolve(r.Filepath)
	fmt.Printf("Filewrite %s (%s)\n", path, err)

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, h.permissions)
	if err != nil {
		fmt.Printf("Error opening file for write: %v", err)
		return nil, err
	}
	return file, nil
}

func (h *VirtualSftpHandler) Filecmd(r *sftp.Request) error {
	path, err := h.resolve(r.Filepath)
	fmt.Printf("Filecmd %s (%s)\n", path, err)

	switch r.Method {
	case "Mkdir":
		if err = os.MkdirAll(path, h.permissions); err != nil {
			fmt.Printf("Failed to create directory %s: %s\n", path, err)
			return err
		}
		return nil
	}

	fmt.Printf("Ignored: Filecmd %s\n", r.Method)
	return os.ErrInvalid
}

func (h *VirtualSftpHandler) Filelist(r *sftp.Request) (sftp.ListerAt, error) {

	var err error

	path, err := h.resolve(r.Filepath)
	fmt.Printf("Filelist %s %s (%s)\n", path, r.Method, err)
	if err != nil {
		fmt.Printf("Filelist: Resolve failed: %s (%s)\n", path, err)
		return nil, err
	}

	var lister sftp.ListerAt

	switch r.Method {
	case "List":
		lister, err = NewFileLister(path)
		break

	case "Stat":
		lister, err = NewStatFileLister(path)
		break
	default:
		return nil, errors.New("unsupported")
	}

	if err != nil {
		fmt.Printf("Filelist: %s (%s)\n", path, err)
		return nil, err
	}
	return lister, nil

}

func (h *VirtualSftpHandler) resolve(path string) (string, error) {

	fmt.Printf(" - Path requested: '%s'\n", path)
	if strings.Contains(path, ":") {
		return "", ErrorInvalidPathRequest
	}

	resolved := path
	if resolved == string(filepath.Separator) || resolved == "/" {
		resolved = ""
	}

	resolved = filepath.FromSlash(resolved)
	resolved = strings.TrimPrefix(resolved, string(filepath.Separator))
	resolved = strings.TrimPrefix(resolved, "/")
	resolved = filepath.Join(h.rootDirectory, resolved)

	fmt.Printf(" - Path resolved: '%s'\n", resolved)
	return resolved, nil
}

type FileLister struct {
	list []os.FileInfo
}

func NewFileLister(path string) (*FileLister, error) {
	lister := &FileLister{}

	fis, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer util.CloseOrLog(fis)

	lister.list, err = fis.Readdir(-1)
	if err != nil {
		return nil, err
	}

	fmt.Printf(" - Path lister opened: '%s' (%d entries)\n", path, len(lister.list))
	return lister, nil
}

func NewStatFileLister(path string) (*FileLister, error) {
	lister := &FileLister{}

	f, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}

	lister.list = []os.FileInfo{f}

	fmt.Printf(" - Path lister opened: '%s' (%d entries)\n", path, len(lister.list))
	return lister, nil
}

func (l *FileLister) ListAt(ls []os.FileInfo, offset int64) (int, error) {
	// For empty directories, return (0, nil) on the first call so some clients
	// (e.g., WinSCP) don't interpret immediate EOF as an error.
	fmt.Printf(" - ListAt offset: %d\n", offset)
	if len(l.list) == 0 && offset == 0 {
		return 0, nil
	}
	if offset >= int64(len(l.list)) {
		return 0, io.EOF
	}
	n := copy(ls, l.list[offset:])
	if n < len(ls) {
		return n, io.EOF
	}
	return n, nil
}
