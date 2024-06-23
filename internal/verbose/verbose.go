package verbose

import (
	"fmt"
	"io"
	"strings"
)

var (
	// NopPrinter is a no-op printer.
	//
	// This generally aligns with the --verbose flag not being set.
	NopPrinter = nopPrinter{}
)

// Printer prints verbose messages.
type Printer interface {
	// Enabled returns true if verbose mode is enabled.
	//
	// This is false if the Printer is a no-op printer.
	Enabled() bool
	// Printf prints a new verbose message.
	//
	// Leading and trailing newlines are not respected.
	//
	// Callers should not rely on the print calls being reliable, i.e. errors to
	// a backing Writer will be ignored.
	Printf(format string, args ...interface{})

	isPrinter()
}

// NewPrinter returns a new Printer using the given Writer.
//
// The trimmed prefix is printed with a : before each line.
//
// This generally aligns with the --verbose flag being set and writer being stderr.
func NewPrinter(writer io.Writer, prefix string) Printer {
	return newWritePrinter(writer, prefix)
}

// NewPrinterForFlagValue returns a new Printer for the given verboseValue flag value.
func NewPrinterForFlagValue(writer io.Writer, prefix string, verboseValue bool) Printer {
	if verboseValue {
		return NewPrinter(writer, prefix)
	}
	return NopPrinter
}

type nopPrinter struct{}

func (nopPrinter) Printf(string, ...interface{}) {}

func (nopPrinter) Enabled() bool {
	return false
}

func (nopPrinter) isPrinter() {}

type writePrinter struct {
	writer io.Writer
	prefix string
}

func newWritePrinter(writer io.Writer, prefix string) *writePrinter {
	prefix = strings.TrimSpace(prefix)
	if prefix != "" {
		prefix = prefix + ": "
	}
	return &writePrinter{
		writer: writer,
		prefix: prefix,
	}
}

func (w *writePrinter) Printf(format string, args ...interface{}) {
	if value := strings.TrimSpace(fmt.Sprintf(format, args...)); value != "" {
		// Errors are ignored per the interface spec.
		_, _ = w.writer.Write([]byte(w.prefix + value + "\n"))
	}
}

func (*writePrinter) Enabled() bool {
	return true
}

func (*writePrinter) isPrinter() {}
