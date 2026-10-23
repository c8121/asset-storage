package ssh_server

import (
	"github.com/c8121/asset-storage/internal/access_control"
	"github.com/c8121/asset-storage/internal/users"
)

func GetAccessContext(h SshHandler) *access_control.AccessContext {
	user := &users.AuthenticatedUser{Username: h.GetUsername()}
	return &access_control.AccessContext{User: user}
}
