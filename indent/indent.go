// Package indent handles printing with indentation, mostly for debug purposes.
package indent

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"unicode"
)

// Printer prints with potential indents.
//
// Not thread-safe.
type Printer interface {
	// P prints the args with fmt.Sprint on a line after applying the current indent.
	P(args ...any)
	// Pf prints the format and args with fmt.Sprintf on a line after applying the current indent.
	Pf(format string, args ...any)
	// In indents by one.
	In()
	// Out unindents by one.
	Out()
	// String gets the resulting string represntation.
	//
	// Returns error if there was an error during printing.
	String() (string, error)
	// Bytes gets the resulting bytes representation.
	//
	// Returns error if there was an error during printing.
	Bytes() ([]byte, error)

	isPrinter()
}

// NewPrinter returns a new Printer.
func NewPrinter(indent string) Printer {
	return newPrinter(indent)
}

// *** PRIVATE ***

type printer struct {
	indent         string
	buffer         *bytes.Buffer
	curIndentCount int
	err            error
}

func newPrinter(indent string) *printer {
	return &printer{
		indent:         indent,
		buffer:         bytes.NewBuffer(nil),
		curIndentCount: 0,
		err:            nil,
	}
}

func (p *printer) P(args ...any) {
	if p.err != nil {
		return
	}
	p.pString(fmt.Sprint(args...))
}

func (p *printer) Pf(format string, args ...any) {
	if p.err != nil {
		return
	}
	p.pString(fmt.Sprintf(format, args...))
}

func (p *printer) In() {
	if p.err != nil {
		return
	}
	p.curIndentCount++
}

func (p *printer) Out() {
	if p.err != nil {
		return
	}
	if p.curIndentCount <= 0 {
		p.err = errors.New("printer indent count is 0 and Out called")
		return
	}
	p.curIndentCount--
}

func (p *printer) String() (string, error) {
	if p.err != nil {
		return "", p.err
	}
	return p.buffer.String(), nil
}

func (p *printer) Bytes() ([]byte, error) {
	if p.err != nil {
		return nil, p.err
	}
	return p.buffer.Bytes(), nil
}

func (p *printer) pString(s string) {
	s = strings.TrimSpace(s)
	if s == "" {
		_, _ = p.buffer.WriteRune('\n')
		return
	}
	if p.curIndentCount > 0 {
		s = strings.Repeat(p.indent, p.curIndentCount) + s
	}
	s = strings.TrimRightFunc(s, unicode.IsSpace)
	if strings.TrimSpace(s) != "" {
		_, _ = p.buffer.WriteString(s)
	}
	_, _ = p.buffer.WriteRune('\n')
}

func (*printer) isPrinter() {}
