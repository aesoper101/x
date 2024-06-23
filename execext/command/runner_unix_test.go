//go:build aix || darwin || dragonfly || freebsd || (js && wasm) || linux || netbsd || openbsd || solaris
// +build aix darwin dragonfly freebsd js,wasm linux netbsd openbsd solaris

package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDoubleWait(t *testing.T) {
	t.Parallel()

	runner := NewRunner()
	process, err := runner.Start("echo")
	require.NoError(t, err)
	ctx := context.Background()
	_ = process.Wait(ctx)
	require.Equal(t, process.Wait(ctx), errWaitAlreadyCalled)
}

func TestNoDeadlock(t *testing.T) {
	t.Parallel()

	runner := NewRunner(RunnerWithParallelism(2))
	processes := make([]Process, 4)
	for i := 0; i < 4; i++ {
		process, err := runner.Start("sleep", StartWithArgs("1"))
		require.NoError(t, err)
		processes[i] = process
	}
	ctx := context.Background()
	for _, process := range processes {
		require.NoError(t, process.Wait(ctx))
	}
}
