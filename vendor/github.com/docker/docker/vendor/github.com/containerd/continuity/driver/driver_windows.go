package driver

import (
	"os"

	"github.com/pkg/errors"
)

func (d *driver) Mknod(path string, mode os.FileMode, major, minor int) error ***REMOVED***
	return errors.Wrap(ErrNotSupported, "cannot create device node on Windows")
***REMOVED***

func (d *driver) Mkfifo(path string, mode os.FileMode) error ***REMOVED***
	return errors.Wrap(ErrNotSupported, "cannot create fifo on Windows")
***REMOVED***

// Lchmod changes the mode of an file not following symlinks.
func (d *driver) Lchmod(path string, mode os.FileMode) (err error) ***REMOVED***
	// TODO: Use Window's equivalent
	return os.Chmod(path, mode)
***REMOVED***
