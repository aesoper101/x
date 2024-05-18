package tablewriterx

import (
	"io"
	"os"

	"github.com/olekukonko/tablewriter"
)

type TableWriter struct {
	*tablewriter.Table
}

type TableWriterOptions struct {
	writer io.Writer
}

type TableWriterOption func(*TableWriterOptions)

func WithWriter(writer io.Writer) TableWriterOption {
	return func(options *TableWriterOptions) {
		options.writer = writer
	}
}

func NewWriter(options ...TableWriterOption) *TableWriter {
	opts := &TableWriterOptions{
		writer: os.Stdout,
	}
	for _, o := range options {
		o(opts)
	}
	return &TableWriter{tablewriter.NewWriter(opts.writer)}
}
