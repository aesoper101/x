package filelock

import (
	"context"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestGlobalBasic(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	tempDirPath := t.TempDir()
	filePath := filepath.Join(tempDirPath, "path/to/lock")
	unlocker, err := Lock(ctx, filePath)
	require.NoError(t, err)
	_, err = Lock(ctx, filePath, LockWithTimeout(100*time.Millisecond), LockWithRetryDelay(10*time.Millisecond))
	require.Error(t, err)
	require.NoError(t, unlocker.Unlock())
	unlocker, err = Lock(ctx, filePath, LockWithTimeout(100*time.Millisecond), LockWithRetryDelay(10*time.Millisecond))
	require.NoError(t, err)
	require.NoError(t, unlocker.Unlock())
	unlocker, err = RLock(ctx, filePath)
	require.NoError(t, err)
	unlocker2, err := RLock(ctx, filePath)
	require.NoError(t, err)
	_, err = Lock(ctx, filePath, LockWithTimeout(100*time.Millisecond), LockWithRetryDelay(10*time.Millisecond))
	require.Error(t, err)
	require.NoError(t, unlocker.Unlock())
	require.NoError(t, unlocker2.Unlock())
}

func TestLockerBasic(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	tempDirPath := t.TempDir()
	filePath := "path/to/lock"
	locker, err := NewLocker(tempDirPath)
	require.NoError(t, err)
	unlocker, err := locker.Lock(ctx, filePath)
	require.NoError(t, err)
	_, err = locker.Lock(ctx, filePath, LockWithTimeout(100*time.Millisecond), LockWithRetryDelay(10*time.Millisecond))
	require.Error(t, err)
	require.NoError(t, unlocker.Unlock())
	unlocker, err = locker.Lock(ctx, filePath, LockWithTimeout(100*time.Millisecond), LockWithRetryDelay(10*time.Millisecond))
	require.NoError(t, err)
	require.NoError(t, unlocker.Unlock())
	unlocker, err = locker.RLock(ctx, filePath)
	require.NoError(t, err)
	unlocker2, err := locker.RLock(ctx, filePath)
	require.NoError(t, err)
	_, err = locker.Lock(ctx, filePath, LockWithTimeout(100*time.Millisecond), LockWithRetryDelay(10*time.Millisecond))
	require.Error(t, err)
	require.NoError(t, unlocker.Unlock())
	require.NoError(t, unlocker2.Unlock())
	absolutePath := "/not/normalized/and/validated"
	if runtime.GOOS == "windows" {
		absolutePath = "C:\\not\\normalized\\and\\validated"
	}
	_, err = locker.Lock(ctx, absolutePath)
	require.Error(t, err)
}
