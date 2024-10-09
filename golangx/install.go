package golangx

import (
	"context"
	"fmt"
	"github.com/aesoper101/x/execext/command"
	"os"
	"strings"

	"github.com/aesoper101/x/execext"
)

// GoInstall is a functionx that installs a package from a remote repository.
// returns an error if something goes wrong.
func GoInstall(repository string) error {
	if err := execext.LookPath("go"); err != nil {
		return err
	}

	env := os.Environ()

	if !IsGO111ModuleOn() {
		env = append(env, "GO111MODULE=on")
	}

	env = append(env, fmt.Sprintf("GOPROXY=%s", GoProxy()))

	if !strings.Contains(repository, "@") {
		repository += "@latest"
	}

	if strings.Contains(repository, "://") {
		repository = strings.TrimPrefix(repository, "https://")
		repository = strings.TrimPrefix(repository, "http://")
	}

	runner := command.NewRunner()

	return runner.Run(
		context.Background(),
		"go",
		command.RunWithArgs("install", repository),
		command.RunWithEnviron(env),
		command.RunWithStdout(os.Stdout),
		command.RunWithStderr(os.Stderr),
	)
}
