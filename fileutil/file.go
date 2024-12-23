package fileutil

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// IsFile checks if the given path is a file.
func IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// IfIsFile executes the given functionx if the given path is a file.
// an error is returned if execution fails.
func IfIsFile(path string, f func(p string) error) error {
	if IsFile(path) {
		return f(path)
	}
	return nil
}

// Extension returns the extension of the file name
func Extension(filePath string) string {
	return strings.TrimLeft(filepath.Ext(filePath), "")
}

// IsExists returns true if the file exists.
func IsExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// IfExists checks if the file exists and executes the functionx.
func IfExists(path string, fn func(path string) error) error {
	if IsExists(path) && fn != nil {
		return fn(path)
	}
	return nil
}

// IsNotExists returns true if the file does not exist.
func IsNotExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return false
	}

	if os.IsNotExist(err) {
		return true
	}

	return false
}

// IfNotExists executes the functionx if the file does not exist.
// returns error if execution fails.
func IfNotExists(path string, fn func(path string) error) error {
	if IsNotExists(path) && fn != nil {
		return fn(path)
	}
	return nil
}

// IsRegular returns true if the file is a regular file.
func IsRegular(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsRegular()
}

// IsSymlink returns true if the file is a symlink.
func IsSymlink(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.Mode() == os.ModeSymlink
}

// IfIsRegular executes the functionx if the file is a regular file.
// returns error if execution fails.
func IfIsRegular(path string, fn func() error) error {
	if IsRegular(path) && fn != nil {
		return fn()
	}
	return nil
}

// MkdirIfNotExist creates a directory if it does not exist.
func MkdirIfNotExist(dir string) error {
	return IfNotExists(
		dir, func(dir string) error {
			return os.MkdirAll(dir, os.ModePerm)
		},
	)
}

// CreateFileFromByteFn creates a file from the functionx that returns a byte slice.
func CreateFileFromByteFn(filename string, overwrite bool, f func() []byte) error {
	if IsExists(filename) && !overwrite {
		return nil
	}

	dir := filepath.Dir(filename)
	if err := MkdirIfNotExist(dir); err != nil {
		return err
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer func(w *os.File) {
		_ = w.Close()
	}(file)

	_, err = file.Write(f())
	return err
}

// CreateFileFromReader creates a file from the reader.
func CreateFileFromReader(file string, overwrite bool, reader io.Reader) error {
	b := bytes.NewBuffer([]byte{})
	if _, err := b.ReadFrom(reader); err != nil {
		return err
	}
	return CreateFileFromByteFn(
		file, overwrite, func() []byte {
			return b.Bytes()
		},
	)
}

// CreateFileFromWriterFunc creates a file from the functionx that returns a writer.
func CreateFileFromWriterFunc(file string, overwrite bool, fn func(w io.Writer) error) error {
	b := bytes.NewBuffer([]byte{})
	if err := fn(b); err != nil {
		return err
	}
	return CreateFileFromByteFn(
		file, overwrite, func() []byte {
			return b.Bytes()
		},
	)
}

// CreateFileFromString creates a file from the string.
func CreateFileFromString(file string, overwrite bool, content string) error {
	return CreateFileFromByteFn(
		file, overwrite, func() []byte {
			return []byte(content)
		},
	)
}

// CreateFileFromBytes creates a file from the byte slice.
func CreateFileFromBytes(file string, overwrite bool, content []byte) error {
	return CreateFileFromByteFn(
		file, overwrite, func() []byte {
			return content
		},
	)
}

// DeleteFile deletes a file.
func DeleteFile(file string) error {
	return os.Remove(file)
}

// DeleteFileFn executes the functionx if the file is deleted.
func DeleteFileFn(file string, fn func() error) error {
	if err := DeleteFile(file); err != nil {
		return err
	}
	return fn()
}

func Getwd() string {
	wd, _ := os.Getwd()
	return wd
}

// AfterFileRead is a function that executes a function when a file is read.
func AfterFileRead(filename string, callback func([]byte) error) error {
	bs, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	return callback(bs)
}

func IsReadableFile(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	if info.IsDir() {
		return false, errors.New("is directory")
	}
	// Is it readable?
	fd, err := os.Open(path)
	if err != nil {
		return false, errors.New("permission denied")
	}
	_ = fd.Close()
	return true, nil // Item exists and is readable.
}
