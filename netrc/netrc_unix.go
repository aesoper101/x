// Matching the unix-like build tags in the Golang source i.e. https://github.com/golang/go/blob/912f0750472dd4f674b69ca1616bfaf377af1805/src/os/file_unix.go#L6

//go:build aix || darwin || dragonfly || freebsd || (js && wasm) || linux || netbsd || openbsd || solaris
// +build aix darwin dragonfly freebsd js,wasm linux netbsd openbsd solaris

package netrc

// netrcFilename is the netrc filename.
//
// This will be .netrc for unix-like platforms including darwin.
// This will be _netrc for windows.
const netrcFilename = ".netrc"
