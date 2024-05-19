package filepathx

import (
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
