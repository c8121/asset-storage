package ssh_server

import (
	"io"
	"log"
	"log/slog"
	"net"
	"os"
	"os/exec"

	"github.com/c8121/go_vue_spa/util"
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
	slog.Info("SSH Server ready", "address", listener.Addr())

	for {
		clientConnection, err := listener.Accept()
		if err != nil {
			log.Fatal("failed to accept incoming connection", err)
		}

		go func(conn net.Conn) {
			defer util.CloseOrLog(conn)

			//Handshake
			sshConnection, chans, reqs, err := ssh.NewServerConn(clientConnection, serverConfig)
			if err != nil {
				slog.Warn("handshake failed", "error", err)
				return
			}
			slog.Info("SSH connection established", "remote", sshConnection.RemoteAddr(), "local", sshConnection.ClientVersion())

			// Discard global requests
			go ssh.DiscardRequests(reqs)

			// Service the incoming channel.
			for newChannel := range chans {
				if newChannel.ChannelType() != "session" {
					if err := newChannel.Reject(ssh.UnknownChannelType, "unknown channel type"); err != nil {
						slog.Warn("failed to reject channel", "error", err)
					}
					slog.Warn("Unknown channel type", "type", newChannel.ChannelType())
					continue
				}
				channel, requests, err := newChannel.Accept()
				if err != nil {
					slog.Warn("Could not accept channel", "error", err)
					continue
				}
				slog.Info("Channel accepted", "remote", sshConnection.RemoteAddr(), "type", newChannel.ChannelType())

				go func(in <-chan *ssh.Request) {
					for req := range in {

						switch req.Type {
						case "subsystem":
							if len(req.Payload) >= 4 && string(req.Payload[4:]) == "sftp" {
								if config.SftpHandlerCreator != nil {
									if err := req.Reply(true, nil); err != nil {
										slog.Warn("Could not reply to client", "error", err)
										return
									}
									if err := executeSFTP(config, *sshConnection, channel); err != nil {
										slog.Warn("Subsystem failed", "error", err)
										return
									}

									continue
								}
							}
						case "exec":
							if len(req.Payload) >= 3 && string(req.Payload[4:9]) == "rsync" {
								if config.RsyncHandlerCreator != nil {
									if err := req.Reply(true, nil); err != nil {
										slog.Warn("Could not reply to client", "error", err)
										return
									}
									if err := executeRsync(config, *sshConnection, channel, string(req.Payload[4:])); err != nil {
										slog.Warn("Execute command failed", "error", err)
										return
									}

									continue
								}
							}
						}

						// request wasn't handled above, send request-not-ok
						if err := req.Reply(false, nil); err != nil {
							slog.Warn("Could not reply to client", "error", err)
							return
						}
						slog.Warn("Decline request", "remote", sshConnection.RemoteAddr(), "type", req.Type)
					}
				}(requests)

			}

		}(clientConnection)
	}
}

// executeRsync executes rsync-command (given in payload)
func executeRsync(config *SshServerConfig, sshConnection ssh.ServerConn, channel ssh.Channel, payload string) error {

	slog.Info("Start RSYNC server", "remote", sshConnection.RemoteAddr())

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

	go func() {
		if _, err := io.Copy(stdin, channel); err != nil {
			slog.Warn("Error copying stdin to channel", "error", err)
		}
	}()
	go func() {
		if _, err := io.Copy(channel, stdout); err != nil {
			slog.Warn("Error copying stdout to channel", "error", err)
		}
	}()

	err = cmd.Run()
	if err != nil {
		util.CloseOrLog(channel)
		return err
	}
	if _, err = channel.SendRequest("exit-status", false, []byte{0, 0, 0, 0}); err != nil {
		slog.Warn("Rsync send command exit-status failed", "error", err)
	}
	util.CloseOrLog(channel)

	slog.Info("Command completed", "remote", sshConnection.RemoteAddr(), "cmd", cmdName, "args", cmdArgs)

	// Add new files to storage
	go AddFilesToArchive(handler)

	return nil
}

// executeSFTP starts the SFTP-Subsystem
func executeSFTP(config *SshServerConfig, sshConnection ssh.ServerConn, channel ssh.Channel) error {

	slog.Info("Start SFTP server", "remote", sshConnection.RemoteAddr())

	handler, err := config.SftpHandlerCreator(sshConnection.Permissions.Extensions["username"])
	if err != nil {
		return err
	}

	server := sftp.NewRequestServer(channel, handler.GetHandlers())
	if err := server.Serve(); err != nil {
		if err != io.EOF {
			slog.Warn("SFTP server completed with error", "error", err, "remote", sshConnection.RemoteAddr())
		}
	}
	if err := server.Close(); err != nil && err != io.EOF {
		slog.Warn("SFTP server close error", "error", err)
	}
	slog.Info("SFTP client exited session", "remote", sshConnection.RemoteAddr())

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
