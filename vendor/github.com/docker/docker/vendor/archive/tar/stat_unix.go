// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux darwin dragonfly freebsd openbsd netbsd solaris

package tar

import (
	"os"
	"syscall"
)

func init() ***REMOVED***
	sysStat = statUnix
***REMOVED***

func statUnix(fi os.FileInfo, h *Header) error ***REMOVED***
	sys, ok := fi.Sys().(*syscall.Stat_t)
	if !ok ***REMOVED***
		return nil
	***REMOVED***
	h.Uid = int(sys.Uid)
	h.Gid = int(sys.Gid)
	// TODO(bradfitz): populate username & group.  os/user
	// doesn't cache LookupId lookups, and lacks group
	// lookup functions.
	h.AccessTime = statAtime(sys)
	h.ChangeTime = statCtime(sys)
	// TODO(bradfitz): major/minor device numbers?
	return nil
***REMOVED***
