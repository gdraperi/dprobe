package console

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

const (
	cmdTcGet = unix.TIOCGETA
	cmdTcSet = unix.TIOCSETA
)

// unlockpt unlocks the slave pseudoterminal device corresponding to the master pseudoterminal referred to by f.
// unlockpt should be called before opening the slave side of a pty.
// This does not exist on FreeBSD, it does not allocate controlling terminals on open
func unlockpt(f *os.File) error ***REMOVED***
	return nil
***REMOVED***

// ptsname retrieves the name of the first available pts for the given master.
func ptsname(f *os.File) (string, error) ***REMOVED***
	n, err := unix.IoctlGetInt(int(f.Fd()), unix.TIOCGPTN)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return fmt.Sprintf("/dev/pts/%d", n), nil
***REMOVED***
