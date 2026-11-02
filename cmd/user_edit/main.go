package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/users"
	"github.com/c8121/go_vue_spa/auth/cli"
)

func main() {

	username := flag.String("u", "", "Username")
	addUser := flag.Bool("a", false, "Add new user")
	pwdCost := flag.Int("cost", 10, "Password hash cost")

	config.LoadDefault()

	err := users.InitLocalFileAppAuthentication(*pwdCost)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		fmt.Println(err)
		return
	}

	auth.CliLocalFileUserEdit(*username, *addUser)
}
