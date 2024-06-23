package command

import (
	"context"
	"errors"
	"os/exec"

	"go.uber.org/multierr"
)

var errWaitAlreadyCalled = errors.New("wait already called on process")

type process struct {
	cmd   *exec.Cmd
	done  func()
	waitC chan error
}

// newProcess wraps an *exec.Cmd and monitors it for exiting.
// When the process exits, done will be called.
//
// This implements the Process interface.
//
// The process is expected to have been started by the caller.
func newProcess(cmd *exec.Cmd, done func()) *process {
	return &process{
		cmd:   cmd,
		done:  done,
		waitC: make(chan error, 1),
	}
}

// Monitor starts monitoring of the *exec.Cmd.
func (p *process) Monitor() {
	go func() {
		p.waitC <- p.cmd.Wait()
		close(p.waitC)
		p.done()
	}()
}

// Wait waits for the process to exit.
func (p *process) Wait(ctx context.Context) error {
	select {
	case err, ok := <-p.waitC:
		// Process exited
		if ok {
			return err
		}
		return errWaitAlreadyCalled
	case <-ctx.Done():
		// Timed out. Send a kill signal and release our handle to it.
		return multierr.Combine(ctx.Err(), p.cmd.Process.Kill())
	}
}
