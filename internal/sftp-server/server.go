package sftp_server

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func RunSftpServer(
	listenAddress string, hostKeyFile string,
	passwordCallback func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error),
	handlerCreator func() sftp.Handlers) {

	config := &ssh.ServerConfig{
		PasswordCallback: passwordCallback,
	}

	config.AddHostKey(loadHostKey(hostKeyFile))

	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatal("failed to listen for connection", err)
	}
	fmt.Printf("Listening on %v\n", listener.Addr())

	for {
		clientConnection, err := listener.Accept()
		if err != nil {
			log.Fatal("failed to accept incoming connection", err)
		}

		go func(conn net.Conn) {
			defer conn.Close()

			//Handshake
			sshConnection, chans, reqs, err := ssh.NewServerConn(clientConnection, config)
			if err != nil {
				fmt.Printf("handshake failed: %s\n", err)
				return
			}
			fmt.Printf("SSH connection established: %s (%s)\n", sshConnection.RemoteAddr(), sshConnection.ClientVersion())

			// Discard global requests
			go ssh.DiscardRequests(reqs)

			// Service the incoming Channel channel.
			for newChannel := range chans {
				if newChannel.ChannelType() != "session" {
					if err := newChannel.Reject(ssh.UnknownChannelType, "unknown channel type"); err != nil {
						fmt.Printf("failed to reject channel %s\n", err)
					}
					fmt.Printf("Unknown channel type: %s\n", newChannel.ChannelType())
					continue
				}
				channel, requests, err := newChannel.Accept()
				if err != nil {
					fmt.Printf("Could not accept channel: %s\n", err)
					continue
				}
				fmt.Printf("Channel accepted for %s: %s\n", sshConnection.RemoteAddr(), newChannel.ChannelType())

				go func(in <-chan *ssh.Request) {
					for req := range in {
						if req.Type == "subsystem" && len(req.Payload) >= 4 && string(req.Payload[4:]) == "sftp" {

							err := req.Reply(true, nil)
							if err != nil {
								fmt.Printf("Could not reply to client: %s\n", err)
								return
							}

							fmt.Printf("Start SFTP server for %s.\n", sshConnection.RemoteAddr())

							root := handlerCreator()
							server := sftp.NewRequestServer(channel, root)
							if err := server.Serve(); err != nil {
								if err != io.EOF {
									fmt.Printf("SFTP server completed with error: %s, %s\n", err, sshConnection.RemoteAddr())
								}
							}
							if err := server.Close(); err != nil && err != io.EOF {
								fmt.Printf("SFTP server close error: %s\n", err)
							}
							fmt.Printf("SFTP client exited session: %s.\n", sshConnection.RemoteAddr())

						} else {
							if err := req.Reply(false, nil); err != nil {
								fmt.Printf("Could not reply to client: %s\n", err)
							}
							fmt.Printf("Decline request type for %s: %s\n", sshConnection.RemoteAddr(), req.Type)
						}
					}
				}(requests)

			}

		}(clientConnection)
	}
}

func loadHostKey(hostKeyFile string) ssh.Signer {
	privateBytes, err := os.ReadFile(hostKeyFile)
	if err != nil {
		log.Fatal("Failed to load private key", err)
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatal("Failed to parse private key", err)
	}

	return private
}
