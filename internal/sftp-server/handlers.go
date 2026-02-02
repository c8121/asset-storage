package sftp_server

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/sftp"
)

var (
	ErrorInvalidPathRequest = errors.New("Invalid path request")
)

func NewVirtualSftpHandler(rootDirectory string) sftp.Handlers {
	handler := &VirtualSftpHandler{
		rootDirectory: rootDirectory,
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
}

func (h *VirtualSftpHandler) Fileread(r *sftp.Request) (io.ReaderAt, error) {
	path, err := h.resolve(r.Filepath)
	fmt.Printf("Fileread %s (%s)\n", path, err)
	return nil, os.ErrPermission
}

func (h *VirtualSftpHandler) Filewrite(r *sftp.Request) (io.WriterAt, error) {
	path, err := h.resolve(r.Filepath)
	fmt.Printf("Filewrite %s (%s)\n", path, err)
	return nil, os.ErrPermission
}

func (h *VirtualSftpHandler) Filecmd(r *sftp.Request) error {
	path, err := h.resolve(r.Filepath)
	fmt.Printf("Filecmd %s (%s)\n", path, err)
	switch r.Method {
	case "Mkdir":

		break
	}
	return os.ErrPermission
}

func (h *VirtualSftpHandler) Filelist(r *sftp.Request) (sftp.ListerAt, error) {
	path, err := h.resolve(r.Filepath)
	fmt.Printf("Filelist %s (%s)\n", path, err)
	return nil, os.ErrPermission
}

func (h *VirtualSftpHandler) ListAt(ls []os.FileInfo, offset int64) (int, error) {
	fmt.Printf("ListAt %d\n", offset)
	return 0, io.EOF
}

func (h *VirtualSftpHandler) Lstat(r *sftp.Request) (sftp.ListerAt, error) {
	path, err := h.resolve(r.Filepath)
	fmt.Printf("Lstat %s (%s)\n", path, err)
	return nil, os.ErrPermission
}

func (h *VirtualSftpHandler) Stat(r *sftp.Request) (sftp.ListerAt, error) {
	path, err := h.resolve(r.Filepath)
	fmt.Printf("Stat %s (%s)\n", path, err)
	return nil, os.ErrPermission
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
