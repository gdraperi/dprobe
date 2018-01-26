// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

// Package terminal provides support functions for dealing with terminals, as
// commonly found on UNIX systems.
//
// Putting a terminal into raw mode is the most common requirement:
//
// 	oldState, err := terminal.MakeRaw(0)
// 	if err != nil ***REMOVED***
// 	        panic(err)
// 	***REMOVED***
// 	defer terminal.Restore(0, oldState)
package terminal

import (
	"syscall"
	"unsafe"
)

const (
	enableLineInput       = 2
	enableEchoInput       = 4
	enableProcessedInput  = 1
	enableWindowInput     = 8
	enableMouseInput      = 16
	enableInsertMode      = 32
	enableQuickEditMode   = 64
	enableExtendedFlags   = 128
	enableAutoPosition    = 256
	enableProcessedOutput = 1
	enableWrapAtEolOutput = 2
)

var kernel32 = syscall.NewLazyDLL("kernel32.dll")

var (
	procGetConsoleMode             = kernel32.NewProc("GetConsoleMode")
	procSetConsoleMode             = kernel32.NewProc("SetConsoleMode")
	procGetConsoleScreenBufferInfo = kernel32.NewProc("GetConsoleScreenBufferInfo")
)

type (
	short int16
	word  uint16

	coord struct ***REMOVED***
		x short
		y short
	***REMOVED***
	smallRect struct ***REMOVED***
		left   short
		top    short
		right  short
		bottom short
	***REMOVED***
	consoleScreenBufferInfo struct ***REMOVED***
		size              coord
		cursorPosition    coord
		attributes        word
		window            smallRect
		maximumWindowSize coord
	***REMOVED***
)

type State struct ***REMOVED***
	mode uint32
***REMOVED***

// IsTerminal returns true if the given file descriptor is a terminal.
func IsTerminal(fd int) bool ***REMOVED***
	var st uint32
	r, _, e := syscall.Syscall(procGetConsoleMode.Addr(), 2, uintptr(fd), uintptr(unsafe.Pointer(&st)), 0)
	return r != 0 && e == 0
***REMOVED***

// MakeRaw put the terminal connected to the given file descriptor into raw
// mode and returns the previous state of the terminal so that it can be
// restored.
func MakeRaw(fd int) (*State, error) ***REMOVED***
	var st uint32
	_, _, e := syscall.Syscall(procGetConsoleMode.Addr(), 2, uintptr(fd), uintptr(unsafe.Pointer(&st)), 0)
	if e != 0 ***REMOVED***
		return nil, error(e)
	***REMOVED***
	raw := st &^ (enableEchoInput | enableProcessedInput | enableLineInput | enableProcessedOutput)
	_, _, e = syscall.Syscall(procSetConsoleMode.Addr(), 2, uintptr(fd), uintptr(raw), 0)
	if e != 0 ***REMOVED***
		return nil, error(e)
	***REMOVED***
	return &State***REMOVED***st***REMOVED***, nil
***REMOVED***

// GetState returns the current state of a terminal which may be useful to
// restore the terminal after a signal.
func GetState(fd int) (*State, error) ***REMOVED***
	var st uint32
	_, _, e := syscall.Syscall(procGetConsoleMode.Addr(), 2, uintptr(fd), uintptr(unsafe.Pointer(&st)), 0)
	if e != 0 ***REMOVED***
		return nil, error(e)
	***REMOVED***
	return &State***REMOVED***st***REMOVED***, nil
***REMOVED***

// Restore restores the terminal connected to the given file descriptor to a
// previous state.
func Restore(fd int, state *State) error ***REMOVED***
	_, _, err := syscall.Syscall(procSetConsoleMode.Addr(), 2, uintptr(fd), uintptr(state.mode), 0)
	return err
***REMOVED***

// GetSize returns the dimensions of the given terminal.
func GetSize(fd int) (width, height int, err error) ***REMOVED***
	var info consoleScreenBufferInfo
	_, _, e := syscall.Syscall(procGetConsoleScreenBufferInfo.Addr(), 2, uintptr(fd), uintptr(unsafe.Pointer(&info)), 0)
	if e != 0 ***REMOVED***
		return 0, 0, error(e)
	***REMOVED***
	return int(info.size.x), int(info.size.y), nil
***REMOVED***

// passwordReader is an io.Reader that reads from a specific Windows HANDLE.
type passwordReader int

func (r passwordReader) Read(buf []byte) (int, error) ***REMOVED***
	return syscall.Read(syscall.Handle(r), buf)
***REMOVED***

// ReadPassword reads a line of input from a terminal without local echo.  This
// is commonly used for inputting passwords and other sensitive data. The slice
// returned does not include the \n.
func ReadPassword(fd int) ([]byte, error) ***REMOVED***
	var st uint32
	_, _, e := syscall.Syscall(procGetConsoleMode.Addr(), 2, uintptr(fd), uintptr(unsafe.Pointer(&st)), 0)
	if e != 0 ***REMOVED***
		return nil, error(e)
	***REMOVED***
	old := st

	st &^= (enableEchoInput)
	st |= (enableProcessedInput | enableLineInput | enableProcessedOutput)
	_, _, e = syscall.Syscall(procSetConsoleMode.Addr(), 2, uintptr(fd), uintptr(st), 0)
	if e != 0 ***REMOVED***
		return nil, error(e)
	***REMOVED***

	defer func() ***REMOVED***
		syscall.Syscall(procSetConsoleMode.Addr(), 2, uintptr(fd), uintptr(old), 0)
	***REMOVED***()

	return readPasswordLine(passwordReader(fd))
***REMOVED***
