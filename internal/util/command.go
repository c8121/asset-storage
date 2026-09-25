package util

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RunSilent is a convenience method to execute a command and capture its output, which will only be printed on errors.
func RunSilent(name string, args ...string) error {
	return RunSilentWithEnv(name, nil, args...)
}

// RunSilentWithEnv is a convenience method to execute a command and capture its output, which will only be printed on errors.
func RunSilentWithEnv(name string, env []string, args ...string) error {

	//fmt.Printf("%s %v\n", name, args)
	cmd := exec.Command(name, args...)

	if len(env) > 0 {
		cmd.Env = os.Environ()
		for _, e := range env {
			cmd.Env = append(cmd.Env, e)
		}
	}

	var o, e strings.Builder
	cmd.Stdout = &o
	cmd.Stderr = &e
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Command failed: %s\n%s\n", o.String(), e.String())
		return err
	}
	//fmt.Printf("Command: %s\n%s\n", o.String(), e.String())

	return nil
}
