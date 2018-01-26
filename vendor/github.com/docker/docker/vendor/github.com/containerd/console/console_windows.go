package console

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"golang.org/x/sys/windows"
)

var (
	vtInputSupported  bool
	ErrNotImplemented = errors.New("not implemented")
)

func (m *master) initStdios() ***REMOVED***
	m.in = windows.Handle(os.Stdin.Fd())
	if err := windows.GetConsoleMode(m.in, &m.inMode); err == nil ***REMOVED***
		// Validate that windows.ENABLE_VIRTUAL_TERMINAL_INPUT is supported, but do not set it.
		if err = windows.SetConsoleMode(m.in, m.inMode|windows.ENABLE_VIRTUAL_TERMINAL_INPUT); err == nil ***REMOVED***
			vtInputSupported = true
		***REMOVED***
		// Unconditionally set the console mode back even on failure because SetConsoleMode
		// remembers invalid bits on input handles.
		windows.SetConsoleMode(m.in, m.inMode)
	***REMOVED*** else ***REMOVED***
		fmt.Printf("failed to get console mode for stdin: %v\n", err)
	***REMOVED***

	m.out = windows.Handle(os.Stdout.Fd())
	if err := windows.GetConsoleMode(m.out, &m.outMode); err == nil ***REMOVED***
		if err := windows.SetConsoleMode(m.out, m.outMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING); err == nil ***REMOVED***
			m.outMode |= windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
		***REMOVED*** else ***REMOVED***
			windows.SetConsoleMode(m.out, m.outMode)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		fmt.Printf("failed to get console mode for stdout: %v\n", err)
	***REMOVED***

	m.err = windows.Handle(os.Stderr.Fd())
	if err := windows.GetConsoleMode(m.err, &m.errMode); err == nil ***REMOVED***
		if err := windows.SetConsoleMode(m.err, m.errMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING); err == nil ***REMOVED***
			m.errMode |= windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
		***REMOVED*** else ***REMOVED***
			windows.SetConsoleMode(m.err, m.errMode)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		fmt.Printf("failed to get console mode for stderr: %v\n", err)
	***REMOVED***
***REMOVED***

type master struct ***REMOVED***
	in     windows.Handle
	inMode uint32

	out     windows.Handle
	outMode uint32

	err     windows.Handle
	errMode uint32
***REMOVED***

func (m *master) SetRaw() error ***REMOVED***
	if err := makeInputRaw(m.in, m.inMode); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Set StdOut and StdErr to raw mode, we ignore failures since
	// windows.DISABLE_NEWLINE_AUTO_RETURN might not be supported on this version of
	// Windows.

	windows.SetConsoleMode(m.out, m.outMode|windows.DISABLE_NEWLINE_AUTO_RETURN)

	windows.SetConsoleMode(m.err, m.errMode|windows.DISABLE_NEWLINE_AUTO_RETURN)

	return nil
***REMOVED***

func (m *master) Reset() error ***REMOVED***
	for _, s := range []struct ***REMOVED***
		fd   windows.Handle
		mode uint32
	***REMOVED******REMOVED***
		***REMOVED***m.in, m.inMode***REMOVED***,
		***REMOVED***m.out, m.outMode***REMOVED***,
		***REMOVED***m.err, m.errMode***REMOVED***,
	***REMOVED*** ***REMOVED***
		if err := windows.SetConsoleMode(s.fd, s.mode); err != nil ***REMOVED***
			return errors.Wrap(err, "unable to restore console mode")
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func (m *master) Size() (WinSize, error) ***REMOVED***
	var info windows.ConsoleScreenBufferInfo
	err := windows.GetConsoleScreenBufferInfo(m.out, &info)
	if err != nil ***REMOVED***
		return WinSize***REMOVED******REMOVED***, errors.Wrap(err, "unable to get console info")
	***REMOVED***

	winsize := WinSize***REMOVED***
		Width:  uint16(info.Window.Right - info.Window.Left + 1),
		Height: uint16(info.Window.Bottom - info.Window.Top + 1),
	***REMOVED***

	return winsize, nil
***REMOVED***

func (m *master) Resize(ws WinSize) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func (m *master) ResizeFrom(c Console) error ***REMOVED***
	return ErrNotImplemented
***REMOVED***

func (m *master) DisableEcho() error ***REMOVED***
	mode := m.inMode &^ windows.ENABLE_ECHO_INPUT
	mode |= windows.ENABLE_PROCESSED_INPUT
	mode |= windows.ENABLE_LINE_INPUT

	if err := windows.SetConsoleMode(m.in, mode); err != nil ***REMOVED***
		return errors.Wrap(err, "unable to set console to disable echo")
	***REMOVED***

	return nil
***REMOVED***

func (m *master) Close() error ***REMOVED***
	return nil
***REMOVED***

func (m *master) Read(b []byte) (int, error) ***REMOVED***
	panic("not implemented on windows")
***REMOVED***

func (m *master) Write(b []byte) (int, error) ***REMOVED***
	panic("not implemented on windows")
***REMOVED***

func (m *master) Fd() uintptr ***REMOVED***
	return uintptr(m.in)
***REMOVED***

// on windows, console can only be made from os.Std***REMOVED***in,out,err***REMOVED***, hence there
// isnt a single name here we can use. Return a dummy "console" value in this
// case should be sufficient.
func (m *master) Name() string ***REMOVED***
	return "console"
***REMOVED***

// makeInputRaw puts the terminal (Windows Console) connected to the given
// file descriptor into raw mode
func makeInputRaw(fd windows.Handle, mode uint32) error ***REMOVED***
	// See
	// -- https://msdn.microsoft.com/en-us/library/windows/desktop/ms686033(v=vs.85).aspx
	// -- https://msdn.microsoft.com/en-us/library/windows/desktop/ms683462(v=vs.85).aspx

	// Disable these modes
	mode &^= windows.ENABLE_ECHO_INPUT
	mode &^= windows.ENABLE_LINE_INPUT
	mode &^= windows.ENABLE_MOUSE_INPUT
	mode &^= windows.ENABLE_WINDOW_INPUT
	mode &^= windows.ENABLE_PROCESSED_INPUT

	// Enable these modes
	mode |= windows.ENABLE_EXTENDED_FLAGS
	mode |= windows.ENABLE_INSERT_MODE
	mode |= windows.ENABLE_QUICK_EDIT_MODE

	if vtInputSupported ***REMOVED***
		mode |= windows.ENABLE_VIRTUAL_TERMINAL_INPUT
	***REMOVED***

	if err := windows.SetConsoleMode(fd, mode); err != nil ***REMOVED***
		return errors.Wrap(err, "unable to set console to raw mode")
	***REMOVED***

	return nil
***REMOVED***

func checkConsole(f *os.File) error ***REMOVED***
	var mode uint32
	if err := windows.GetConsoleMode(windows.Handle(f.Fd()), &mode); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

func newMaster(f *os.File) (Console, error) ***REMOVED***
	if f != os.Stdin && f != os.Stdout && f != os.Stderr ***REMOVED***
		return nil, errors.New("creating a console from a file is not supported on windows")
	***REMOVED***
	m := &master***REMOVED******REMOVED***
	m.initStdios()
	return m, nil
***REMOVED***
