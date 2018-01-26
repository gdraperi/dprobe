package pty

import (
	"errors"
	"os"
	"syscall"
	"unsafe"
)

func open() (pty, tty *os.File, err error) ***REMOVED***
	p, err := os.OpenFile("/dev/ptmx", os.O_RDWR, 0)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	sname, err := ptsname(p)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	err = grantpt(p)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	err = unlockpt(p)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	t, err := os.OpenFile(sname, os.O_RDWR, 0)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	return p, t, nil
***REMOVED***

func ptsname(f *os.File) (string, error) ***REMOVED***
	n := make([]byte, _IOC_PARM_LEN(syscall.TIOCPTYGNAME))

	err := ioctl(f.Fd(), syscall.TIOCPTYGNAME, uintptr(unsafe.Pointer(&n[0])))
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	for i, c := range n ***REMOVED***
		if c == 0 ***REMOVED***
			return string(n[:i]), nil
		***REMOVED***
	***REMOVED***
	return "", errors.New("TIOCPTYGNAME string not NUL-terminated")
***REMOVED***

func grantpt(f *os.File) error ***REMOVED***
	return ioctl(f.Fd(), syscall.TIOCPTYGRANT, 0)
***REMOVED***

func unlockpt(f *os.File) error ***REMOVED***
	return ioctl(f.Fd(), syscall.TIOCPTYUNLK, 0)
***REMOVED***
