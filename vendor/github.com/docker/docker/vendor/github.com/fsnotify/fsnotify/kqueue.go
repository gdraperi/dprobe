// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build freebsd openbsd netbsd dragonfly darwin

package fsnotify

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"golang.org/x/sys/unix"
)

// Watcher watches a set of files, delivering events to a channel.
type Watcher struct ***REMOVED***
	Events chan Event
	Errors chan error
	done   chan bool // Channel for sending a "quit message" to the reader goroutine

	kq int // File descriptor (as returned by the kqueue() syscall).

	mu              sync.Mutex        // Protects access to watcher data
	watches         map[string]int    // Map of watched file descriptors (key: path).
	externalWatches map[string]bool   // Map of watches added by user of the library.
	dirFlags        map[string]uint32 // Map of watched directories to fflags used in kqueue.
	paths           map[int]pathInfo  // Map file descriptors to path names for processing kqueue events.
	fileExists      map[string]bool   // Keep track of if we know this file exists (to stop duplicate create events).
	isClosed        bool              // Set to true when Close() is first called
***REMOVED***

type pathInfo struct ***REMOVED***
	name  string
	isDir bool
***REMOVED***

// NewWatcher establishes a new watcher with the underlying OS and begins waiting for events.
func NewWatcher() (*Watcher, error) ***REMOVED***
	kq, err := kqueue()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	w := &Watcher***REMOVED***
		kq:              kq,
		watches:         make(map[string]int),
		dirFlags:        make(map[string]uint32),
		paths:           make(map[int]pathInfo),
		fileExists:      make(map[string]bool),
		externalWatches: make(map[string]bool),
		Events:          make(chan Event),
		Errors:          make(chan error),
		done:            make(chan bool),
	***REMOVED***

	go w.readEvents()
	return w, nil
***REMOVED***

// Close removes all watches and closes the events channel.
func (w *Watcher) Close() error ***REMOVED***
	w.mu.Lock()
	if w.isClosed ***REMOVED***
		w.mu.Unlock()
		return nil
	***REMOVED***
	w.isClosed = true
	w.mu.Unlock()

	// copy paths to remove while locked
	w.mu.Lock()
	var pathsToRemove = make([]string, 0, len(w.watches))
	for name := range w.watches ***REMOVED***
		pathsToRemove = append(pathsToRemove, name)
	***REMOVED***
	w.mu.Unlock()
	// unlock before calling Remove, which also locks

	var err error
	for _, name := range pathsToRemove ***REMOVED***
		if e := w.Remove(name); e != nil && err == nil ***REMOVED***
			err = e
		***REMOVED***
	***REMOVED***

	// Send "quit" message to the reader goroutine:
	w.done <- true

	return nil
***REMOVED***

// Add starts watching the named file or directory (non-recursively).
func (w *Watcher) Add(name string) error ***REMOVED***
	w.mu.Lock()
	w.externalWatches[name] = true
	w.mu.Unlock()
	_, err := w.addWatch(name, noteAllEvents)
	return err
***REMOVED***

// Remove stops watching the the named file or directory (non-recursively).
func (w *Watcher) Remove(name string) error ***REMOVED***
	name = filepath.Clean(name)
	w.mu.Lock()
	watchfd, ok := w.watches[name]
	w.mu.Unlock()
	if !ok ***REMOVED***
		return fmt.Errorf("can't remove non-existent kevent watch for: %s", name)
	***REMOVED***

	const registerRemove = unix.EV_DELETE
	if err := register(w.kq, []int***REMOVED***watchfd***REMOVED***, registerRemove, 0); err != nil ***REMOVED***
		return err
	***REMOVED***

	unix.Close(watchfd)

	w.mu.Lock()
	isDir := w.paths[watchfd].isDir
	delete(w.watches, name)
	delete(w.paths, watchfd)
	delete(w.dirFlags, name)
	w.mu.Unlock()

	// Find all watched paths that are in this directory that are not external.
	if isDir ***REMOVED***
		var pathsToRemove []string
		w.mu.Lock()
		for _, path := range w.paths ***REMOVED***
			wdir, _ := filepath.Split(path.name)
			if filepath.Clean(wdir) == name ***REMOVED***
				if !w.externalWatches[path.name] ***REMOVED***
					pathsToRemove = append(pathsToRemove, path.name)
				***REMOVED***
			***REMOVED***
		***REMOVED***
		w.mu.Unlock()
		for _, name := range pathsToRemove ***REMOVED***
			// Since these are internal, not much sense in propagating error
			// to the user, as that will just confuse them with an error about
			// a path they did not explicitly watch themselves.
			w.Remove(name)
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

// Watch all events (except NOTE_EXTEND, NOTE_LINK, NOTE_REVOKE)
const noteAllEvents = unix.NOTE_DELETE | unix.NOTE_WRITE | unix.NOTE_ATTRIB | unix.NOTE_RENAME

// keventWaitTime to block on each read from kevent
var keventWaitTime = durationToTimespec(100 * time.Millisecond)

// addWatch adds name to the watched file set.
// The flags are interpreted as described in kevent(2).
// Returns the real path to the file which was added, if any, which may be different from the one passed in the case of symlinks.
func (w *Watcher) addWatch(name string, flags uint32) (string, error) ***REMOVED***
	var isDir bool
	// Make ./name and name equivalent
	name = filepath.Clean(name)

	w.mu.Lock()
	if w.isClosed ***REMOVED***
		w.mu.Unlock()
		return "", errors.New("kevent instance already closed")
	***REMOVED***
	watchfd, alreadyWatching := w.watches[name]
	// We already have a watch, but we can still override flags.
	if alreadyWatching ***REMOVED***
		isDir = w.paths[watchfd].isDir
	***REMOVED***
	w.mu.Unlock()

	if !alreadyWatching ***REMOVED***
		fi, err := os.Lstat(name)
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***

		// Don't watch sockets.
		if fi.Mode()&os.ModeSocket == os.ModeSocket ***REMOVED***
			return "", nil
		***REMOVED***

		// Don't watch named pipes.
		if fi.Mode()&os.ModeNamedPipe == os.ModeNamedPipe ***REMOVED***
			return "", nil
		***REMOVED***

		// Follow Symlinks
		// Unfortunately, Linux can add bogus symlinks to watch list without
		// issue, and Windows can't do symlinks period (AFAIK). To  maintain
		// consistency, we will act like everything is fine. There will simply
		// be no file events for broken symlinks.
		// Hence the returns of nil on errors.
		if fi.Mode()&os.ModeSymlink == os.ModeSymlink ***REMOVED***
			name, err = filepath.EvalSymlinks(name)
			if err != nil ***REMOVED***
				return "", nil
			***REMOVED***

			w.mu.Lock()
			_, alreadyWatching = w.watches[name]
			w.mu.Unlock()

			if alreadyWatching ***REMOVED***
				return name, nil
			***REMOVED***

			fi, err = os.Lstat(name)
			if err != nil ***REMOVED***
				return "", nil
			***REMOVED***
		***REMOVED***

		watchfd, err = unix.Open(name, openMode, 0700)
		if watchfd == -1 ***REMOVED***
			return "", err
		***REMOVED***

		isDir = fi.IsDir()
	***REMOVED***

	const registerAdd = unix.EV_ADD | unix.EV_CLEAR | unix.EV_ENABLE
	if err := register(w.kq, []int***REMOVED***watchfd***REMOVED***, registerAdd, flags); err != nil ***REMOVED***
		unix.Close(watchfd)
		return "", err
	***REMOVED***

	if !alreadyWatching ***REMOVED***
		w.mu.Lock()
		w.watches[name] = watchfd
		w.paths[watchfd] = pathInfo***REMOVED***name: name, isDir: isDir***REMOVED***
		w.mu.Unlock()
	***REMOVED***

	if isDir ***REMOVED***
		// Watch the directory if it has not been watched before,
		// or if it was watched before, but perhaps only a NOTE_DELETE (watchDirectoryFiles)
		w.mu.Lock()

		watchDir := (flags&unix.NOTE_WRITE) == unix.NOTE_WRITE &&
			(!alreadyWatching || (w.dirFlags[name]&unix.NOTE_WRITE) != unix.NOTE_WRITE)
		// Store flags so this watch can be updated later
		w.dirFlags[name] = flags
		w.mu.Unlock()

		if watchDir ***REMOVED***
			if err := w.watchDirectoryFiles(name); err != nil ***REMOVED***
				return "", err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return name, nil
***REMOVED***

// readEvents reads from kqueue and converts the received kevents into
// Event values that it sends down the Events channel.
func (w *Watcher) readEvents() ***REMOVED***
	eventBuffer := make([]unix.Kevent_t, 10)

	for ***REMOVED***
		// See if there is a message on the "done" channel
		select ***REMOVED***
		case <-w.done:
			err := unix.Close(w.kq)
			if err != nil ***REMOVED***
				w.Errors <- err
			***REMOVED***
			close(w.Events)
			close(w.Errors)
			return
		default:
		***REMOVED***

		// Get new events
		kevents, err := read(w.kq, eventBuffer, &keventWaitTime)
		// EINTR is okay, the syscall was interrupted before timeout expired.
		if err != nil && err != unix.EINTR ***REMOVED***
			w.Errors <- err
			continue
		***REMOVED***

		// Flush the events we received to the Events channel
		for len(kevents) > 0 ***REMOVED***
			kevent := &kevents[0]
			watchfd := int(kevent.Ident)
			mask := uint32(kevent.Fflags)
			w.mu.Lock()
			path := w.paths[watchfd]
			w.mu.Unlock()
			event := newEvent(path.name, mask)

			if path.isDir && !(event.Op&Remove == Remove) ***REMOVED***
				// Double check to make sure the directory exists. This can happen when
				// we do a rm -fr on a recursively watched folders and we receive a
				// modification event first but the folder has been deleted and later
				// receive the delete event
				if _, err := os.Lstat(event.Name); os.IsNotExist(err) ***REMOVED***
					// mark is as delete event
					event.Op |= Remove
				***REMOVED***
			***REMOVED***

			if event.Op&Rename == Rename || event.Op&Remove == Remove ***REMOVED***
				w.Remove(event.Name)
				w.mu.Lock()
				delete(w.fileExists, event.Name)
				w.mu.Unlock()
			***REMOVED***

			if path.isDir && event.Op&Write == Write && !(event.Op&Remove == Remove) ***REMOVED***
				w.sendDirectoryChangeEvents(event.Name)
			***REMOVED*** else ***REMOVED***
				// Send the event on the Events channel
				w.Events <- event
			***REMOVED***

			if event.Op&Remove == Remove ***REMOVED***
				// Look for a file that may have overwritten this.
				// For example, mv f1 f2 will delete f2, then create f2.
				if path.isDir ***REMOVED***
					fileDir := filepath.Clean(event.Name)
					w.mu.Lock()
					_, found := w.watches[fileDir]
					w.mu.Unlock()
					if found ***REMOVED***
						// make sure the directory exists before we watch for changes. When we
						// do a recursive watch and perform rm -fr, the parent directory might
						// have gone missing, ignore the missing directory and let the
						// upcoming delete event remove the watch from the parent directory.
						if _, err := os.Lstat(fileDir); err == nil ***REMOVED***
							w.sendDirectoryChangeEvents(fileDir)
						***REMOVED***
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					filePath := filepath.Clean(event.Name)
					if fileInfo, err := os.Lstat(filePath); err == nil ***REMOVED***
						w.sendFileCreatedEventIfNew(filePath, fileInfo)
					***REMOVED***
				***REMOVED***
			***REMOVED***

			// Move to next event
			kevents = kevents[1:]
		***REMOVED***
	***REMOVED***
***REMOVED***

// newEvent returns an platform-independent Event based on kqueue Fflags.
func newEvent(name string, mask uint32) Event ***REMOVED***
	e := Event***REMOVED***Name: name***REMOVED***
	if mask&unix.NOTE_DELETE == unix.NOTE_DELETE ***REMOVED***
		e.Op |= Remove
	***REMOVED***
	if mask&unix.NOTE_WRITE == unix.NOTE_WRITE ***REMOVED***
		e.Op |= Write
	***REMOVED***
	if mask&unix.NOTE_RENAME == unix.NOTE_RENAME ***REMOVED***
		e.Op |= Rename
	***REMOVED***
	if mask&unix.NOTE_ATTRIB == unix.NOTE_ATTRIB ***REMOVED***
		e.Op |= Chmod
	***REMOVED***
	return e
***REMOVED***

func newCreateEvent(name string) Event ***REMOVED***
	return Event***REMOVED***Name: name, Op: Create***REMOVED***
***REMOVED***

// watchDirectoryFiles to mimic inotify when adding a watch on a directory
func (w *Watcher) watchDirectoryFiles(dirPath string) error ***REMOVED***
	// Get all files
	files, err := ioutil.ReadDir(dirPath)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, fileInfo := range files ***REMOVED***
		filePath := filepath.Join(dirPath, fileInfo.Name())
		filePath, err = w.internalWatch(filePath, fileInfo)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		w.mu.Lock()
		w.fileExists[filePath] = true
		w.mu.Unlock()
	***REMOVED***

	return nil
***REMOVED***

// sendDirectoryEvents searches the directory for newly created files
// and sends them over the event channel. This functionality is to have
// the BSD version of fsnotify match Linux inotify which provides a
// create event for files created in a watched directory.
func (w *Watcher) sendDirectoryChangeEvents(dirPath string) ***REMOVED***
	// Get all files
	files, err := ioutil.ReadDir(dirPath)
	if err != nil ***REMOVED***
		w.Errors <- err
	***REMOVED***

	// Search for new files
	for _, fileInfo := range files ***REMOVED***
		filePath := filepath.Join(dirPath, fileInfo.Name())
		err := w.sendFileCreatedEventIfNew(filePath, fileInfo)

		if err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// sendFileCreatedEvent sends a create event if the file isn't already being tracked.
func (w *Watcher) sendFileCreatedEventIfNew(filePath string, fileInfo os.FileInfo) (err error) ***REMOVED***
	w.mu.Lock()
	_, doesExist := w.fileExists[filePath]
	w.mu.Unlock()
	if !doesExist ***REMOVED***
		// Send create event
		w.Events <- newCreateEvent(filePath)
	***REMOVED***

	// like watchDirectoryFiles (but without doing another ReadDir)
	filePath, err = w.internalWatch(filePath, fileInfo)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	w.mu.Lock()
	w.fileExists[filePath] = true
	w.mu.Unlock()

	return nil
***REMOVED***

func (w *Watcher) internalWatch(name string, fileInfo os.FileInfo) (string, error) ***REMOVED***
	if fileInfo.IsDir() ***REMOVED***
		// mimic Linux providing delete events for subdirectories
		// but preserve the flags used if currently watching subdirectory
		w.mu.Lock()
		flags := w.dirFlags[name]
		w.mu.Unlock()

		flags |= unix.NOTE_DELETE | unix.NOTE_RENAME
		return w.addWatch(name, flags)
	***REMOVED***

	// watch file to mimic Linux inotify
	return w.addWatch(name, noteAllEvents)
***REMOVED***

// kqueue creates a new kernel event queue and returns a descriptor.
func kqueue() (kq int, err error) ***REMOVED***
	kq, err = unix.Kqueue()
	if kq == -1 ***REMOVED***
		return kq, err
	***REMOVED***
	return kq, nil
***REMOVED***

// register events with the queue
func register(kq int, fds []int, flags int, fflags uint32) error ***REMOVED***
	changes := make([]unix.Kevent_t, len(fds))

	for i, fd := range fds ***REMOVED***
		// SetKevent converts int to the platform-specific types:
		unix.SetKevent(&changes[i], fd, unix.EVFILT_VNODE, flags)
		changes[i].Fflags = fflags
	***REMOVED***

	// register the events
	success, err := unix.Kevent(kq, changes, nil, nil)
	if success == -1 ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// read retrieves pending events, or waits until an event occurs.
// A timeout of nil blocks indefinitely, while 0 polls the queue.
func read(kq int, events []unix.Kevent_t, timeout *unix.Timespec) ([]unix.Kevent_t, error) ***REMOVED***
	n, err := unix.Kevent(kq, nil, events, timeout)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return events[0:n], nil
***REMOVED***

// durationToTimespec prepares a timeout value
func durationToTimespec(d time.Duration) unix.Timespec ***REMOVED***
	return unix.NsecToTimespec(d.Nanoseconds())
***REMOVED***
