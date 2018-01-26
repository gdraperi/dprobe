package term

import (
	"io"
	"os"
	"os/signal"
	"syscall" // used for STD_INPUT_HANDLE, STD_OUTPUT_HANDLE and STD_ERROR_HANDLE

	"github.com/Azure/go-ansiterm/winterm"
	"github.com/docker/docker/pkg/term/windows"
)

// State holds the console mode for the terminal.
type State struct ***REMOVED***
	mode uint32
***REMOVED***

// Winsize is used for window size.
type Winsize struct ***REMOVED***
	Height uint16
	Width  uint16
***REMOVED***

// vtInputSupported is true if winterm.ENABLE_VIRTUAL_TERMINAL_INPUT is supported by the console
var vtInputSupported bool

// StdStreams returns the standard streams (stdin, stdout, stderr).
func StdStreams() (stdIn io.ReadCloser, stdOut, stdErr io.Writer) ***REMOVED***
	// Turn on VT handling on all std handles, if possible. This might
	// fail, in which case we will fall back to terminal emulation.
	var emulateStdin, emulateStdout, emulateStderr bool
	fd := os.Stdin.Fd()
	if mode, err := winterm.GetConsoleMode(fd); err == nil ***REMOVED***
		// Validate that winterm.ENABLE_VIRTUAL_TERMINAL_INPUT is supported, but do not set it.
		if err = winterm.SetConsoleMode(fd, mode|winterm.ENABLE_VIRTUAL_TERMINAL_INPUT); err != nil ***REMOVED***
			emulateStdin = true
		***REMOVED*** else ***REMOVED***
			vtInputSupported = true
		***REMOVED***
		// Unconditionally set the console mode back even on failure because SetConsoleMode
		// remembers invalid bits on input handles.
		winterm.SetConsoleMode(fd, mode)
	***REMOVED***

	fd = os.Stdout.Fd()
	if mode, err := winterm.GetConsoleMode(fd); err == nil ***REMOVED***
		// Validate winterm.DISABLE_NEWLINE_AUTO_RETURN is supported, but do not set it.
		if err = winterm.SetConsoleMode(fd, mode|winterm.ENABLE_VIRTUAL_TERMINAL_PROCESSING|winterm.DISABLE_NEWLINE_AUTO_RETURN); err != nil ***REMOVED***
			emulateStdout = true
		***REMOVED*** else ***REMOVED***
			winterm.SetConsoleMode(fd, mode|winterm.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
		***REMOVED***
	***REMOVED***

	fd = os.Stderr.Fd()
	if mode, err := winterm.GetConsoleMode(fd); err == nil ***REMOVED***
		// Validate winterm.DISABLE_NEWLINE_AUTO_RETURN is supported, but do not set it.
		if err = winterm.SetConsoleMode(fd, mode|winterm.ENABLE_VIRTUAL_TERMINAL_PROCESSING|winterm.DISABLE_NEWLINE_AUTO_RETURN); err != nil ***REMOVED***
			emulateStderr = true
		***REMOVED*** else ***REMOVED***
			winterm.SetConsoleMode(fd, mode|winterm.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
		***REMOVED***
	***REMOVED***

	if os.Getenv("ConEmuANSI") == "ON" || os.Getenv("ConsoleZVersion") != "" ***REMOVED***
		// The ConEmu and ConsoleZ terminals emulate ANSI on output streams well.
		emulateStdin = true
		emulateStdout = false
		emulateStderr = false
	***REMOVED***

	// Temporarily use STD_INPUT_HANDLE, STD_OUTPUT_HANDLE and
	// STD_ERROR_HANDLE from syscall rather than x/sys/windows as long as
	// go-ansiterm hasn't switch to x/sys/windows.
	// TODO: switch back to x/sys/windows once go-ansiterm has switched
	if emulateStdin ***REMOVED***
		stdIn = windowsconsole.NewAnsiReader(syscall.STD_INPUT_HANDLE)
	***REMOVED*** else ***REMOVED***
		stdIn = os.Stdin
	***REMOVED***

	if emulateStdout ***REMOVED***
		stdOut = windowsconsole.NewAnsiWriter(syscall.STD_OUTPUT_HANDLE)
	***REMOVED*** else ***REMOVED***
		stdOut = os.Stdout
	***REMOVED***

	if emulateStderr ***REMOVED***
		stdErr = windowsconsole.NewAnsiWriter(syscall.STD_ERROR_HANDLE)
	***REMOVED*** else ***REMOVED***
		stdErr = os.Stderr
	***REMOVED***

	return
***REMOVED***

// GetFdInfo returns the file descriptor for an os.File and indicates whether the file represents a terminal.
func GetFdInfo(in interface***REMOVED******REMOVED***) (uintptr, bool) ***REMOVED***
	return windowsconsole.GetHandleInfo(in)
***REMOVED***

// GetWinsize returns the window size based on the specified file descriptor.
func GetWinsize(fd uintptr) (*Winsize, error) ***REMOVED***
	info, err := winterm.GetConsoleScreenBufferInfo(fd)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	winsize := &Winsize***REMOVED***
		Width:  uint16(info.Window.Right - info.Window.Left + 1),
		Height: uint16(info.Window.Bottom - info.Window.Top + 1),
	***REMOVED***

	return winsize, nil
***REMOVED***

// IsTerminal returns true if the given file descriptor is a terminal.
func IsTerminal(fd uintptr) bool ***REMOVED***
	return windowsconsole.IsConsole(fd)
***REMOVED***

// RestoreTerminal restores the terminal connected to the given file descriptor
// to a previous state.
func RestoreTerminal(fd uintptr, state *State) error ***REMOVED***
	return winterm.SetConsoleMode(fd, state.mode)
***REMOVED***

// SaveState saves the state of the terminal connected to the given file descriptor.
func SaveState(fd uintptr) (*State, error) ***REMOVED***
	mode, e := winterm.GetConsoleMode(fd)
	if e != nil ***REMOVED***
		return nil, e
	***REMOVED***

	return &State***REMOVED***mode: mode***REMOVED***, nil
***REMOVED***

// DisableEcho disables echo for the terminal connected to the given file descriptor.
// -- See https://msdn.microsoft.com/en-us/library/windows/desktop/ms683462(v=vs.85).aspx
func DisableEcho(fd uintptr, state *State) error ***REMOVED***
	mode := state.mode
	mode &^= winterm.ENABLE_ECHO_INPUT
	mode |= winterm.ENABLE_PROCESSED_INPUT | winterm.ENABLE_LINE_INPUT
	err := winterm.SetConsoleMode(fd, mode)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// Register an interrupt handler to catch and restore prior state
	restoreAtInterrupt(fd, state)
	return nil
***REMOVED***

// SetRawTerminal puts the terminal connected to the given file descriptor into
// raw mode and returns the previous state. On UNIX, this puts both the input
// and output into raw mode. On Windows, it only puts the input into raw mode.
func SetRawTerminal(fd uintptr) (*State, error) ***REMOVED***
	state, err := MakeRaw(fd)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Register an interrupt handler to catch and restore prior state
	restoreAtInterrupt(fd, state)
	return state, err
***REMOVED***

// SetRawTerminalOutput puts the output of terminal connected to the given file
// descriptor into raw mode. On UNIX, this does nothing and returns nil for the
// state. On Windows, it disables LF -> CRLF translation.
func SetRawTerminalOutput(fd uintptr) (*State, error) ***REMOVED***
	state, err := SaveState(fd)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Ignore failures, since winterm.DISABLE_NEWLINE_AUTO_RETURN might not be supported on this
	// version of Windows.
	winterm.SetConsoleMode(fd, state.mode|winterm.DISABLE_NEWLINE_AUTO_RETURN)
	return state, err
***REMOVED***

// MakeRaw puts the terminal (Windows Console) connected to the given file descriptor into raw
// mode and returns the previous state of the terminal so that it can be restored.
func MakeRaw(fd uintptr) (*State, error) ***REMOVED***
	state, err := SaveState(fd)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	mode := state.mode

	// See
	// -- https://msdn.microsoft.com/en-us/library/windows/desktop/ms686033(v=vs.85).aspx
	// -- https://msdn.microsoft.com/en-us/library/windows/desktop/ms683462(v=vs.85).aspx

	// Disable these modes
	mode &^= winterm.ENABLE_ECHO_INPUT
	mode &^= winterm.ENABLE_LINE_INPUT
	mode &^= winterm.ENABLE_MOUSE_INPUT
	mode &^= winterm.ENABLE_WINDOW_INPUT
	mode &^= winterm.ENABLE_PROCESSED_INPUT

	// Enable these modes
	mode |= winterm.ENABLE_EXTENDED_FLAGS
	mode |= winterm.ENABLE_INSERT_MODE
	mode |= winterm.ENABLE_QUICK_EDIT_MODE
	if vtInputSupported ***REMOVED***
		mode |= winterm.ENABLE_VIRTUAL_TERMINAL_INPUT
	***REMOVED***

	err = winterm.SetConsoleMode(fd, mode)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return state, nil
***REMOVED***

func restoreAtInterrupt(fd uintptr, state *State) ***REMOVED***
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)

	go func() ***REMOVED***
		_ = <-sigchan
		RestoreTerminal(fd, state)
		os.Exit(0)
	***REMOVED***()
***REMOVED***
