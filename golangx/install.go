package golangx

import (
	"fmt"
	"os"
	"os/exec"
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

	cmd := exec.Command("go", "install", repository)
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
