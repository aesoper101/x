package golangx

import (
	"go/build"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/aesoper101/x/filex"
	"github.com/aesoper101/x/str"
)

var defaultGoProxy = "https://goproxy.cn,direct"

// GoBin returns the path to the go binary.
func GoBin() string {
	bin := os.Getenv("GOBIN")
	if bin != "" {
		if filex.IsExists(bin) {
			return bin
		}
	}

	gopath := GoPath()
	bin = filepath.Join(gopath, "bin")
	if filex.IsExists(bin) {
		return bin
	}

	goroot := GoRoot()
	bin = filepath.Join(goroot, "bin")
	if filex.IsExists(bin) {
		return bin
	}

	ctx := build.Default

	return filepath.Join(ctx.GOPATH, "bin")
}

// GoProxy returns the proxy setting for go. if not set,
// return a default proxy setting: "https://goproxy.cn,direct".
func GoProxy() string {
	proxy := os.Getenv("GOPROXY")
	if proxy != "" {
		return proxy
	}

	if output, _ := exec.Command("go", "env", "GOPROXY").Output(); str.IsNotEmpty(string(output)) {
		return string(output)
	}

	return defaultGoProxy
}

// IsGO111ModuleOn returns true if go111module is set to on.
func IsGO111ModuleOn() bool {
	v := os.Getenv("GO111MODULE")
	if str.IsNotEmpty(v) {
		return strings.ToLower(v) == "on"
	}

	if output, _ := exec.Command("go", "env", "GO111MODULE").Output(); str.IsNotEmpty(string(output)) {
		return strings.ToLower(string(output)) == "on"
	}

	return true
}

// GoPath returns the GOPATH.
func GoPath() string {
	path := os.Getenv("GOPATH")
	if str.IsNotEmpty(path) {
		return splitGoPath(path)[0]
	}

	if output, _ := exec.Command("go", "env", "GOPATH").Output(); str.IsNotEmpty(string(output)) {
		return string(output)
	}

	return build.Default.GOPATH
}

// GoRoot returns the GOROOT.
func GoRoot() string {
	path := os.Getenv("GOROOT")
	if str.IsNotEmpty(path) {
		return path
	}

	if output, _ := exec.Command("go", "env", "GOROOT").Output(); str.IsNotEmpty(string(output)) {
		return string(output)
	}

	return build.Default.GOROOT
}

// GoVersion returns the go version.
func GoVersion() string {
	if output, _ := exec.Command("go", "env", "GOVERSION").Output(); str.IsNotEmpty(string(output)) {
		return strings.ReplaceAll(strings.Replace(string(output), "\n", "", -1), "go", "")
	}

	return ""
}

// GoModCache returns the path of Go module cache.
func GoModCache() string {
	cacheOut, _ := exec.Command("go", "env", "GOMODCACHE").Output()
	cachePath := strings.Trim(string(cacheOut), "\n")
	if str.IsEmpty(cachePath) {
		return filepath.Join(GoPath(), "pkg", "mod")
	}
	return cachePath
}

func splitGoPath(goPath string) []string {
	var sep string
	goos := strings.ToLower(runtime.GOOS)
	switch goos {
	case "windows":
		sep = ";"
	case "linux", "darwin":
		sep = ":"
	}

	return strings.Split(goPath, sep)
}
