// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux,!appengine netbsd openbsd

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
package terminal // import "golang.org/x/crypto/ssh/terminal"

import (
	"syscall"
	"unsafe"
)

// State contains the state of a terminal.
type State struct ***REMOVED***
	termios syscall.Termios
***REMOVED***

// IsTerminal returns true if the given file descriptor is a terminal.
func IsTerminal(fd int) bool ***REMOVED***
	var termios syscall.Termios
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), ioctlReadTermios, uintptr(unsafe.Pointer(&termios)), 0, 0, 0)
	return err == 0
***REMOVED***

// MakeRaw put the terminal connected to the given file descriptor into raw
// mode and returns the previous state of the terminal so that it can be
// restored.
func MakeRaw(fd int) (*State, error) ***REMOVED***
	var oldState State
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), ioctlReadTermios, uintptr(unsafe.Pointer(&oldState.termios)), 0, 0, 0); err != 0 ***REMOVED***
		return nil, err
	***REMOVED***

	newState := oldState.termios
	// This attempts to replicate the behaviour documented for cfmakeraw in
	// the termios(3) manpage.
	newState.Iflag &^= syscall.IGNBRK | syscall.BRKINT | syscall.PARMRK | syscall.ISTRIP | syscall.INLCR | syscall.IGNCR | syscall.ICRNL | syscall.IXON
	newState.Oflag &^= syscall.OPOST
	newState.Lflag &^= syscall.ECHO | syscall.ECHONL | syscall.ICANON | syscall.ISIG | syscall.IEXTEN
	newState.Cflag &^= syscall.CSIZE | syscall.PARENB
	newState.Cflag |= syscall.CS8
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), ioctlWriteTermios, uintptr(unsafe.Pointer(&newState)), 0, 0, 0); err != 0 ***REMOVED***
		return nil, err
	***REMOVED***

	return &oldState, nil
***REMOVED***

// GetState returns the current state of a terminal which may be useful to
// restore the terminal after a signal.
func GetState(fd int) (*State, error) ***REMOVED***
	var oldState State
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), ioctlReadTermios, uintptr(unsafe.Pointer(&oldState.termios)), 0, 0, 0); err != 0 ***REMOVED***
		return nil, err
	***REMOVED***

	return &oldState, nil
***REMOVED***

// Restore restores the terminal connected to the given file descriptor to a
// previous state.
func Restore(fd int, state *State) error ***REMOVED***
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), ioctlWriteTermios, uintptr(unsafe.Pointer(&state.termios)), 0, 0, 0); err != 0 ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// GetSize returns the dimensions of the given terminal.
func GetSize(fd int) (width, height int, err error) ***REMOVED***
	var dimensions [4]uint16

	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&dimensions)), 0, 0, 0); err != 0 ***REMOVED***
		return -1, -1, err
	***REMOVED***
	return int(dimensions[1]), int(dimensions[0]), nil
***REMOVED***

// passwordReader is an io.Reader that reads from a specific file descriptor.
type passwordReader int

func (r passwordReader) Read(buf []byte) (int, error) ***REMOVED***
	return syscall.Read(int(r), buf)
***REMOVED***

// ReadPassword reads a line of input from a terminal without local echo.  This
// is commonly used for inputting passwords and other sensitive data. The slice
// returned does not include the \n.
func ReadPassword(fd int) ([]byte, error) ***REMOVED***
	var oldState syscall.Termios
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), ioctlReadTermios, uintptr(unsafe.Pointer(&oldState)), 0, 0, 0); err != 0 ***REMOVED***
		return nil, err
	***REMOVED***

	newState := oldState
	newState.Lflag &^= syscall.ECHO
	newState.Lflag |= syscall.ICANON | syscall.ISIG
	newState.Iflag |= syscall.ICRNL
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), ioctlWriteTermios, uintptr(unsafe.Pointer(&newState)), 0, 0, 0); err != 0 ***REMOVED***
		return nil, err
	***REMOVED***

	defer func() ***REMOVED***
		syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), ioctlWriteTermios, uintptr(unsafe.Pointer(&oldState)), 0, 0, 0)
	***REMOVED***()

	return readPasswordLine(passwordReader(fd))
***REMOVED***
