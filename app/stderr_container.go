package app

import (
	"github.com/aesoper101/x/internal/ioext"
	"io"
)

type stderrContainer struct {
	writer io.Writer
}

func newStderrContainer(writer io.Writer) *stderrContainer {
	if writer == nil {
		writer = io.Discard
	}
	return &stderrContainer{
		writer: ioext.LockedWriter(writer),
	}
}

func (s *stderrContainer) Stderr() io.Writer {
	return s.writer
}
