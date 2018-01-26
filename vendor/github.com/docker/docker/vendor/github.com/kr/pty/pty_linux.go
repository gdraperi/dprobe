package pty

import (
	"os"
	"strconv"
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

	err = unlockpt(p)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	t, err := os.OpenFile(sname, os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	return p, t, nil
***REMOVED***

func ptsname(f *os.File) (string, error) ***REMOVED***
	var n _C_uint
	err := ioctl(f.Fd(), syscall.TIOCGPTN, uintptr(unsafe.Pointer(&n)))
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return "/dev/pts/" + strconv.Itoa(int(n)), nil
***REMOVED***

func unlockpt(f *os.File) error ***REMOVED***
	var u _C_int
	// use TIOCSPTLCK with a zero valued arg to clear the slave pty lock
	return ioctl(f.Fd(), syscall.TIOCSPTLCK, uintptr(unsafe.Pointer(&u)))
***REMOVED***
