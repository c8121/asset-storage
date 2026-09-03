package util

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func CliConfirm(s string) bool {
	r := bufio.NewReader(os.Stdin)

	fmt.Printf("%s [y/n]? ", s)

	res, err := r.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	return strings.ToLower(strings.TrimSpace(res))[0] == 'y'
}
