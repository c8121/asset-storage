package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/c8121/asset-storage/internal/users"
	"github.com/c8121/asset-storage/shared/config"
	auth "github.com/c8121/go_vue_spa/auth/cli"
)

func main() {

	data := flag.String("data", "", "Data directory for storage, meta-data, db...")
	username := flag.String("u", "", "Username")
	addUser := flag.Bool("a", false, "Add new user")
	pwdCost := flag.Int("cost", 10, "Password hash cost")
	flag.Parse()

	dir, err := config.ResolveDataDir(*data)
	if err != nil {
		fmt.Println(err)
		return
	}

	config.AppConfig = config.NewConfigFromPath(dir)

	err = users.InitLocalFileAppAuthentication(*pwdCost)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		fmt.Println(err)
		return
	}

	auth.CliLocalFileUserEdit(*username, *addUser)
}
