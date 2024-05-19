package execx

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

func LookPath(commandName ...string) error {
	if len(commandName) == 0 {
		return errors.New("no command specified")
	}
	for _, c := range commandName {
		if _, err := exec.LookPath(c); err != nil {
			return err
		}
	}
	return nil
}

func ExecCmd(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("`%v %v` failed: %v %v (%v)", name, strings.Join(args, " "), stderr.String(), stdout.String(), err)
	}

	return stdout.String(), nil
}
