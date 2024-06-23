package app

import (
	"github.com/aesoper101/x/internal/ioext"
	"io"
)

type stdinContainer struct {
	reader io.Reader
}

func newStdinContainer(reader io.Reader) *stdinContainer {
	if reader == nil {
		reader = ioext.DiscardReader
	}
	return &stdinContainer{
		reader: reader,
	}
}

func (s *stdinContainer) Stdin() io.Reader {
	return s.reader
}
