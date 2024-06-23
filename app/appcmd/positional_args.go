package appcmd

import "github.com/spf13/cobra"

var (
	// NoArgs matches cobra.NoArgs.
	NoArgs = newPositionalArgs(cobra.NoArgs)
	// OnlyValidArgs matches cobra.OnlyValidArgs.
	OnlyValidArgs = newPositionalArgs(cobra.OnlyValidArgs)
	// ArbitraryArgs matches cobra.ArbitraryArgs.
	ArbitraryArgs = newPositionalArgs(cobra.ArbitraryArgs)
)

// MinimumNArgs matches cobra.MinimumNArgs.
func MinimumNArgs(n int) PositionalArgs {
	return newPositionalArgs(cobra.MinimumNArgs(n))
}

// MaximumNArgs matches cobra.MaximumNArgs.
func MaximumNArgs(n int) PositionalArgs {
	return newPositionalArgs(cobra.MaximumNArgs(n))
}

// ExactArgs matches cobra.ExactArgs.
func ExactArgs(n int) PositionalArgs {
	return newPositionalArgs(cobra.ExactArgs(n))
}

// RangeArgs matches cobra.RangeArgs.
func RangeArgs(min int, max int) PositionalArgs {
	return newPositionalArgs(cobra.RangeArgs(min, max))
}

// PostionalArgs matches cobra.PositionalArgs so that importers of appcmd do
// not need to reference cobra (and shouldn't).
type PositionalArgs interface {
	cobra() cobra.PositionalArgs
}

// *** PRIVATE ***

type positionalArgs struct {
	args cobra.PositionalArgs
}

func newPositionalArgs(args cobra.PositionalArgs) *positionalArgs {
	return &positionalArgs{
		args: args,
	}
}

func (p *positionalArgs) cobra() cobra.PositionalArgs {
	return p.args
}
