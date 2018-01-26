// +build darwin freebsd

package sys

import (
	"syscall"
)

// StatAtime returns the access time from a stat struct
func StatAtime(st *syscall.Stat_t) syscall.Timespec ***REMOVED***
	return st.Atimespec
***REMOVED***

// StatCtime returns the created time from a stat struct
func StatCtime(st *syscall.Stat_t) syscall.Timespec ***REMOVED***
	return st.Ctimespec
***REMOVED***

// StatMtime returns the modified time from a stat struct
func StatMtime(st *syscall.Stat_t) syscall.Timespec ***REMOVED***
	return st.Mtimespec
***REMOVED***
