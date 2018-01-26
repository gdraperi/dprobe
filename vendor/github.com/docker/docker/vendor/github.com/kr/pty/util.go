package pty

import (
	"os"
	"syscall"
	"unsafe"
)

// Getsize returns the number of rows (lines) and cols (positions
// in each line) in terminal t.
func Getsize(t *os.File) (rows, cols int, err error) ***REMOVED***
	var ws winsize
	err = windowrect(&ws, t.Fd())
	return int(ws.ws_row), int(ws.ws_col), err
***REMOVED***

type winsize struct ***REMOVED***
	ws_row    uint16
	ws_col    uint16
	ws_xpixel uint16
	ws_ypixel uint16
***REMOVED***

func windowrect(ws *winsize, fd uintptr) error ***REMOVED***
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		fd,
		syscall.TIOCGWINSZ,
		uintptr(unsafe.Pointer(ws)),
	)
	if errno != 0 ***REMOVED***
		return syscall.Errno(errno)
	***REMOVED***
	return nil
***REMOVED***
