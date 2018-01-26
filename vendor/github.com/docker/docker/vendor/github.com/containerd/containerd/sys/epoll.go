// +build linux

package sys

import "golang.org/x/sys/unix"

// EpollCreate1 directly calls unix.EpollCreate1
func EpollCreate1(flag int) (int, error) ***REMOVED***
	return unix.EpollCreate1(flag)
***REMOVED***

// EpollCtl directly calls unix.EpollCtl
func EpollCtl(epfd int, op int, fd int, event *unix.EpollEvent) error ***REMOVED***
	return unix.EpollCtl(epfd, op, fd, event)
***REMOVED***

// EpollWait directly calls unix.EpollWait
func EpollWait(epfd int, events []unix.EpollEvent, msec int) (int, error) ***REMOVED***
	return unix.EpollWait(epfd, events, msec)
***REMOVED***
