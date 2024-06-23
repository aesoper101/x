package filelock

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gofrs/flock"
)

func lock(
	ctx context.Context,
	filePath string,
	options ...LockOption,
) (Unlocker, error) {
	return lockForFunc(
		ctx,
		filePath,
		(*flock.Flock).TryLockContext,
		options...,
	)
}

func rlock(
	ctx context.Context,
	filePath string,
	options ...LockOption,
) (Unlocker, error) {
	return lockForFunc(
		ctx,
		filePath,
		(*flock.Flock).TryRLockContext,
		options...,
	)
}

func lockForFunc(
	ctx context.Context,
	filePath string,
	tryLockContextFunc func(*flock.Flock, context.Context, time.Duration) (bool, error),
	options ...LockOption,
) (Unlocker, error) {
	lockOptions := newLockOptions()
	for _, option := range options {
		option(lockOptions)
	}
	// mkdir is an atomic operation
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return nil, err
	}
	var cancel context.CancelFunc
	if lockOptions.timeout != 0 {
		ctx, cancel = context.WithTimeout(ctx, lockOptions.timeout)
		defer cancel()
	}
	flock := flock.New(filePath)
	locked, err := tryLockContextFunc(flock, ctx, lockOptions.retryDelay)
	if err != nil {
		return nil, fmt.Errorf("could not get file lock %q: %w", filePath, err)
	}
	if !locked {
		return nil, fmt.Errorf("could not lock %q", filePath)
	}
	return flock, nil
}

type lockOptions struct {
	timeout    time.Duration
	retryDelay time.Duration
}

func newLockOptions() *lockOptions {
	return &lockOptions{
		timeout:    DefaultLockTimeout,
		retryDelay: DefaultLockRetryDelay,
	}
}
