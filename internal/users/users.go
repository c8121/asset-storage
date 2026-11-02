package users

import (
	"os"
	"path/filepath"

	"github.com/c8121/asset-storage/internal/util"
	"github.com/c8121/asset-storage/shared/config"
	"github.com/c8121/go_vue_spa/auth"
)

var (
	UserFile            = "users"
	DefaultPasswordCost = 10
)

func InitAppAuthentication(authenticator auth.Authenticator) {
	auth.SetAuthenticator(authenticator)
}

func InitLocalFileAppAuthentication(pwdCost int) error {

	util.CreateDirIfNotExists(config.AppConfig.AssetStorageConfigDir, auth.DefaultFilePermissions)

	usersFile := filepath.Join(config.AppConfig.AssetStorageConfigDir, UserFile)

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
	usersFile := filepath.Join(config.AppConfig.AssetStorageConfigDir, UserFile)
	_, err := os.Stat(usersFile)
	return err == nil
}
