package tmp

import (
	"github.com/aesoper101/x/interrupt"
	"github.com/aesoper101/x/uuidutil"
	"io"
	"os"
	"path/filepath"

	"go.uber.org/multierr"
)

// File is a temporary file
//
// It must be closed when done.
type File interface {
	io.Closer

	AbsPath() string
}

// NewFileWithData returns a new temporary file with the given data.
//
// It must be closed when done.
// This file will be deleted on interrupt signals.
//
// Usage of this function requires eng approval - ask before using.
func NewFileWithData(data []byte) (File, error) {
	id, err := uuidutil.New()
	if err != nil {
		return nil, err
	}
	file, err := os.CreateTemp("", id.String())
	if err != nil {
		return nil, err
	}
	path := file.Name()
	// just in case
	absPath, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	signalC, closer := interrupt.NewSignalChannel()
	go func() {
		<-signalC
		_ = os.Remove(absPath)
	}()
	_, err = file.Write(data)
	err = multierr.Append(err, file.Close())
	if err != nil {
		err = multierr.Append(err, os.Remove(absPath))
		closer()
		return nil, err
	}
	return newFile(absPath, closer), nil
}

// Dir is a temporary directory.
//
// It must be closed when done.
type Dir interface {
	io.Closer

	AbsPath() string
}

// NewDir returns a new temporary directory.
//
// It must be closed when done.
// This file will be deleted on interrupt signals.
//
// Usage of this function requires eng approval - ask before using.
func NewDir(options ...DirOption) (Dir, error) {
	dirOptions := newDirOptions()
	for _, option := range options {
		option(dirOptions)
	}
	id, err := uuidutil.New()
	if err != nil {
		return nil, err
	}
	path, err := os.MkdirTemp(dirOptions.basePath, id.String())
	if err != nil {
		return nil, err
	}
	// just in case
	absPath, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	signalC, closer := interrupt.NewSignalChannel()
	go func() {
		<-signalC
		_ = os.RemoveAll(absPath)
	}()
	return newDir(absPath, closer), nil
}

// DirOption is an option for NewDir.
type DirOption func(*dirOptions)

// DirWithBasePath returns a new DirOption that sets the base path to create
// the temporary directory in.
//
// The default is to use os.TempDir().
func DirWithBasePath(basePath string) DirOption {
	return func(dirOptions *dirOptions) {
		dirOptions.basePath = basePath
	}
}

type file struct {
	absPath string
	closer  func()
}

func newFile(absPath string, closer func()) *file {
	return &file{
		absPath: absPath,
		closer:  closer,
	}
}

func (f *file) AbsPath() string {
	return f.absPath
}

func (f *file) Close() error {
	err := os.Remove(f.absPath)
	f.closer()
	return err
}

type dir struct {
	absPath string
	closer  func()
}

func newDir(absPath string, closer func()) *dir {
	return &dir{
		absPath: absPath,
		closer:  closer,
	}
}

func (d *dir) AbsPath() string {
	return d.absPath
}

func (d *dir) Close() error {
	err := os.RemoveAll(d.absPath)
	d.closer()
	return err
}

type dirOptions struct {
	basePath string
}

func newDirOptions() *dirOptions {
	return &dirOptions{}
}
