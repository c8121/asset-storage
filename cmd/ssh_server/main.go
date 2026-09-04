package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/metadata"
	"github.com/c8121/asset-storage/internal/metadata_sqlite"
	"github.com/c8121/asset-storage/internal/ssh_server"
	"github.com/c8121/asset-storage/internal/storage"
	"github.com/c8121/asset-storage/internal/users"
	"golang.org/x/crypto/ssh"
)

var (
	hostKeyFile = flag.String("host-key", "id_rsa", "SSH-Host Key-File")
)

func main() {

	config.LoadDefault()
	storage.CreateDirectories()
	metadata.CreateDirectories()

	metadata_sqlite.Open()
	defer metadata_sqlite.Close()

	sshServerConfig := &ssh_server.SshServerConfig{

		ListenAddress: config.ListenAddress,

		HostKeyFile: *hostKeyFile,

		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {

			fmt.Printf("Login attempt from %s: %s\n", c.RemoteAddr(), c.User())
			if len(c.User()) == 0 {
				return nil, fmt.Errorf("user rejected for %q", c.User())
			} else if err := users.Authenticate(c.User(), pass); err != nil {
				return nil, fmt.Errorf("password rejected for %q", c.User())
			}

			fmt.Printf("Login successful %s: %s\n", c.RemoteAddr(), c.User())
			perms := &ssh.Permissions{Extensions: map[string]string{"username": c.User()}}
			return perms, nil
		},

		SftpHandlerCreator: func(username string) (*ssh_server.VirtualSftpHandler, error) {
			root, err := createVirtualRootDir(username)
			if err != nil {
				return nil, err
			}
			return ssh_server.NewVirtualSftpHandler(root, username), nil
		},

		RsyncHandlerCreator: func(username string) (*ssh_server.VirtualRsyncHandler, error) {
			root, err := createVirtualRootDir(username)
			if err != nil {
				return nil, err
			}
			return ssh_server.NewVirtualRsyncHandler(root, username), nil
		},
	}

	ssh_server.RunSftpServer(sshServerConfig)
}

// create a home-/root-directory for the given user.
func createVirtualRootDir(username string) (string, error) {
	root := filepath.Join(config.AssetStorageTempDir, "virtual-users", username)
	if err := os.MkdirAll(root, storage.FilePermissions); err != nil {
		return "", err
	}
	return root, nil
}
