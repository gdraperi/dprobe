package ioutils

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// todo: split docker/pkg/ioutils into a separate repo

// AtomicWriteFile atomically writes data to a file specified by filename.
func AtomicWriteFile(filename string, data []byte, perm os.FileMode) error ***REMOVED***
	f, err := ioutil.TempFile(filepath.Dir(filename), ".tmp-"+filepath.Base(filename))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	err = os.Chmod(f.Name(), perm)
	if err != nil ***REMOVED***
		f.Close()
		return err
	***REMOVED***
	n, err := f.Write(data)
	if err == nil && n < len(data) ***REMOVED***
		f.Close()
		return io.ErrShortWrite
	***REMOVED***
	if err != nil ***REMOVED***
		f.Close()
		return err
	***REMOVED***
	if err := f.Sync(); err != nil ***REMOVED***
		f.Close()
		return err
	***REMOVED***
	if err := f.Close(); err != nil ***REMOVED***
		return err
	***REMOVED***
	return os.Rename(f.Name(), filename)
***REMOVED***
