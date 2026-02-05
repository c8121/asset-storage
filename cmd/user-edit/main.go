package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"syscall"

	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/users"
	"golang.org/x/term"
)

func main() {

	username := flag.String("u", "", "Username")
	addUser := flag.Bool("a", false, "Add new user")
	pwdCost := flag.Int("cost", 10, "Password hash cost")
	config.LoadDefault()

	users.CreateDirectories()

	if len(*username) == 0 {
		fmt.Printf("Missing: Usernamen (-u)")
		return
	}

	if *addUser {
		fmt.Printf("Add user: %s\n", *username)
		if users.UserExists(*username) {
			fmt.Printf("User already exists: %s\n", *username)
			return
		}
	} else {
		fmt.Printf("User edit: %s\n", *username)
		if !users.UserExists(*username) {
			fmt.Printf("User not found: %s\n", *username)
			return
		}
	}

	hash, err := readPasswordHash(*pwdCost)
	if err != nil {
		return
	}
	users.AddOrUpdateUser(*username, hash)

	users.SaveUsers()
}

func readPasswordHash(pwdCost int) (string, error) {
	fmt.Printf("Enter password: ")
	password, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Printf("Cannot read password input: %v\n", err)
		return "", err
	}
	fmt.Println()

	fmt.Printf("Confirm password: ")
	checkPassword, _ := term.ReadPassword(int(syscall.Stdin))
	if !bytes.Equal(password, checkPassword) {
		fmt.Printf("Password missmatch\n")
		return "", errors.New("Password mismatch")
	}
	fmt.Println()

	return users.CreatePasswordHash(password, pwdCost)
}
