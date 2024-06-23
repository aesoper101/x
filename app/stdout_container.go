package app

import (
	"github.com/aesoper101/x/internal/ioext"
	"io"
)

type stdoutContainer struct {
	writer io.Writer
}

func newStdoutContainer(writer io.Writer) *stdoutContainer {
	if writer == nil {
		writer = io.Discard
	}
	return &stdoutContainer{
		writer: ioext.LockedWriter(writer),
	}
}

func (s *stdoutContainer) Stdout() io.Writer {
	return s.writer
}
