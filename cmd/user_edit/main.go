package main

import (
	"flag"

	"github.com/c8121/asset-storage/internal/config"
	auth "github.com/c8121/asset-storage/shared/auth/cli"
)

func main() {

	username := flag.String("u", "", "Username")
	addUser := flag.Bool("a", false, "Add new user")
	pwdCost := flag.Int("cost", 10, "Password hash cost")

	config.LoadDefault()

	auth.CliLocalFileUserEdit(*username, *addUser, *pwdCost)
}
