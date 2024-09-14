package interrupt

import (
	"context"
	"os"
	"os/signal"
)

var signals = append(
	[]os.Signal{
		os.Interrupt,
	},
	extraSignals...,
)

// WithCancel returns a context that is cancelled if interrupt signals are sent.
func WithCancel(ctx context.Context) (context.Context, context.CancelFunc) {
	signalC, closer := NewSignalChannel()
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		select {
		case <-signalC:
			cancel()
		case <-ctx.Done():
			closer()
		}
	}()
	return ctx, cancel
}

// NewSignalChannel returns a new channel for interrupt signals.
//
// Call the returned function to cancel sending to this channel.
func NewSignalChannel() (<-chan os.Signal, func()) {
	signalC := make(chan os.Signal, 1)
	signal.Notify(signalC, signals...)
	return signalC, func() {
		signal.Stop(signalC)
		close(signalC)
	}
}
