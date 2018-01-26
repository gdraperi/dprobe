// +build !windows

// Package term provides structures and helper functions to work with
// terminal (state, sizes).
package term

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"

	"golang.org/x/sys/unix"
)

var (
	// ErrInvalidState is returned if the state of the terminal is invalid.
	ErrInvalidState = errors.New("Invalid terminal state")
)

// State represents the state of the terminal.
type State struct ***REMOVED***
	termios Termios
***REMOVED***

// Winsize represents the size of the terminal window.
type Winsize struct ***REMOVED***
	Height uint16
	Width  uint16
	x      uint16
	y      uint16
***REMOVED***

// StdStreams returns the standard streams (stdin, stdout, stderr).
func StdStreams() (stdIn io.ReadCloser, stdOut, stdErr io.Writer) ***REMOVED***
	return os.Stdin, os.Stdout, os.Stderr
***REMOVED***

// GetFdInfo returns the file descriptor for an os.File and indicates whether the file represents a terminal.
func GetFdInfo(in interface***REMOVED******REMOVED***) (uintptr, bool) ***REMOVED***
	var inFd uintptr
	var isTerminalIn bool
	if file, ok := in.(*os.File); ok ***REMOVED***
		inFd = file.Fd()
		isTerminalIn = IsTerminal(inFd)
	***REMOVED***
	return inFd, isTerminalIn
***REMOVED***

// IsTerminal returns true if the given file descriptor is a terminal.
func IsTerminal(fd uintptr) bool ***REMOVED***
	var termios Termios
	return tcget(fd, &termios) == 0
***REMOVED***

// RestoreTerminal restores the terminal connected to the given file descriptor
// to a previous state.
func RestoreTerminal(fd uintptr, state *State) error ***REMOVED***
	if state == nil ***REMOVED***
		return ErrInvalidState
	***REMOVED***
	if err := tcset(fd, &state.termios); err != 0 ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// SaveState saves the state of the terminal connected to the given file descriptor.
func SaveState(fd uintptr) (*State, error) ***REMOVED***
	var oldState State
	if err := tcget(fd, &oldState.termios); err != 0 ***REMOVED***
		return nil, err
	***REMOVED***

	return &oldState, nil
***REMOVED***

// DisableEcho applies the specified state to the terminal connected to the file
// descriptor, with echo disabled.
func DisableEcho(fd uintptr, state *State) error ***REMOVED***
	newState := state.termios
	newState.Lflag &^= unix.ECHO

	if err := tcset(fd, &newState); err != 0 ***REMOVED***
		return err
	***REMOVED***
	handleInterrupt(fd, state)
	return nil
***REMOVED***

// SetRawTerminal puts the terminal connected to the given file descriptor into
// raw mode and returns the previous state. On UNIX, this puts both the input
// and output into raw mode. On Windows, it only puts the input into raw mode.
func SetRawTerminal(fd uintptr) (*State, error) ***REMOVED***
	oldState, err := MakeRaw(fd)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	handleInterrupt(fd, oldState)
	return oldState, err
***REMOVED***

// SetRawTerminalOutput puts the output of terminal connected to the given file
// descriptor into raw mode. On UNIX, this does nothing and returns nil for the
// state. On Windows, it disables LF -> CRLF translation.
func SetRawTerminalOutput(fd uintptr) (*State, error) ***REMOVED***
	return nil, nil
***REMOVED***

func handleInterrupt(fd uintptr, state *State) ***REMOVED***
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	go func() ***REMOVED***
		for range sigchan ***REMOVED***
			// quit cleanly and the new terminal item is on a new line
			fmt.Println()
			signal.Stop(sigchan)
			close(sigchan)
			RestoreTerminal(fd, state)
			os.Exit(1)
		***REMOVED***
	***REMOVED***()
***REMOVED***
