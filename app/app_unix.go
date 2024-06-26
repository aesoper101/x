//go:build aix || darwin || dragonfly || freebsd || (js && wasm) || linux || netbsd || openbsd || solaris
// +build aix darwin dragonfly freebsd js,wasm linux netbsd openbsd solaris

package app

import (
	"errors"
	"fmt"
	"path/filepath"
)

const (
	// DevStdinFilePath is the equivalent of /dev/stdin.
	//
	// This will be /dev/stdin for darwin and linux.
	// This does not exist for windows.
	DevStdinFilePath = "/dev/stdin"
	// DevStdoutFilePath is the equivalent of /dev/stdout.
	//
	// This will be /dev/stdout for darwin and linux.
	// This does not exist for windows.
	DevStdoutFilePath = "/dev/stdout"
	// DevStderrFilePath is the equivalent of /dev/stderr.
	//
	// This will be /dev/stderr for darwin and linux.
	// This does not exist for windows.
	DevStderrFilePath = "/dev/stderr"
	// DevNullFilePath is the equivalent of /dev/null.
	//
	// This will be /dev/null for darwin and linux.
	// This will be nul for windows.
	DevNullFilePath = "/dev/null"
)

// HomeDirPath returns the home directory path.
//
// This will be $HOME for darwin and linux.
// This will be %USERPROFILE% for windows.
func HomeDirPath(envContainer EnvContainer) (string, error) {
	if home := envContainer.Env("HOME"); home != "" {
		return home, nil
	}
	return "", errors.New("$HOME is not set")
}

// CacheDirPath returns the cache directory path.
//
// This will be $XDG_CACHE_HOME for darwin and linux, falling back to $HOME/.cache.
// This will be %LocalAppData% for windows.
//
// See https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html
// for darwin and linux. Note that we use the same for darwin and linux as this is
// what developers expect, as opposed to ~/Library/Preferences etc as the stdlib
// does for Go.
//
// Users cannot assume that CacheDirPath, ConfigDirPath, and DataDirPath are unique.
func CacheDirPath(envContainer EnvContainer) (string, error) {
	return xdgDirPath(envContainer, "XDG_CACHE_HOME", ".cache")
}

// ConfigDirPath returns the config directory path.
//
// This will be $XDG_CONFIG_HOME for darwin and linux, falling back to $HOME/.config.
// This will be %AppData% for windows.
//
// See https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html
// for darwin and linux. Note that we use the same for darwin and linux as this is
// what developers expect, as opposed to ~/Library/Preferences etc as the stdlib
// does for Go.
//
// Users cannot assume that CacheDirPath, ConfigDirPath, and DataDirPath are unique.
func ConfigDirPath(envContainer EnvContainer) (string, error) {
	return xdgDirPath(envContainer, "XDG_CONFIG_HOME", ".config")
}

// DataDirPath returns the data directory path.
//
// This will be $XDG_DATA_HOME for darwin and linux, falling back to $HOME/.local/share.
// This will be %LocalAppData% for windows.
//
// See https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html
// for darwin and linux. Note that we use the same for darwin and linux as this is
// what developers expect, as opposed to ~/Library/Preferences etc as the stdlib
// does for Go.
//
// Users cannot assume that CacheDirPath, ConfigDirPath, and DataDirPath are unique.
func DataDirPath(envContainer EnvContainer) (string, error) {
	return xdgDirPath(envContainer, "XDG_DATA_HOME", filepath.Join(".local", "share"))
}

func xdgDirPath(envContainer EnvContainer, key string, fallbackRelHomeDirPath string) (string, error) {
	if value := envContainer.Env(key); value != "" {
		return value, nil
	}
	if home := envContainer.Env("HOME"); home != "" {
		return filepath.Join(home, fallbackRelHomeDirPath), nil
	}
	return "", fmt.Errorf("$%s and $HOME are not set", key)
}
