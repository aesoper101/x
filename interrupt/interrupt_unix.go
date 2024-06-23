// Matching the unix-like build tags in the Golang source i.e. https://github.com/golang/go/blob/912f0750472dd4f674b69ca1616bfaf377af1805/src/os/file_unix.go#L6

//go:build aix || darwin || dragonfly || freebsd || (js && wasm) || linux || netbsd || openbsd || solaris
// +build aix darwin dragonfly freebsd js,wasm linux netbsd openbsd solaris

package interrupt

import (
	"os"
	"syscall"
)

// extraSignals are signals beyond os.Interrupt that we want to be handled
// as interrupts.
//
// For unix-like platforms, this adds syscall.SIGTERM, although this is only
// tested on darwin and linux, which buf officially supports. Other unix-like
// platforms should have this as well, however.
var extraSignals = []os.Signal{
	syscall.SIGTERM,
}
