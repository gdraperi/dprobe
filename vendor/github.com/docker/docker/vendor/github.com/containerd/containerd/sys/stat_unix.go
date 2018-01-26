// +build linux solaris

package sys

import (
	"syscall"
)

// StatAtime returns the Atim
func StatAtime(st *syscall.Stat_t) syscall.Timespec ***REMOVED***
	return st.Atim
***REMOVED***

// StatCtime returns the Ctim
func StatCtime(st *syscall.Stat_t) syscall.Timespec ***REMOVED***
	return st.Ctim
***REMOVED***

// StatMtime returns the Mtim
func StatMtime(st *syscall.Stat_t) syscall.Timespec ***REMOVED***
	return st.Mtim
***REMOVED***
