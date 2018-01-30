// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package nettest

import (
	"fmt"
	"os"
	"runtime"
	"syscall"
)

func maxOpenFiles() int ***REMOVED***
	var rlim syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlim); err != nil ***REMOVED***
		return defaultMaxOpenFiles
	***REMOVED***
	return int(rlim.Cur)
***REMOVED***

func supportsRawIPSocket() (string, bool) ***REMOVED***
	if os.Getuid() != 0 ***REMOVED***
		return fmt.Sprintf("must be root on %s", runtime.GOOS), false
	***REMOVED***
	return "", true
***REMOVED***
