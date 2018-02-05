// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package fsnotify

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"
	"unsafe"
)

// Watcher watches a set of files, delivering events to a channel.
type Watcher struct ***REMOVED***
	Events   chan Event
	Errors   chan error
	isClosed bool           // Set to true when Close() is first called
	mu       sync.Mutex     // Map access
	port     syscall.Handle // Handle to completion port
	watches  watchMap       // Map of watches (key: i-number)
	input    chan *input    // Inputs to the reader are sent on this channel
	quit     chan chan<- error
***REMOVED***

// NewWatcher establishes a new watcher with the underlying OS and begins waiting for events.
func NewWatcher() (*Watcher, error) ***REMOVED***
	port, e := syscall.CreateIoCompletionPort(syscall.InvalidHandle, 0, 0, 0)
	if e != nil ***REMOVED***
		return nil, os.NewSyscallError("CreateIoCompletionPort", e)
	***REMOVED***
	w := &Watcher***REMOVED***
		port:    port,
		watches: make(watchMap),
		input:   make(chan *input, 1),
		Events:  make(chan Event, 50),
		Errors:  make(chan error),
		quit:    make(chan chan<- error, 1),
	***REMOVED***
	go w.readEvents()
	return w, nil
***REMOVED***

// Close removes all watches and closes the events channel.
func (w *Watcher) Close() error ***REMOVED***
	if w.isClosed ***REMOVED***
		return nil
	***REMOVED***
	w.isClosed = true

	// Send "quit" message to the reader goroutine
	ch := make(chan error)
	w.quit <- ch
	if err := w.wakeupReader(); err != nil ***REMOVED***
		return err
	***REMOVED***
	return <-ch
***REMOVED***

// Add starts watching the named file or directory (non-recursively).
func (w *Watcher) Add(name string) error ***REMOVED***
	if w.isClosed ***REMOVED***
		return errors.New("watcher already closed")
	***REMOVED***
	in := &input***REMOVED***
		op:    opAddWatch,
		path:  filepath.Clean(name),
		flags: sysFSALLEVENTS,
		reply: make(chan error),
	***REMOVED***
	w.input <- in
	if err := w.wakeupReader(); err != nil ***REMOVED***
		return err
	***REMOVED***
	return <-in.reply
***REMOVED***

// Remove stops watching the the named file or directory (non-recursively).
func (w *Watcher) Remove(name string) error ***REMOVED***
	in := &input***REMOVED***
		op:    opRemoveWatch,
		path:  filepath.Clean(name),
		reply: make(chan error),
	***REMOVED***
	w.input <- in
	if err := w.wakeupReader(); err != nil ***REMOVED***
		return err
	***REMOVED***
	return <-in.reply
***REMOVED***

const (
	// Options for AddWatch
	sysFSONESHOT = 0x80000000
	sysFSONLYDIR = 0x1000000

	// Events
	sysFSACCESS     = 0x1
	sysFSALLEVENTS  = 0xfff
	sysFSATTRIB     = 0x4
	sysFSCLOSE      = 0x18
	sysFSCREATE     = 0x100
	sysFSDELETE     = 0x200
	sysFSDELETESELF = 0x400
	sysFSMODIFY     = 0x2
	sysFSMOVE       = 0xc0
	sysFSMOVEDFROM  = 0x40
	sysFSMOVEDTO    = 0x80
	sysFSMOVESELF   = 0x800

	// Special events
	sysFSIGNORED   = 0x8000
	sysFSQOVERFLOW = 0x4000
)

func newEvent(name string, mask uint32) Event ***REMOVED***
	e := Event***REMOVED***Name: name***REMOVED***
	if mask&sysFSCREATE == sysFSCREATE || mask&sysFSMOVEDTO == sysFSMOVEDTO ***REMOVED***
		e.Op |= Create
	***REMOVED***
	if mask&sysFSDELETE == sysFSDELETE || mask&sysFSDELETESELF == sysFSDELETESELF ***REMOVED***
		e.Op |= Remove
	***REMOVED***
	if mask&sysFSMODIFY == sysFSMODIFY ***REMOVED***
		e.Op |= Write
	***REMOVED***
	if mask&sysFSMOVE == sysFSMOVE || mask&sysFSMOVESELF == sysFSMOVESELF || mask&sysFSMOVEDFROM == sysFSMOVEDFROM ***REMOVED***
		e.Op |= Rename
	***REMOVED***
	if mask&sysFSATTRIB == sysFSATTRIB ***REMOVED***
		e.Op |= Chmod
	***REMOVED***
	return e
***REMOVED***

const (
	opAddWatch = iota
	opRemoveWatch
)

const (
	provisional uint64 = 1 << (32 + iota)
)

type input struct ***REMOVED***
	op    int
	path  string
	flags uint32
	reply chan error
***REMOVED***

type inode struct ***REMOVED***
	handle syscall.Handle
	volume uint32
	index  uint64
***REMOVED***

type watch struct ***REMOVED***
	ov     syscall.Overlapped
	ino    *inode            // i-number
	path   string            // Directory path
	mask   uint64            // Directory itself is being watched with these notify flags
	names  map[string]uint64 // Map of names being watched and their notify flags
	rename string            // Remembers the old name while renaming a file
	buf    [4096]byte
***REMOVED***

type indexMap map[uint64]*watch
type watchMap map[uint32]indexMap

func (w *Watcher) wakeupReader() error ***REMOVED***
	e := syscall.PostQueuedCompletionStatus(w.port, 0, 0, nil)
	if e != nil ***REMOVED***
		return os.NewSyscallError("PostQueuedCompletionStatus", e)
	***REMOVED***
	return nil
***REMOVED***

func getDir(pathname string) (dir string, err error) ***REMOVED***
	attr, e := syscall.GetFileAttributes(syscall.StringToUTF16Ptr(pathname))
	if e != nil ***REMOVED***
		return "", os.NewSyscallError("GetFileAttributes", e)
	***REMOVED***
	if attr&syscall.FILE_ATTRIBUTE_DIRECTORY != 0 ***REMOVED***
		dir = pathname
	***REMOVED*** else ***REMOVED***
		dir, _ = filepath.Split(pathname)
		dir = filepath.Clean(dir)
	***REMOVED***
	return
***REMOVED***

func getIno(path string) (ino *inode, err error) ***REMOVED***
	h, e := syscall.CreateFile(syscall.StringToUTF16Ptr(path),
		syscall.FILE_LIST_DIRECTORY,
		syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE|syscall.FILE_SHARE_DELETE,
		nil, syscall.OPEN_EXISTING,
		syscall.FILE_FLAG_BACKUP_SEMANTICS|syscall.FILE_FLAG_OVERLAPPED, 0)
	if e != nil ***REMOVED***
		return nil, os.NewSyscallError("CreateFile", e)
	***REMOVED***
	var fi syscall.ByHandleFileInformation
	if e = syscall.GetFileInformationByHandle(h, &fi); e != nil ***REMOVED***
		syscall.CloseHandle(h)
		return nil, os.NewSyscallError("GetFileInformationByHandle", e)
	***REMOVED***
	ino = &inode***REMOVED***
		handle: h,
		volume: fi.VolumeSerialNumber,
		index:  uint64(fi.FileIndexHigh)<<32 | uint64(fi.FileIndexLow),
	***REMOVED***
	return ino, nil
***REMOVED***

// Must run within the I/O thread.
func (m watchMap) get(ino *inode) *watch ***REMOVED***
	if i := m[ino.volume]; i != nil ***REMOVED***
		return i[ino.index]
	***REMOVED***
	return nil
***REMOVED***

// Must run within the I/O thread.
func (m watchMap) set(ino *inode, watch *watch) ***REMOVED***
	i := m[ino.volume]
	if i == nil ***REMOVED***
		i = make(indexMap)
		m[ino.volume] = i
	***REMOVED***
	i[ino.index] = watch
***REMOVED***

// Must run within the I/O thread.
func (w *Watcher) addWatch(pathname string, flags uint64) error ***REMOVED***
	dir, err := getDir(pathname)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if flags&sysFSONLYDIR != 0 && pathname != dir ***REMOVED***
		return nil
	***REMOVED***
	ino, err := getIno(dir)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	w.mu.Lock()
	watchEntry := w.watches.get(ino)
	w.mu.Unlock()
	if watchEntry == nil ***REMOVED***
		if _, e := syscall.CreateIoCompletionPort(ino.handle, w.port, 0, 0); e != nil ***REMOVED***
			syscall.CloseHandle(ino.handle)
			return os.NewSyscallError("CreateIoCompletionPort", e)
		***REMOVED***
		watchEntry = &watch***REMOVED***
			ino:   ino,
			path:  dir,
			names: make(map[string]uint64),
		***REMOVED***
		w.mu.Lock()
		w.watches.set(ino, watchEntry)
		w.mu.Unlock()
		flags |= provisional
	***REMOVED*** else ***REMOVED***
		syscall.CloseHandle(ino.handle)
	***REMOVED***
	if pathname == dir ***REMOVED***
		watchEntry.mask |= flags
	***REMOVED*** else ***REMOVED***
		watchEntry.names[filepath.Base(pathname)] |= flags
	***REMOVED***
	if err = w.startRead(watchEntry); err != nil ***REMOVED***
		return err
	***REMOVED***
	if pathname == dir ***REMOVED***
		watchEntry.mask &= ^provisional
	***REMOVED*** else ***REMOVED***
		watchEntry.names[filepath.Base(pathname)] &= ^provisional
	***REMOVED***
	return nil
***REMOVED***

// Must run within the I/O thread.
func (w *Watcher) remWatch(pathname string) error ***REMOVED***
	dir, err := getDir(pathname)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	ino, err := getIno(dir)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	w.mu.Lock()
	watch := w.watches.get(ino)
	w.mu.Unlock()
	if watch == nil ***REMOVED***
		return fmt.Errorf("can't remove non-existent watch for: %s", pathname)
	***REMOVED***
	if pathname == dir ***REMOVED***
		w.sendEvent(watch.path, watch.mask&sysFSIGNORED)
		watch.mask = 0
	***REMOVED*** else ***REMOVED***
		name := filepath.Base(pathname)
		w.sendEvent(filepath.Join(watch.path, name), watch.names[name]&sysFSIGNORED)
		delete(watch.names, name)
	***REMOVED***
	return w.startRead(watch)
***REMOVED***

// Must run within the I/O thread.
func (w *Watcher) deleteWatch(watch *watch) ***REMOVED***
	for name, mask := range watch.names ***REMOVED***
		if mask&provisional == 0 ***REMOVED***
			w.sendEvent(filepath.Join(watch.path, name), mask&sysFSIGNORED)
		***REMOVED***
		delete(watch.names, name)
	***REMOVED***
	if watch.mask != 0 ***REMOVED***
		if watch.mask&provisional == 0 ***REMOVED***
			w.sendEvent(watch.path, watch.mask&sysFSIGNORED)
		***REMOVED***
		watch.mask = 0
	***REMOVED***
***REMOVED***

// Must run within the I/O thread.
func (w *Watcher) startRead(watch *watch) error ***REMOVED***
	if e := syscall.CancelIo(watch.ino.handle); e != nil ***REMOVED***
		w.Errors <- os.NewSyscallError("CancelIo", e)
		w.deleteWatch(watch)
	***REMOVED***
	mask := toWindowsFlags(watch.mask)
	for _, m := range watch.names ***REMOVED***
		mask |= toWindowsFlags(m)
	***REMOVED***
	if mask == 0 ***REMOVED***
		if e := syscall.CloseHandle(watch.ino.handle); e != nil ***REMOVED***
			w.Errors <- os.NewSyscallError("CloseHandle", e)
		***REMOVED***
		w.mu.Lock()
		delete(w.watches[watch.ino.volume], watch.ino.index)
		w.mu.Unlock()
		return nil
	***REMOVED***
	e := syscall.ReadDirectoryChanges(watch.ino.handle, &watch.buf[0],
		uint32(unsafe.Sizeof(watch.buf)), false, mask, nil, &watch.ov, 0)
	if e != nil ***REMOVED***
		err := os.NewSyscallError("ReadDirectoryChanges", e)
		if e == syscall.ERROR_ACCESS_DENIED && watch.mask&provisional == 0 ***REMOVED***
			// Watched directory was probably removed
			if w.sendEvent(watch.path, watch.mask&sysFSDELETESELF) ***REMOVED***
				if watch.mask&sysFSONESHOT != 0 ***REMOVED***
					watch.mask = 0
				***REMOVED***
			***REMOVED***
			err = nil
		***REMOVED***
		w.deleteWatch(watch)
		w.startRead(watch)
		return err
	***REMOVED***
	return nil
***REMOVED***

// readEvents reads from the I/O completion port, converts the
// received events into Event objects and sends them via the Events channel.
// Entry point to the I/O thread.
func (w *Watcher) readEvents() ***REMOVED***
	var (
		n, key uint32
		ov     *syscall.Overlapped
	)
	runtime.LockOSThread()

	for ***REMOVED***
		e := syscall.GetQueuedCompletionStatus(w.port, &n, &key, &ov, syscall.INFINITE)
		watch := (*watch)(unsafe.Pointer(ov))

		if watch == nil ***REMOVED***
			select ***REMOVED***
			case ch := <-w.quit:
				w.mu.Lock()
				var indexes []indexMap
				for _, index := range w.watches ***REMOVED***
					indexes = append(indexes, index)
				***REMOVED***
				w.mu.Unlock()
				for _, index := range indexes ***REMOVED***
					for _, watch := range index ***REMOVED***
						w.deleteWatch(watch)
						w.startRead(watch)
					***REMOVED***
				***REMOVED***
				var err error
				if e := syscall.CloseHandle(w.port); e != nil ***REMOVED***
					err = os.NewSyscallError("CloseHandle", e)
				***REMOVED***
				close(w.Events)
				close(w.Errors)
				ch <- err
				return
			case in := <-w.input:
				switch in.op ***REMOVED***
				case opAddWatch:
					in.reply <- w.addWatch(in.path, uint64(in.flags))
				case opRemoveWatch:
					in.reply <- w.remWatch(in.path)
				***REMOVED***
			default:
			***REMOVED***
			continue
		***REMOVED***

		switch e ***REMOVED***
		case syscall.ERROR_MORE_DATA:
			if watch == nil ***REMOVED***
				w.Errors <- errors.New("ERROR_MORE_DATA has unexpectedly null lpOverlapped buffer")
			***REMOVED*** else ***REMOVED***
				// The i/o succeeded but the buffer is full.
				// In theory we should be building up a full packet.
				// In practice we can get away with just carrying on.
				n = uint32(unsafe.Sizeof(watch.buf))
			***REMOVED***
		case syscall.ERROR_ACCESS_DENIED:
			// Watched directory was probably removed
			w.sendEvent(watch.path, watch.mask&sysFSDELETESELF)
			w.deleteWatch(watch)
			w.startRead(watch)
			continue
		case syscall.ERROR_OPERATION_ABORTED:
			// CancelIo was called on this handle
			continue
		default:
			w.Errors <- os.NewSyscallError("GetQueuedCompletionPort", e)
			continue
		case nil:
		***REMOVED***

		var offset uint32
		for ***REMOVED***
			if n == 0 ***REMOVED***
				w.Events <- newEvent("", sysFSQOVERFLOW)
				w.Errors <- errors.New("short read in readEvents()")
				break
			***REMOVED***

			// Point "raw" to the event in the buffer
			raw := (*syscall.FileNotifyInformation)(unsafe.Pointer(&watch.buf[offset]))
			buf := (*[syscall.MAX_PATH]uint16)(unsafe.Pointer(&raw.FileName))
			name := syscall.UTF16ToString(buf[:raw.FileNameLength/2])
			fullname := filepath.Join(watch.path, name)

			var mask uint64
			switch raw.Action ***REMOVED***
			case syscall.FILE_ACTION_REMOVED:
				mask = sysFSDELETESELF
			case syscall.FILE_ACTION_MODIFIED:
				mask = sysFSMODIFY
			case syscall.FILE_ACTION_RENAMED_OLD_NAME:
				watch.rename = name
			case syscall.FILE_ACTION_RENAMED_NEW_NAME:
				if watch.names[watch.rename] != 0 ***REMOVED***
					watch.names[name] |= watch.names[watch.rename]
					delete(watch.names, watch.rename)
					mask = sysFSMOVESELF
				***REMOVED***
			***REMOVED***

			sendNameEvent := func() ***REMOVED***
				if w.sendEvent(fullname, watch.names[name]&mask) ***REMOVED***
					if watch.names[name]&sysFSONESHOT != 0 ***REMOVED***
						delete(watch.names, name)
					***REMOVED***
				***REMOVED***
			***REMOVED***
			if raw.Action != syscall.FILE_ACTION_RENAMED_NEW_NAME ***REMOVED***
				sendNameEvent()
			***REMOVED***
			if raw.Action == syscall.FILE_ACTION_REMOVED ***REMOVED***
				w.sendEvent(fullname, watch.names[name]&sysFSIGNORED)
				delete(watch.names, name)
			***REMOVED***
			if w.sendEvent(fullname, watch.mask&toFSnotifyFlags(raw.Action)) ***REMOVED***
				if watch.mask&sysFSONESHOT != 0 ***REMOVED***
					watch.mask = 0
				***REMOVED***
			***REMOVED***
			if raw.Action == syscall.FILE_ACTION_RENAMED_NEW_NAME ***REMOVED***
				fullname = filepath.Join(watch.path, watch.rename)
				sendNameEvent()
			***REMOVED***

			// Move to the next event in the buffer
			if raw.NextEntryOffset == 0 ***REMOVED***
				break
			***REMOVED***
			offset += raw.NextEntryOffset

			// Error!
			if offset >= n ***REMOVED***
				w.Errors <- errors.New("Windows system assumed buffer larger than it is, events have likely been missed.")
				break
			***REMOVED***
		***REMOVED***

		if err := w.startRead(watch); err != nil ***REMOVED***
			w.Errors <- err
		***REMOVED***
	***REMOVED***
***REMOVED***

func (w *Watcher) sendEvent(name string, mask uint64) bool ***REMOVED***
	if mask == 0 ***REMOVED***
		return false
	***REMOVED***
	event := newEvent(name, uint32(mask))
	select ***REMOVED***
	case ch := <-w.quit:
		w.quit <- ch
	case w.Events <- event:
	***REMOVED***
	return true
***REMOVED***

func toWindowsFlags(mask uint64) uint32 ***REMOVED***
	var m uint32
	if mask&sysFSACCESS != 0 ***REMOVED***
		m |= syscall.FILE_NOTIFY_CHANGE_LAST_ACCESS
	***REMOVED***
	if mask&sysFSMODIFY != 0 ***REMOVED***
		m |= syscall.FILE_NOTIFY_CHANGE_LAST_WRITE
	***REMOVED***
	if mask&sysFSATTRIB != 0 ***REMOVED***
		m |= syscall.FILE_NOTIFY_CHANGE_ATTRIBUTES
	***REMOVED***
	if mask&(sysFSMOVE|sysFSCREATE|sysFSDELETE) != 0 ***REMOVED***
		m |= syscall.FILE_NOTIFY_CHANGE_FILE_NAME | syscall.FILE_NOTIFY_CHANGE_DIR_NAME
	***REMOVED***
	return m
***REMOVED***

func toFSnotifyFlags(action uint32) uint64 ***REMOVED***
	switch action ***REMOVED***
	case syscall.FILE_ACTION_ADDED:
		return sysFSCREATE
	case syscall.FILE_ACTION_REMOVED:
		return sysFSDELETE
	case syscall.FILE_ACTION_MODIFIED:
		return sysFSMODIFY
	case syscall.FILE_ACTION_RENAMED_OLD_NAME:
		return sysFSMOVEDFROM
	case syscall.FILE_ACTION_RENAMED_NEW_NAME:
		return sysFSMOVEDTO
	***REMOVED***
	return 0
***REMOVED***
