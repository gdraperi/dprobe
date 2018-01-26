// +build darwin freebsd linux solaris

package console

import (
	"os"

	"golang.org/x/sys/unix"
)

// NewPty creates a new pty pair
// The master is returned as the first console and a string
// with the path to the pty slave is returned as the second
func NewPty() (Console, string, error) ***REMOVED***
	f, err := os.OpenFile("/dev/ptmx", unix.O_RDWR|unix.O_NOCTTY|unix.O_CLOEXEC, 0)
	if err != nil ***REMOVED***
		return nil, "", err
	***REMOVED***
	slave, err := ptsname(f)
	if err != nil ***REMOVED***
		return nil, "", err
	***REMOVED***
	if err := unlockpt(f); err != nil ***REMOVED***
		return nil, "", err
	***REMOVED***
	m, err := newMaster(f)
	if err != nil ***REMOVED***
		return nil, "", err
	***REMOVED***
	return m, slave, nil
***REMOVED***

type master struct ***REMOVED***
	f        *os.File
	original *unix.Termios
***REMOVED***

func (m *master) Read(b []byte) (int, error) ***REMOVED***
	return m.f.Read(b)
***REMOVED***

func (m *master) Write(b []byte) (int, error) ***REMOVED***
	return m.f.Write(b)
***REMOVED***

func (m *master) Close() error ***REMOVED***
	return m.f.Close()
***REMOVED***

func (m *master) Resize(ws WinSize) error ***REMOVED***
	return tcswinsz(m.f.Fd(), ws)
***REMOVED***

func (m *master) ResizeFrom(c Console) error ***REMOVED***
	ws, err := c.Size()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return m.Resize(ws)
***REMOVED***

func (m *master) Reset() error ***REMOVED***
	if m.original == nil ***REMOVED***
		return nil
	***REMOVED***
	return tcset(m.f.Fd(), m.original)
***REMOVED***

func (m *master) getCurrent() (unix.Termios, error) ***REMOVED***
	var termios unix.Termios
	if err := tcget(m.f.Fd(), &termios); err != nil ***REMOVED***
		return unix.Termios***REMOVED******REMOVED***, err
	***REMOVED***
	return termios, nil
***REMOVED***

func (m *master) SetRaw() error ***REMOVED***
	rawState, err := m.getCurrent()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	rawState = cfmakeraw(rawState)
	rawState.Oflag = rawState.Oflag | unix.OPOST
	return tcset(m.f.Fd(), &rawState)
***REMOVED***

func (m *master) DisableEcho() error ***REMOVED***
	rawState, err := m.getCurrent()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	rawState.Lflag = rawState.Lflag &^ unix.ECHO
	return tcset(m.f.Fd(), &rawState)
***REMOVED***

func (m *master) Size() (WinSize, error) ***REMOVED***
	return tcgwinsz(m.f.Fd())
***REMOVED***

func (m *master) Fd() uintptr ***REMOVED***
	return m.f.Fd()
***REMOVED***

func (m *master) Name() string ***REMOVED***
	return m.f.Name()
***REMOVED***

// checkConsole checks if the provided file is a console
func checkConsole(f *os.File) error ***REMOVED***
	var termios unix.Termios
	if tcget(f.Fd(), &termios) != nil ***REMOVED***
		return ErrNotAConsole
	***REMOVED***
	return nil
***REMOVED***

func newMaster(f *os.File) (Console, error) ***REMOVED***
	m := &master***REMOVED***
		f: f,
	***REMOVED***
	t, err := m.getCurrent()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	m.original = &t
	return m, nil
***REMOVED***

// ClearONLCR sets the necessary tty_ioctl(4)s to ensure that a pty pair
// created by us acts normally. In particular, a not-very-well-known default of
// Linux unix98 ptys is that they have +onlcr by default. While this isn't a
// problem for terminal emulators, because we relay data from the terminal we
// also relay that funky line discipline.
func ClearONLCR(fd uintptr) error ***REMOVED***
	return setONLCR(fd, false)
***REMOVED***

// SetONLCR sets the necessary tty_ioctl(4)s to ensure that a pty pair
// created by us acts as intended for a terminal emulator.
func SetONLCR(fd uintptr) error ***REMOVED***
	return setONLCR(fd, true)
***REMOVED***
