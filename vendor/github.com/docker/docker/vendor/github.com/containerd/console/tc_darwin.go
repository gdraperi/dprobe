package console

import (
	"fmt"
	"os"
	"unsafe"

	"golang.org/x/sys/unix"
)

const (
	cmdTcGet = unix.TIOCGETA
	cmdTcSet = unix.TIOCSETA
)

func ioctl(fd, flag, data uintptr) error ***REMOVED***
	if _, _, err := unix.Syscall(unix.SYS_IOCTL, fd, flag, data); err != 0 ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// unlockpt unlocks the slave pseudoterminal device corresponding to the master pseudoterminal referred to by f.
// unlockpt should be called before opening the slave side of a pty.
func unlockpt(f *os.File) error ***REMOVED***
	var u int32
	return ioctl(f.Fd(), unix.TIOCPTYUNLK, uintptr(unsafe.Pointer(&u)))
***REMOVED***

// ptsname retrieves the name of the first available pts for the given master.
func ptsname(f *os.File) (string, error) ***REMOVED***
	n, err := unix.IoctlGetInt(int(f.Fd()), unix.TIOCPTYGNAME)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return fmt.Sprintf("/dev/pts/%d", n), nil
***REMOVED***
