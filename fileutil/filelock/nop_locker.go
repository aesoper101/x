package filelock

import "context"

type nopLocker struct{}

func newNopLocker() *nopLocker {
	return &nopLocker{}
}

func (l *nopLocker) Lock(ctx context.Context, path string, options ...LockOption) (Unlocker, error) {
	return newNopUnlocker(), nil
}

func (l *nopLocker) RLock(ctx context.Context, path string, options ...LockOption) (Unlocker, error) {
	return newNopUnlocker(), nil
}
