package filelock

type nopUnlocker struct{}

func newNopUnlocker() *nopUnlocker {
	return &nopUnlocker{}
}

func (l *nopUnlocker) Unlock() error {
	return nil
}
