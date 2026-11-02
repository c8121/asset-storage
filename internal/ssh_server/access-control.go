package ssh_server

import (
	"github.com/c8121/asset-storage/internal/access_control"
	"github.com/c8121/go_vue_spa/auth"
)

func GetAccessContext(h SshHandler) *access_control.AccessContext {
	user := &auth.AuthenticatedUser{Username: h.GetUsername()}
	return &access_control.AccessContext{User: user}
}
