package users

import (
	"os"
	"path/filepath"

	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/util"
	"github.com/c8121/asset-storage/shared/auth"
)

var (
	UserFile                               = "users"
	AppAuthenticator    auth.Authenticator = nil
	DefaultPasswordCost                    = 10
)

func InitAppAuthentication(authenticator auth.Authenticator) {
	AppAuthenticator = authenticator
}

func InitLocalFileAppAuthentication(pwdCost int) error {

	util.CreateDirIfNotExists(config.AssetStorageConfigDir, auth.DefaultFilePermissions)

	usersFile := filepath.Join(config.AssetStorageConfigDir, UserFile)

	authenticator, err := auth.NewLocalFileAuthenticator(usersFile, pwdCost)
	if err != nil {
		return err
	}

	err = authenticator.ReadUsers()
	if err != nil {
		return err
	}

	InitAppAuthentication(authenticator)
	return nil
}

func UsersFileExists() bool {
	usersFile := filepath.Join(config.AssetStorageConfigDir, UserFile)
	_, err := os.Stat(usersFile)
	return err == nil
}
