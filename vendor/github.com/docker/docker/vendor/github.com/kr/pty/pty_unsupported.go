// +build !linux,!darwin,!freebsd

package pty

import (
	"os"
)

func open() (pty, tty *os.File, err error) ***REMOVED***
	return nil, nil, ErrUnsupported
***REMOVED***
