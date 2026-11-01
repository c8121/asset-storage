package auth

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"syscall"

	"github.com/c8121/asset-storage/internal/users"
	"github.com/c8121/asset-storage/shared/auth"
	"golang.org/x/term"
)

func CliLocalFileUserEdit(username string, addUser bool, pwdCost int) {

	flag.Parse()

	err := users.InitLocalFileAppAuthentication(pwdCost)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		fmt.Println(err)
		return
	}

	authenticator, ok := users.AppAuthenticator.(*auth.LocalFileAuthenticator)
	if !ok {
		fmt.Println("Init error: Not a local file authenticator")
		return
	}

	if len(username) == 0 {
		fmt.Printf("Missing: Username (-u)\n")
		return
	}

	if addUser {
		fmt.Printf("Add user: %s\n", username)
		if authenticator.UserExists(username) {
			fmt.Printf("User already exists: %s\n", username)
			return
		}
	} else {
		fmt.Printf("User edit: %s\n", username)
		if !authenticator.UserExists(username) {
			fmt.Printf("User not found: %s\n", username)
			return
		}
	}

	password, err := readPassword()
	if err != nil {
		return
	}
	err = authenticator.AddOrUpdateUser(username, password)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = authenticator.SaveUsers()
	if err != nil {
		fmt.Println(err)
	}
}

func readPassword() ([]byte, error) {
	fmt.Printf("Enter password: ")
	password, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Printf("Cannot read password input: %v\n", err)
		return nil, err
	}
	fmt.Println()

	fmt.Printf("Confirm password: ")
	checkPassword, _ := term.ReadPassword(int(syscall.Stdin))
	if !bytes.Equal(password, checkPassword) {
		fmt.Printf("Password missmatch\n")
		return nil, errors.New("password mismatch")
	}
	fmt.Println()

	return password, nil
}
