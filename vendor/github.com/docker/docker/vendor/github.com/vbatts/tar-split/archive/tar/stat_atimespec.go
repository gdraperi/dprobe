// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin freebsd netbsd

package tar

import (
	"syscall"
	"time"
)

func statAtime(st *syscall.Stat_t) time.Time ***REMOVED***
	return time.Unix(st.Atimespec.Unix())
***REMOVED***

func statCtime(st *syscall.Stat_t) time.Time ***REMOVED***
	return time.Unix(st.Ctimespec.Unix())
***REMOVED***
