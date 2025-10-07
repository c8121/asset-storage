package thumbnails

import (
	"fmt"
	"os/exec"
	"strings"
)

func run(name string, args ...string) error {

	fmt.Printf("%s %v\n", name, args)
	cmd := exec.Command(name, args...)
	var o, e strings.Builder
	cmd.Stdout = &o
	cmd.Stderr = &e
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Command failed: %s\n%s\n", o.String(), e.String())
		return err
	}
	fmt.Printf("Command: %s\n%s\n", o.String(), e.String())

	return nil
}
