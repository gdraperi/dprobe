// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux

package fsnotify

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"unsafe"

	"golang.org/x/sys/unix"
)

// Watcher watches a set of files, delivering events to a channel.
type Watcher struct ***REMOVED***
	Events   chan Event
	Errors   chan error
	mu       sync.Mutex // Map access
	fd       int
	poller   *fdPoller
	watches  map[string]*watch // Map of inotify watches (key: path)
	paths    map[int]string    // Map of watched paths (key: watch descriptor)
	done     chan struct***REMOVED******REMOVED***     // Channel for sending a "quit message" to the reader goroutine
	doneResp chan struct***REMOVED******REMOVED***     // Channel to respond to Close
***REMOVED***

// NewWatcher establishes a new watcher with the underlying OS and begins waiting for events.
func NewWatcher() (*Watcher, error) ***REMOVED***
	// Create inotify fd
	fd, errno := unix.InotifyInit1(unix.IN_CLOEXEC)
	if fd == -1 ***REMOVED***
		return nil, errno
	***REMOVED***
	// Create epoll
	poller, err := newFdPoller(fd)
	if err != nil ***REMOVED***
		unix.Close(fd)
		return nil, err
	***REMOVED***
	w := &Watcher***REMOVED***
		fd:       fd,
		poller:   poller,
		watches:  make(map[string]*watch),
		paths:    make(map[int]string),
		Events:   make(chan Event),
		Errors:   make(chan error),
		done:     make(chan struct***REMOVED******REMOVED***),
		doneResp: make(chan struct***REMOVED******REMOVED***),
	***REMOVED***

	go w.readEvents()
	return w, nil
***REMOVED***

func (w *Watcher) isClosed() bool ***REMOVED***
	select ***REMOVED***
	case <-w.done:
		return true
	default:
		return false
	***REMOVED***
***REMOVED***

// Close removes all watches and closes the events channel.
func (w *Watcher) Close() error ***REMOVED***
	if w.isClosed() ***REMOVED***
		return nil
	***REMOVED***

	// Send 'close' signal to goroutine, and set the Watcher to closed.
	close(w.done)

	// Wake up goroutine
	w.poller.wake()

	// Wait for goroutine to close
	<-w.doneResp

	return nil
***REMOVED***

// Add starts watching the named file or directory (non-recursively).
func (w *Watcher) Add(name string) error ***REMOVED***
	name = filepath.Clean(name)
	if w.isClosed() ***REMOVED***
		return errors.New("inotify instance already closed")
	***REMOVED***

	const agnosticEvents = unix.IN_MOVED_TO | unix.IN_MOVED_FROM |
		unix.IN_CREATE | unix.IN_ATTRIB | unix.IN_MODIFY |
		unix.IN_MOVE_SELF | unix.IN_DELETE | unix.IN_DELETE_SELF

	var flags uint32 = agnosticEvents

	w.mu.Lock()
	defer w.mu.Unlock()
	watchEntry := w.watches[name]
	if watchEntry != nil ***REMOVED***
		flags |= watchEntry.flags | unix.IN_MASK_ADD
	***REMOVED***
	wd, errno := unix.InotifyAddWatch(w.fd, name, flags)
	if wd == -1 ***REMOVED***
		return errno
	***REMOVED***

	if watchEntry == nil ***REMOVED***
		w.watches[name] = &watch***REMOVED***wd: uint32(wd), flags: flags***REMOVED***
		w.paths[wd] = name
	***REMOVED*** else ***REMOVED***
		watchEntry.wd = uint32(wd)
		watchEntry.flags = flags
	***REMOVED***

	return nil
***REMOVED***

// Remove stops watching the named file or directory (non-recursively).
func (w *Watcher) Remove(name string) error ***REMOVED***
	name = filepath.Clean(name)

	// Fetch the watch.
	w.mu.Lock()
	defer w.mu.Unlock()
	watch, ok := w.watches[name]

	// Remove it from inotify.
	if !ok ***REMOVED***
		return fmt.Errorf("can't remove non-existent inotify watch for: %s", name)
	***REMOVED***

	// We successfully removed the watch if InotifyRmWatch doesn't return an
	// error, we need to clean up our internal state to ensure it matches
	// inotify's kernel state.
	delete(w.paths, int(watch.wd))
	delete(w.watches, name)

	// inotify_rm_watch will return EINVAL if the file has been deleted;
	// the inotify will already have been removed.
	// watches and pathes are deleted in ignoreLinux() implicitly and asynchronously
	// by calling inotify_rm_watch() below. e.g. readEvents() goroutine receives IN_IGNORE
	// so that EINVAL means that the wd is being rm_watch()ed or its file removed
	// by another thread and we have not received IN_IGNORE event.
	success, errno := unix.InotifyRmWatch(w.fd, watch.wd)
	if success == -1 ***REMOVED***
		// TODO: Perhaps it's not helpful to return an error here in every case.
		// the only two possible errors are:
		// EBADF, which happens when w.fd is not a valid file descriptor of any kind.
		// EINVAL, which is when fd is not an inotify descriptor or wd is not a valid watch descriptor.
		// Watch descriptors are invalidated when they are removed explicitly or implicitly;
		// explicitly by inotify_rm_watch, implicitly when the file they are watching is deleted.
		return errno
	***REMOVED***

	return nil
***REMOVED***

type watch struct ***REMOVED***
	wd    uint32 // Watch descriptor (as returned by the inotify_add_watch() syscall)
	flags uint32 // inotify flags of this watch (see inotify(7) for the list of valid flags)
***REMOVED***

// readEvents reads from the inotify file descriptor, converts the
// received events into Event objects and sends them via the Events channel
func (w *Watcher) readEvents() ***REMOVED***
	var (
		buf   [unix.SizeofInotifyEvent * 4096]byte // Buffer for a maximum of 4096 raw events
		n     int                                  // Number of bytes read with read()
		errno error                                // Syscall errno
		ok    bool                                 // For poller.wait
	)

	defer close(w.doneResp)
	defer close(w.Errors)
	defer close(w.Events)
	defer unix.Close(w.fd)
	defer w.poller.close()

	for ***REMOVED***
		// See if we have been closed.
		if w.isClosed() ***REMOVED***
			return
		***REMOVED***

		ok, errno = w.poller.wait()
		if errno != nil ***REMOVED***
			select ***REMOVED***
			case w.Errors <- errno:
			case <-w.done:
				return
			***REMOVED***
			continue
		***REMOVED***

		if !ok ***REMOVED***
			continue
		***REMOVED***

		n, errno = unix.Read(w.fd, buf[:])
		// If a signal interrupted execution, see if we've been asked to close, and try again.
		// http://man7.org/linux/man-pages/man7/signal.7.html :
		// "Before Linux 3.8, reads from an inotify(7) file descriptor were not restartable"
		if errno == unix.EINTR ***REMOVED***
			continue
		***REMOVED***

		// unix.Read might have been woken up by Close. If so, we're done.
		if w.isClosed() ***REMOVED***
			return
		***REMOVED***

		if n < unix.SizeofInotifyEvent ***REMOVED***
			var err error
			if n == 0 ***REMOVED***
				// If EOF is received. This should really never happen.
				err = io.EOF
			***REMOVED*** else if n < 0 ***REMOVED***
				// If an error occurred while reading.
				err = errno
			***REMOVED*** else ***REMOVED***
				// Read was too short.
				err = errors.New("notify: short read in readEvents()")
			***REMOVED***
			select ***REMOVED***
			case w.Errors <- err:
			case <-w.done:
				return
			***REMOVED***
			continue
		***REMOVED***

		var offset uint32
		// We don't know how many events we just read into the buffer
		// While the offset points to at least one whole event...
		for offset <= uint32(n-unix.SizeofInotifyEvent) ***REMOVED***
			// Point "raw" to the event in the buffer
			raw := (*unix.InotifyEvent)(unsafe.Pointer(&buf[offset]))

			mask := uint32(raw.Mask)
			nameLen := uint32(raw.Len)

			if mask&unix.IN_Q_OVERFLOW != 0 ***REMOVED***
				select ***REMOVED***
				case w.Errors <- ErrEventOverflow:
				case <-w.done:
					return
				***REMOVED***
			***REMOVED***

			// If the event happened to the watched directory or the watched file, the kernel
			// doesn't append the filename to the event, but we would like to always fill the
			// the "Name" field with a valid filename. We retrieve the path of the watch from
			// the "paths" map.
			w.mu.Lock()
			name, ok := w.paths[int(raw.Wd)]
			// IN_DELETE_SELF occurs when the file/directory being watched is removed.
			// This is a sign to clean up the maps, otherwise we are no longer in sync
			// with the inotify kernel state which has already deleted the watch
			// automatically.
			if ok && mask&unix.IN_DELETE_SELF == unix.IN_DELETE_SELF ***REMOVED***
				delete(w.paths, int(raw.Wd))
				delete(w.watches, name)
			***REMOVED***
			w.mu.Unlock()

			if nameLen > 0 ***REMOVED***
				// Point "bytes" at the first byte of the filename
				bytes := (*[unix.PathMax]byte)(unsafe.Pointer(&buf[offset+unix.SizeofInotifyEvent]))
				// The filename is padded with NULL bytes. TrimRight() gets rid of those.
				name += "/" + strings.TrimRight(string(bytes[0:nameLen]), "\000")
			***REMOVED***

			event := newEvent(name, mask)

			// Send the events that are not ignored on the events channel
			if !event.ignoreLinux(mask) ***REMOVED***
				select ***REMOVED***
				case w.Events <- event:
				case <-w.done:
					return
				***REMOVED***
			***REMOVED***

			// Move to the next event in the buffer
			offset += unix.SizeofInotifyEvent + nameLen
		***REMOVED***
	***REMOVED***
***REMOVED***

// Certain types of events can be "ignored" and not sent over the Events
// channel. Such as events marked ignore by the kernel, or MODIFY events
// against files that do not exist.
func (e *Event) ignoreLinux(mask uint32) bool ***REMOVED***
	// Ignore anything the inotify API says to ignore
	if mask&unix.IN_IGNORED == unix.IN_IGNORED ***REMOVED***
		return true
	***REMOVED***

	// If the event is not a DELETE or RENAME, the file must exist.
	// Otherwise the event is ignored.
	// *Note*: this was put in place because it was seen that a MODIFY
	// event was sent after the DELETE. This ignores that MODIFY and
	// assumes a DELETE will come or has come if the file doesn't exist.
	if !(e.Op&Remove == Remove || e.Op&Rename == Rename) ***REMOVED***
		_, statErr := os.Lstat(e.Name)
		return os.IsNotExist(statErr)
	***REMOVED***
	return false
***REMOVED***

// newEvent returns an platform-independent Event based on an inotify mask.
func newEvent(name string, mask uint32) Event ***REMOVED***
	e := Event***REMOVED***Name: name***REMOVED***
	if mask&unix.IN_CREATE == unix.IN_CREATE || mask&unix.IN_MOVED_TO == unix.IN_MOVED_TO ***REMOVED***
		e.Op |= Create
	***REMOVED***
	if mask&unix.IN_DELETE_SELF == unix.IN_DELETE_SELF || mask&unix.IN_DELETE == unix.IN_DELETE ***REMOVED***
		e.Op |= Remove
	***REMOVED***
	if mask&unix.IN_MODIFY == unix.IN_MODIFY ***REMOVED***
		e.Op |= Write
	***REMOVED***
	if mask&unix.IN_MOVE_SELF == unix.IN_MOVE_SELF || mask&unix.IN_MOVED_FROM == unix.IN_MOVED_FROM ***REMOVED***
		e.Op |= Rename
	***REMOVED***
	if mask&unix.IN_ATTRIB == unix.IN_ATTRIB ***REMOVED***
		e.Op |= Chmod
	***REMOVED***
	return e
***REMOVED***
