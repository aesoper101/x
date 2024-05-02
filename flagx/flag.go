package flagx

import "github.com/spf13/pflag"

func NewFlagSet(name string) *pflag.FlagSet {
	return pflag.NewFlagSet(name, pflag.ContinueOnError)
}
