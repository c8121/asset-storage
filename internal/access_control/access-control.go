package access_control

import (
	"errors"

	"github.com/c8121/go_vue_spa/auth"
)

/*type AccessLevel int

const (
	AccessLevelDeny  = 1
	AccessLevelRead  = 2
	AccessLevelWrite = 4
)

type AccessControlEntry struct {
	Username string
	//Groupname string
	Level AccessLevel
}*/

var (
	ErrAccessDenied = errors.New("access denied")
)

type AccessContext struct {
	User *auth.AuthenticatedUser
}

type AccessController interface {
	CanRead(assetHash string, ctx *AccessContext) error
	CanWritePath(path string, ctx *AccessContext) error
}

var (
	accessController AccessController
)

func SetAccessController(ctrl AccessController) {
	accessController = ctrl
}

func GetAccessController() AccessController {
	return accessController
}
