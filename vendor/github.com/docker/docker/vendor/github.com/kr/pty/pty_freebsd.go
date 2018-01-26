package pty

import (
	"errors"
	"os"
	"syscall"
	"unsafe"
)

func posix_openpt(oflag int) (fd int, err error) ***REMOVED***
	r0, _, e1 := syscall.Syscall(syscall.SYS_POSIX_OPENPT, uintptr(oflag), 0, 0)
	fd = int(r0)
	if e1 != 0 ***REMOVED***
		err = e1
	***REMOVED***
	return
***REMOVED***

func open() (pty, tty *os.File, err error) ***REMOVED***
	fd, err := posix_openpt(syscall.O_RDWR | syscall.O_CLOEXEC)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	p := os.NewFile(uintptr(fd), "/dev/pts")
	sname, err := ptsname(p)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***

	t, err := os.OpenFile("/dev/"+sname, os.O_RDWR, 0)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	return p, t, nil
***REMOVED***

func isptmaster(fd uintptr) (bool, error) ***REMOVED***
	err := ioctl(fd, syscall.TIOCPTMASTER, 0)
	return err == nil, err
***REMOVED***

var (
	emptyFiodgnameArg fiodgnameArg
	ioctl_FIODGNAME   = _IOW('f', 120, unsafe.Sizeof(emptyFiodgnameArg))
)

func ptsname(f *os.File) (string, error) ***REMOVED***
	master, err := isptmaster(f.Fd())
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	if !master ***REMOVED***
		return "", syscall.EINVAL
	***REMOVED***

	const n = _C_SPECNAMELEN + 1
	var (
		buf = make([]byte, n)
		arg = fiodgnameArg***REMOVED***Len: n, Buf: (*byte)(unsafe.Pointer(&buf[0]))***REMOVED***
	)
	err = ioctl(f.Fd(), ioctl_FIODGNAME, uintptr(unsafe.Pointer(&arg)))
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	for i, c := range buf ***REMOVED***
		if c == 0 ***REMOVED***
			return string(buf[:i]), nil
		***REMOVED***
	***REMOVED***
	return "", errors.New("FIODGNAME string not NUL-terminated")
***REMOVED***
