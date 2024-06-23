package filelock

import (
	"context"
	"time"
)

const (
	// DefaultLockTimeout is the default lock timeout.
	DefaultLockTimeout = 3 * time.Second
	// DefaultLockRetryDelay is the default lock retry delay.
	DefaultLockRetryDelay = 200 * time.Millisecond
)

// Unlocker unlocks a file lock.
type Unlocker interface {
	Unlock() error
}

// Lock locks a file lock.
//
// Use in cases where you need a lock for a specific system file, such as in testing,
// otherwise use a Locker to manage your file locks.
func Lock(ctx context.Context, filePath string, options ...LockOption) (Unlocker, error) {
	return lock(ctx, filePath, options...)
}

// RLock read-locks a file lock.
//
// Use in cases where you need a lock for a specific system file, such as in testing,
// otherwise use a Locker to manage your file locks.
func RLock(ctx context.Context, filePath string, options ...LockOption) (Unlocker, error) {
	return rlock(ctx, filePath, options...)
}

// Locker provides file locks.
type Locker interface {
	// Lock locks a file lock within the root directory of the Locker.
	//
	// The given path must be normalized and relative.
	Lock(ctx context.Context, path string, options ...LockOption) (Unlocker, error)
	// RLock read-locks a file lock within the root directory of the Locker.
	//
	// The given path must be normalized and relative.
	RLock(ctx context.Context, path string, options ...LockOption) (Unlocker, error)
}

// NewLocker returns a new Locker for the given root directory path.
//
// The root directory path should generally be a data directory path.
// The root directory must exist.
func NewLocker(rootDirPath string) (Locker, error) {
	return newLocker(rootDirPath)
}

// LockOption is an option for lock.
type LockOption func(*lockOptions)

// LockWithTimeout returns a new LockOption that sets the lock timeout.
//
// Lock returns error if the lock cannot be acquired after this amount of time.
// If this is set to 0, the lock will never timeout.
func LockWithTimeout(timeout time.Duration) LockOption {
	return func(lockOptions *lockOptions) {
		lockOptions.timeout = timeout
	}
}

// LockWithRetryDelay returns a new LockOption that sets the lock retry delay.
//
// Lock will try to lock on this delay up until the lock timeout.
func LockWithRetryDelay(retryDelay time.Duration) LockOption {
	return func(lockOptions *lockOptions) {
		lockOptions.retryDelay = retryDelay
	}
}

// NewNopLocker returns a new no-op Locker.
func NewNopLocker() Locker {
	return newNopLocker()
}
