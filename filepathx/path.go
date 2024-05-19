package filepathx

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func JoinPath(elem ...string) string {
	if runtime.GOOS == "windows" {
		return strings.ReplaceAll(filepath.Join(elem...), "\\", "/")
	}
	return filepath.Join(elem...)
}

func RelativePath(path string) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	ret, _ := filepath.Rel(cwd, path)
	return ret, nil
}
