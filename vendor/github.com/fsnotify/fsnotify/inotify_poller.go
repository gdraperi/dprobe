// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux

package fsnotify

import (
	"errors"

	"golang.org/x/sys/unix"
)

type fdPoller struct ***REMOVED***
	fd   int    // File descriptor (as returned by the inotify_init() syscall)
	epfd int    // Epoll file descriptor
	pipe [2]int // Pipe for waking up
***REMOVED***

func emptyPoller(fd int) *fdPoller ***REMOVED***
	poller := new(fdPoller)
	poller.fd = fd
	poller.epfd = -1
	poller.pipe[0] = -1
	poller.pipe[1] = -1
	return poller
***REMOVED***

// Create a new inotify poller.
// This creates an inotify handler, and an epoll handler.
func newFdPoller(fd int) (*fdPoller, error) ***REMOVED***
	var errno error
	poller := emptyPoller(fd)
	defer func() ***REMOVED***
		if errno != nil ***REMOVED***
			poller.close()
		***REMOVED***
	***REMOVED***()
	poller.fd = fd

	// Create epoll fd
	poller.epfd, errno = unix.EpollCreate1(0)
	if poller.epfd == -1 ***REMOVED***
		return nil, errno
	***REMOVED***
	// Create pipe; pipe[0] is the read end, pipe[1] the write end.
	errno = unix.Pipe2(poller.pipe[:], unix.O_NONBLOCK)
	if errno != nil ***REMOVED***
		return nil, errno
	***REMOVED***

	// Register inotify fd with epoll
	event := unix.EpollEvent***REMOVED***
		Fd:     int32(poller.fd),
		Events: unix.EPOLLIN,
	***REMOVED***
	errno = unix.EpollCtl(poller.epfd, unix.EPOLL_CTL_ADD, poller.fd, &event)
	if errno != nil ***REMOVED***
		return nil, errno
	***REMOVED***

	// Register pipe fd with epoll
	event = unix.EpollEvent***REMOVED***
		Fd:     int32(poller.pipe[0]),
		Events: unix.EPOLLIN,
	***REMOVED***
	errno = unix.EpollCtl(poller.epfd, unix.EPOLL_CTL_ADD, poller.pipe[0], &event)
	if errno != nil ***REMOVED***
		return nil, errno
	***REMOVED***

	return poller, nil
***REMOVED***

// Wait using epoll.
// Returns true if something is ready to be read,
// false if there is not.
func (poller *fdPoller) wait() (bool, error) ***REMOVED***
	// 3 possible events per fd, and 2 fds, makes a maximum of 6 events.
	// I don't know whether epoll_wait returns the number of events returned,
	// or the total number of events ready.
	// I decided to catch both by making the buffer one larger than the maximum.
	events := make([]unix.EpollEvent, 7)
	for ***REMOVED***
		n, errno := unix.EpollWait(poller.epfd, events, -1)
		if n == -1 ***REMOVED***
			if errno == unix.EINTR ***REMOVED***
				continue
			***REMOVED***
			return false, errno
		***REMOVED***
		if n == 0 ***REMOVED***
			// If there are no events, try again.
			continue
		***REMOVED***
		if n > 6 ***REMOVED***
			// This should never happen. More events were returned than should be possible.
			return false, errors.New("epoll_wait returned more events than I know what to do with")
		***REMOVED***
		ready := events[:n]
		epollhup := false
		epollerr := false
		epollin := false
		for _, event := range ready ***REMOVED***
			if event.Fd == int32(poller.fd) ***REMOVED***
				if event.Events&unix.EPOLLHUP != 0 ***REMOVED***
					// This should not happen, but if it does, treat it as a wakeup.
					epollhup = true
				***REMOVED***
				if event.Events&unix.EPOLLERR != 0 ***REMOVED***
					// If an error is waiting on the file descriptor, we should pretend
					// something is ready to read, and let unix.Read pick up the error.
					epollerr = true
				***REMOVED***
				if event.Events&unix.EPOLLIN != 0 ***REMOVED***
					// There is data to read.
					epollin = true
				***REMOVED***
			***REMOVED***
			if event.Fd == int32(poller.pipe[0]) ***REMOVED***
				if event.Events&unix.EPOLLHUP != 0 ***REMOVED***
					// Write pipe descriptor was closed, by us. This means we're closing down the
					// watcher, and we should wake up.
				***REMOVED***
				if event.Events&unix.EPOLLERR != 0 ***REMOVED***
					// If an error is waiting on the pipe file descriptor.
					// This is an absolute mystery, and should never ever happen.
					return false, errors.New("Error on the pipe descriptor.")
				***REMOVED***
				if event.Events&unix.EPOLLIN != 0 ***REMOVED***
					// This is a regular wakeup, so we have to clear the buffer.
					err := poller.clearWake()
					if err != nil ***REMOVED***
						return false, err
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***

		if epollhup || epollerr || epollin ***REMOVED***
			return true, nil
		***REMOVED***
		return false, nil
	***REMOVED***
***REMOVED***

// Close the write end of the poller.
func (poller *fdPoller) wake() error ***REMOVED***
	buf := make([]byte, 1)
	n, errno := unix.Write(poller.pipe[1], buf)
	if n == -1 ***REMOVED***
		if errno == unix.EAGAIN ***REMOVED***
			// Buffer is full, poller will wake.
			return nil
		***REMOVED***
		return errno
	***REMOVED***
	return nil
***REMOVED***

func (poller *fdPoller) clearWake() error ***REMOVED***
	// You have to be woken up a LOT in order to get to 100!
	buf := make([]byte, 100)
	n, errno := unix.Read(poller.pipe[0], buf)
	if n == -1 ***REMOVED***
		if errno == unix.EAGAIN ***REMOVED***
			// Buffer is empty, someone else cleared our wake.
			return nil
		***REMOVED***
		return errno
	***REMOVED***
	return nil
***REMOVED***

// Close all poller file descriptors, but not the one passed to it.
func (poller *fdPoller) close() ***REMOVED***
	if poller.pipe[1] != -1 ***REMOVED***
		unix.Close(poller.pipe[1])
	***REMOVED***
	if poller.pipe[0] != -1 ***REMOVED***
		unix.Close(poller.pipe[0])
	***REMOVED***
	if poller.epfd != -1 ***REMOVED***
		unix.Close(poller.epfd)
	***REMOVED***
***REMOVED***
