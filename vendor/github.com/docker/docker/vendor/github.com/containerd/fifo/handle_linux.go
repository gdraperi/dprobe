// +build linux

package fifo

import (
	"fmt"
	"os"
	"sync"
	"syscall"

	"github.com/pkg/errors"
)

const O_PATH = 010000000

type handle struct ***REMOVED***
	f         *os.File
	fd        uintptr
	dev       uint64
	ino       uint64
	closeOnce sync.Once
	name      string
***REMOVED***

func getHandle(fn string) (*handle, error) ***REMOVED***
	f, err := os.OpenFile(fn, O_PATH, 0)
	if err != nil ***REMOVED***
		return nil, errors.Wrapf(err, "failed to open %v with O_PATH", fn)
	***REMOVED***

	var (
		stat syscall.Stat_t
		fd   = f.Fd()
	)
	if err := syscall.Fstat(int(fd), &stat); err != nil ***REMOVED***
		f.Close()
		return nil, errors.Wrapf(err, "failed to stat handle %v", fd)
	***REMOVED***

	h := &handle***REMOVED***
		f:    f,
		name: fn,
		dev:  uint64(stat.Dev),
		ino:  stat.Ino,
		fd:   fd,
	***REMOVED***

	// check /proc just in case
	if _, err := os.Stat(h.procPath()); err != nil ***REMOVED***
		f.Close()
		return nil, errors.Wrapf(err, "couldn't stat %v", h.procPath())
	***REMOVED***

	return h, nil
***REMOVED***

func (h *handle) procPath() string ***REMOVED***
	return fmt.Sprintf("/proc/self/fd/%d", h.fd)
***REMOVED***

func (h *handle) Name() string ***REMOVED***
	return h.name
***REMOVED***

func (h *handle) Path() (string, error) ***REMOVED***
	var stat syscall.Stat_t
	if err := syscall.Stat(h.procPath(), &stat); err != nil ***REMOVED***
		return "", errors.Wrapf(err, "path %v could not be statted", h.procPath())
	***REMOVED***
	if uint64(stat.Dev) != h.dev || stat.Ino != h.ino ***REMOVED***
		return "", errors.Errorf("failed to verify handle %v/%v %v/%v", stat.Dev, h.dev, stat.Ino, h.ino)
	***REMOVED***
	return h.procPath(), nil
***REMOVED***

func (h *handle) Close() error ***REMOVED***
	h.closeOnce.Do(func() ***REMOVED***
		h.f.Close()
	***REMOVED***)
	return nil
***REMOVED***
