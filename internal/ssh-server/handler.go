package ssh_server

import "fmt"

type SshHandler interface {
	GetUsername() string
	GetNewFiles() []string
}

func AddFilesToArchive(h SshHandler) {

	files := h.GetNewFiles()
	for _, file := range files {
		fmt.Printf("TODO: Add file from user %s: %s\n", h.GetUsername(), file)
	}

}
