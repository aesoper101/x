package thread

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"
)

func TestParallelizeWithImmediateCancellation(t *testing.T) {
	t.Parallel()
	// The bulk of the code relies on subtle timing that's difficult to
	// reproduce, but we can test the most basic use case.
	t.Run("RegularRun", func(t *testing.T) {
		t.Parallel()
		const jobsToExecute = 10
		var (
			executed atomic.Int64
			jobs     = make([]func(context.Context) error, 0, jobsToExecute)
		)
		for i := 0; i < jobsToExecute; i++ {
			jobs = append(jobs, func(_ context.Context) error {
				executed.Inc()
				return nil
			})
		}
		err := Parallelize(context.Background(), jobs)
		assert.NoError(t, err)
		assert.Equal(t, int64(jobsToExecute), executed.Load(), "jobs executed")
	})
	t.Run("WithCtxCancellation", func(t *testing.T) {
		t.Parallel()
		var executed atomic.Int64
		var jobs []func(context.Context) error
		for i := 0; i < 10; i++ {
			jobs = append(jobs, func(_ context.Context) error {
				executed.Inc()
				return nil
			})
		}
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err := Parallelize(ctx, jobs)
		assert.Error(t, err)
		assert.Equal(t, int64(0), executed.Load(), "jobs executed")
	})
}
