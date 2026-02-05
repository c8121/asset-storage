package ssh_server

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SshServerConfig struct {
	ListenAddress       string
	HostKeyFile         string
	PasswordCallback    func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error)
	SftpHandlerCreator  func(username string) (*VirtualSftpHandler, error)
	RsyncHandlerCreator func(username string) (*VirtualRsyncHandler, error)
}

func RunSftpServer(config *SshServerConfig) {

	serverConfig := &ssh.ServerConfig{
		PasswordCallback: config.PasswordCallback,
	}

	serverConfig.AddHostKey(loadHostKey(config.HostKeyFile))

	listener, err := net.Listen("tcp", config.ListenAddress)
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
			sshConnection, chans, reqs, err := ssh.NewServerConn(clientConnection, serverConfig)
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

						switch req.Type {
						case "subsystem":
							if len(req.Payload) >= 4 && string(req.Payload[4:]) == "sftp" {
								if config.SftpHandlerCreator != nil {
									if err := req.Reply(true, nil); err != nil {
										fmt.Printf("Could not reply to client: %s\n", err)
										return
									}
									if err := executeSFTP(config, *sshConnection, channel); err != nil {
										fmt.Printf("Subsystem failed: %s\n", err)
										return
									}

									continue
								}
							}
						case "exec":
							if len(req.Payload) >= 3 && string(req.Payload[4:9]) == "rsync" {
								if config.RsyncHandlerCreator != nil {
									if err := req.Reply(true, nil); err != nil {
										fmt.Printf("Could not reply to client: %s\n", err)
										return
									}
									if err := executeRsync(config, *sshConnection, channel, string(req.Payload[4:])); err != nil {
										fmt.Printf("Execute command failed: %s\n", err)
										return
									}

									continue
								}
							}
						}

						// request wasn't handled above, send request-not-ok
						if err := req.Reply(false, nil); err != nil {
							fmt.Printf("Could not reply to client: %s\n", err)
							return
						}
						fmt.Printf("Decline request type for %s: %s\n", sshConnection.RemoteAddr(), req.Type)

					}
				}(requests)

			}

		}(clientConnection)
	}
}

// executeRsync executes rsync-command (given in payload)
func executeRsync(config *SshServerConfig, sshConnection ssh.ServerConn, channel ssh.Channel, payload string) error {

	fmt.Printf("Start RSYNC server for %s.\n", sshConnection.RemoteAddr())

	handler, err := config.RsyncHandlerCreator(sshConnection.Permissions.Extensions["username"])
	if err != nil {
		return err
	}

	cmdName, cmdArgs, err := handler.parseCommand(payload)
	if err != nil {
		return err
	}

	cmd := exec.Command(cmdName, cmdArgs...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	go io.Copy(stdin, channel)
	go io.Copy(channel, stdout)

	cmd.Run()
	channel.SendRequest("exit-status", false, []byte{0, 0, 0, 0})
	channel.Close()

	fmt.Printf("Command completed for %s: %s %v.\n", sshConnection.RemoteAddr(), cmdName, cmdArgs)

	// Add new files to storage
	go AddFilesToArchive(handler)

	return nil
}

// executeSFTP starts the SFTP-Subsystem
func executeSFTP(config *SshServerConfig, sshConnection ssh.ServerConn, channel ssh.Channel) error {

	fmt.Printf("Start SFTP server for %s.\n", sshConnection.RemoteAddr())

	handler, err := config.SftpHandlerCreator(sshConnection.Permissions.Extensions["username"])
	if err != nil {
		return err
	}

	server := sftp.NewRequestServer(channel, handler.GetHandlers())
	if err := server.Serve(); err != nil {
		if err != io.EOF {
			fmt.Printf("SFTP server completed with error: %s, %s\n", err, sshConnection.RemoteAddr())
		}
	}
	if err := server.Close(); err != nil && err != io.EOF {
		fmt.Printf("SFTP server close error: %s\n", err)
	}
	fmt.Printf("SFTP client exited session: %s.\n", sshConnection.RemoteAddr())

	// Add new files to storage
	go AddFilesToArchive(handler)

	return nil
}

// loadHostKey reads & parses the key-file
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
