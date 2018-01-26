// +build !linux

package fifo

import (
	"syscall"

	"github.com/pkg/errors"
)

type handle struct ***REMOVED***
	fn  string
	dev uint64
	ino uint64
***REMOVED***

func getHandle(fn string) (*handle, error) ***REMOVED***
	var stat syscall.Stat_t
	if err := syscall.Stat(fn, &stat); err != nil ***REMOVED***
		return nil, errors.Wrapf(err, "failed to stat %v", fn)
	***REMOVED***

	h := &handle***REMOVED***
		fn:  fn,
		dev: uint64(stat.Dev),
		ino: uint64(stat.Ino),
	***REMOVED***

	return h, nil
***REMOVED***

func (h *handle) Path() (string, error) ***REMOVED***
	var stat syscall.Stat_t
	if err := syscall.Stat(h.fn, &stat); err != nil ***REMOVED***
		return "", errors.Wrapf(err, "path %v could not be statted", h.fn)
	***REMOVED***
	if uint64(stat.Dev) != h.dev || uint64(stat.Ino) != h.ino ***REMOVED***
		return "", errors.Errorf("failed to verify handle %v/%v %v/%v for %v", stat.Dev, h.dev, stat.Ino, h.ino, h.fn)
	***REMOVED***
	return h.fn, nil
***REMOVED***

func (h *handle) Name() string ***REMOVED***
	return h.fn
***REMOVED***

func (h *handle) Close() error ***REMOVED***
	return nil
***REMOVED***
