package users

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/c8121/asset-storage/internal/config"
	"github.com/c8121/asset-storage/internal/util"
	"golang.org/x/crypto/bcrypt"
)

const (
	UserFile        = "users"
	FilePermissions = 0700
)

var (
	users map[string]string = nil
)

// CreateDirectories creates required directories
func CreateDirectories() {
	util.CreateDirIfNotExists(config.AssetStorageConfigDir, FilePermissions)
}

func UserExists(name string) bool {
	readUsers()

	_, ok := users[name]
	return ok
}

func Authenticate(username string, password []byte) error {

	if !UserExists(username) {
		return errors.New("No such user")
	}

	hash, _ := users[username]

	return bcrypt.CompareHashAndPassword([]byte(hash), password)
}

func AddOrUpdateUser(name string, pwdHash string) {
	readUsers()

	_, ok := users[name]
	if ok {
		fmt.Printf("Update user: %s\n", name)
	} else {
		fmt.Printf("Add user: %s\n", name)
	}
	users[name] = pwdHash
}

func readUsers() {

	if users != nil {
		return
	}
	users = make(map[string]string)

	file := filepath.Join(config.AssetStorageConfigDir, UserFile)
	stat, err := os.Stat(file)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			fmt.Printf("Cannot read users file")
		}
		return //File doesnt exist yet, nothing to read
	} else if !stat.Mode().IsRegular() {
		panic("Cannot read users file: Not a file\n")
	}

	f, err := os.Open(file)
	if err != nil {
		fmt.Printf("Cannot open users file: %s\n", err)
		return
	}
	defer util.CloseOrLog(f)

	reader := bufio.NewReaderSize(f, 1000)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Failed to read line: %s\n", err)
			return
		}

		if len(line) == 0 || line[0] == '#' {
			continue
		}

		s := string(line)
		p := strings.Index(s, "\t")
		if p < 1 {
			fmt.Printf("Invalid line: %s\n", s)
			continue
		}

		users[s[0:p]] = s[p+1:]
	}

}

func SaveUsers() {

	if users == nil {
		return
	}

	file := filepath.Join(config.AssetStorageConfigDir, UserFile)
	stat, err := os.Stat(file)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		panic(err)
	} else if stat != nil && !stat.Mode().IsRegular() {
		panic("Cannot read users file: Not a file\n")
	}

	f, err := os.Create(file)
	if err != nil {
		panic(err)
	}
	defer util.CloseOrLog(f)

	w := bufio.NewWriter(f)
	for name, hash := range users {
		_, err := w.WriteString(name + "\t" + hash + "\n")
		if err != nil {
			panic(err)
		}
	}
	w.Flush()

	fmt.Printf("Updated user file: %s\n", file)
}

func CreatePasswordHash(password []byte, pwdCost int) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(password, pwdCost)
	if err != nil {
		fmt.Printf("Cannot create hash: %s\n", err)
		return "", err
	}

	return string(hash), nil
}
