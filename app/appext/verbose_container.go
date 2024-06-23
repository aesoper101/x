package appext

import "github.com/aesoper101/x/internal/verbose"

type verboseContainer struct {
	verbosePrinter verbose.Printer
}

func newVerboseContainer(verbosePrinter verbose.Printer) *verboseContainer {
	return &verboseContainer{
		verbosePrinter: verbosePrinter,
	}
}

func (c *verboseContainer) VerboseEnabled() bool {
	return c.verbosePrinter.Enabled()
}

func (c *verboseContainer) VerbosePrinter() verbose.Printer {
	return c.verbosePrinter
}
