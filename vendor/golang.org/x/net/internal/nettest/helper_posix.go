// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux netbsd openbsd solaris windows

package nettest

import (
	"os"
	"syscall"
)

func protocolNotSupported(err error) bool ***REMOVED***
	switch err := err.(type) ***REMOVED***
	case syscall.Errno:
		switch err ***REMOVED***
		case syscall.EPROTONOSUPPORT, syscall.ENOPROTOOPT:
			return true
		***REMOVED***
	case *os.SyscallError:
		switch err := err.Err.(type) ***REMOVED***
		case syscall.Errno:
			switch err ***REMOVED***
			case syscall.EPROTONOSUPPORT, syscall.ENOPROTOOPT:
				return true
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***
