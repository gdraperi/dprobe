// +build !windows

package term

import (
	"golang.org/x/sys/unix"
)

// GetWinsize returns the window size based on the specified file descriptor.
func GetWinsize(fd uintptr) (*Winsize, error) ***REMOVED***
	uws, err := unix.IoctlGetWinsize(int(fd), unix.TIOCGWINSZ)
	ws := &Winsize***REMOVED***Height: uws.Row, Width: uws.Col, x: uws.Xpixel, y: uws.Ypixel***REMOVED***
	return ws, err
***REMOVED***

// SetWinsize tries to set the specified window size for the specified file descriptor.
func SetWinsize(fd uintptr, ws *Winsize) error ***REMOVED***
	uws := &unix.Winsize***REMOVED***Row: ws.Height, Col: ws.Width, Xpixel: ws.x, Ypixel: ws.y***REMOVED***
	return unix.IoctlSetWinsize(int(fd), unix.TIOCSWINSZ, uws)
***REMOVED***
