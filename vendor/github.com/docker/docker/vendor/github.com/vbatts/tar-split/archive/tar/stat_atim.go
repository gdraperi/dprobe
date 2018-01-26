// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux dragonfly openbsd solaris

package tar

import (
	"syscall"
	"time"
)

func statAtime(st *syscall.Stat_t) time.Time ***REMOVED***
	return time.Unix(st.Atim.Unix())
***REMOVED***

func statCtime(st *syscall.Stat_t) time.Time ***REMOVED***
	return time.Unix(st.Ctim.Unix())
***REMOVED***
