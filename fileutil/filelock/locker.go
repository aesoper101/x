package filelock

import (
	"context"
	"fmt"
	"github.com/aesoper101/x/filepathext/normalpath"
	"os"
)

type locker struct {
	rootDirPath string
}

func newLocker(rootDirPath string) (*locker, error) {
	// allow symlinks
	fileInfo, err := os.Stat(normalpath.Unnormalize(rootDirPath))
	if err != nil {
		return nil, err
	}
	if !fileInfo.IsDir() {
		return nil, fmt.Errorf("%q is not a directory", rootDirPath)
	}
	return &locker{
		// do not validate - allow anything including absolute paths and jumping context
		rootDirPath: normalpath.Normalize(rootDirPath),
	}, nil
}

func (l *locker) Lock(ctx context.Context, path string, options ...LockOption) (Unlocker, error) {
	if err := validatePath(path); err != nil {
		return nil, err
	}
	return lock(
		ctx,
		normalpath.Unnormalize(normalpath.Join(l.rootDirPath, path)),
		options...,
	)
}

func (l *locker) RLock(ctx context.Context, path string, options ...LockOption) (Unlocker, error) {
	if err := validatePath(path); err != nil {
		return nil, err
	}
	return rlock(
		ctx,
		normalpath.Unnormalize(normalpath.Join(l.rootDirPath, path)),
		options...,
	)
}

func validatePath(path string) error {
	normalPath, err := normalpath.NormalizeAndValidate(path)
	if err != nil {
		return err
	}
	if path != normalPath {
		// just extra safety
		return fmt.Errorf("expected file lock path %q to be equal to normalized path %q", path, normalPath)
	}
	return nil
}
