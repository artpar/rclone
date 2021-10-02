// This makes the open test suite. It tries to open a file (existing
// or not existing) with all possible file modes and writes a test
// matrix.
//
// The behaviour is as run on Linux, with the small modification that
// O_TRUNC with O_RDONLY does **not** truncate the file.
//
// Run with go generate (defined in vfs.go)
//
//go:build none
// +build none

// FIXME include read too?

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/artpar/rclone/lib/file"
)

// Interprets err into a vfs error
func whichError(err error) string {
	switch err {
	case nil:
		return "nil"
	case io.EOF:
		return "io.EOF"
	case os.ErrInvalid:
		return "EINVAL"
	}
	s := err.Error()
	switch {
	case strings.Contains(s, "no such file or directory"):
		return "ENOENT"
	case strings.Contains(s, "bad file descriptor"):
		return "EBADF"
	case strings.Contains(s, "file exists"):
		return "EEXIST"
	}
	log.Printf("Unknown error: %v", err)
	return ""
}

const accessModeMask = (os.O_RDONLY | os.O_WRONLY | os.O_RDWR)

// test Opening, reading and writing the file handle with the flags given
func test(fileName string, flags int, mode string) {
	// first try with file not existing
	_, err := os.Stat(fileName)
	if !os.IsNotExist(err) {
		log.Printf("File must not exist")
	}
	f, openNonExistentErr := file.OpenFile(fileName, flags, 0666)

	var readNonExistentErr error
	var writeNonExistentErr error
	if openNonExistentErr == nil {
		// read some bytes
		buf := []byte{0, 0}
		_, readNonExistentErr = f.Read(buf)

		// write some bytes
		_, writeNonExistentErr = f.Write([]byte("hello"))

		// close
		err = f.Close()
		if err != nil {
			log.Printf("failed to close: %v", err)
		}
	}

	// write the file
	f, err = file.Create(fileName)
	if err != nil {
		log.Printf("failed to create: %v", err)
	}
	n, err := f.Write([]byte("hello"))
	if n != 5 || err != nil {
		log.Printf("failed to write n=%d: %v", n, err)
	}
	// close
	err = f.Close()
	if err != nil {
		log.Printf("failed to close: %v", err)
	}

	// then open file and try with file existing

	f, openExistingErr := file.OpenFile(fileName, flags, 0666)
	var readExistingErr error
	var writeExistingErr error
	if openExistingErr == nil {
		// read some bytes
		buf := []byte{0, 0}
		_, readExistingErr = f.Read(buf)

		// write some bytes
		_, writeExistingErr = f.Write([]byte("HEL"))

		// close
		err = f.Close()
		if err != nil {
			log.Printf("failed to close: %v", err)
		}
	}

	// read the file
	f, err = file.Open(fileName)
	if err != nil {
		log.Printf("failed to open: %v", err)
	}
	var buf = make([]byte, 64)
	n, err = f.Read(buf)
	if err != nil && err != io.EOF {
		log.Printf("failed to read n=%d: %v", n, err)
	}
	err = f.Close()
	if err != nil {
		log.Printf("failed to close: %v", err)
	}
	contents := string(buf[:n])

	// remove file
	err = os.Remove(fileName)
	if err != nil {
		log.Printf("failed to remove: %v", err)
	}

	// http://pubs.opengroup.org/onlinepubs/7908799/xsh/open.html
	// The result of using O_TRUNC with O_RDONLY is undefined.
	// Linux seems to truncate the file, but we prefer to return EINVAL
	if (flags&accessModeMask) == os.O_RDONLY && flags&os.O_TRUNC != 0 {
		openNonExistentErr = os.ErrInvalid // EINVAL
		readNonExistentErr = nil
		writeNonExistentErr = nil
		openExistingErr = os.ErrInvalid // EINVAL
		readExistingErr = nil
		writeExistingErr = nil
		contents = "hello"
	}

	// output the struct
	fmt.Printf(`{
	flags: %s,
	what: %q,
	openNonExistentErr: %s,
	readNonExistentErr: %s,
	writeNonExistentErr: %s,
	openExistingErr: %s,
	readExistingErr: %s,
	writeExistingErr: %s,
	contents: %q,
},`,
		mode,
		mode,
		whichError(openNonExistentErr),
		whichError(readNonExistentErr),
		whichError(writeNonExistentErr),
		whichError(openExistingErr),
		whichError(readExistingErr),
		whichError(writeExistingErr),
		contents)
}

func main() {
	fmt.Printf(`// Code generated by make_open_tests.go - use go generate to rebuild - DO NOT EDIT

package vfs

import (
	"os"
	"io"
)

// openTest describes a test of OpenFile
type openTest struct{
	flags int
	what string
	openNonExistentErr error
	readNonExistentErr error
	writeNonExistentErr error
	openExistingErr error
	readExistingErr error
	writeExistingErr error
	contents string
}

// openTests is a suite of tests for OpenFile with all possible
// combination of flags.  This obeys Unix semantics even on Windows.
var openTests = []openTest{
`)
	f, err := ioutil.TempFile("", "open-test")
	if err != nil {
		log.Print(err)
	}
	fileName := f.Name()
	_ = f.Close()
	err = os.Remove(fileName)
	if err != nil {
		log.Printf("failed to remove: %v", err)
	}
	for _, rwMode := range []int{os.O_RDONLY, os.O_WRONLY, os.O_RDWR} {
		flags0 := rwMode
		parts0 := []string{"os.O_RDONLY", "os.O_WRONLY", "os.O_RDWR"}[rwMode : rwMode+1]
		for _, appendMode := range []int{0, os.O_APPEND} {
			flags1 := flags0 | appendMode
			parts1 := parts0
			if appendMode != 0 {
				parts1 = append(parts1, "os.O_APPEND")
			}
			for _, createMode := range []int{0, os.O_CREATE} {
				flags2 := flags1 | createMode
				parts2 := parts1
				if createMode != 0 {
					parts2 = append(parts2, "os.O_CREATE")
				}
				for _, exclMode := range []int{0, os.O_EXCL} {
					flags3 := flags2 | exclMode
					parts3 := parts2
					if exclMode != 0 {
						parts3 = append(parts2, "os.O_EXCL")
					}
					for _, syncMode := range []int{0, os.O_SYNC} {
						flags4 := flags3 | syncMode
						parts4 := parts3
						if syncMode != 0 {
							parts4 = append(parts4, "os.O_SYNC")
						}
						for _, truncMode := range []int{0, os.O_TRUNC} {
							flags5 := flags4 | truncMode
							parts5 := parts4
							if truncMode != 0 {
								parts5 = append(parts5, "os.O_TRUNC")
							}
							textMode := strings.Join(parts5, "|")
							flags := flags5

							test(fileName, flags, textMode)
						}
					}
				}
			}
		}
	}
	fmt.Printf("\n}\n")
}
