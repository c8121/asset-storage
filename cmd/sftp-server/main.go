package main

import (
	"flag"
	"fmt"

	sftp_server "github.com/c8121/asset-storage/internal/sftp-server"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

var (
	ListenAddress = "127.0.0.1:2022"
	HostKeyFile   = "id_rsa"
)

func main() {

	flag.StringVar(&ListenAddress, "listen", ListenAddress, "Listen Address (ip:port)")
	flag.StringVar(&HostKeyFile, "host-key", HostKeyFile, "SSH-Host Key-File")
	flag.Parse()

	passwordCallback := func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
		fmt.Printf("Login from %s: %s\n", c.RemoteAddr(), c.User())
		/*if c.User() == "test" && string(pass) == "test" {
			return nil, nil
		}
		return nil, fmt.Errorf("password rejected for %q", c.User())*/
		return nil, nil
	}

	handlerCreator := func() sftp.Handlers {
		return sftp_server.NewVirtualSftpHandler("/tmp")
	}

	sftp_server.RunSftpServer(ListenAddress, HostKeyFile, passwordCallback, handlerCreator)
}
