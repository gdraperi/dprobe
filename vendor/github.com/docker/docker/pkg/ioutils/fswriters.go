package ioutils

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// NewAtomicFileWriter returns WriteCloser so that writing to it writes to a
// temporary file and closing it atomically changes the temporary file to
// destination path. Writing and closing concurrently is not allowed.
func NewAtomicFileWriter(filename string, perm os.FileMode) (io.WriteCloser, error) ***REMOVED***
	f, err := ioutil.TempFile(filepath.Dir(filename), ".tmp-"+filepath.Base(filename))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	abspath, err := filepath.Abs(filename)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &atomicFileWriter***REMOVED***
		f:    f,
		fn:   abspath,
		perm: perm,
	***REMOVED***, nil
***REMOVED***

// AtomicWriteFile atomically writes data to a file named by filename.
func AtomicWriteFile(filename string, data []byte, perm os.FileMode) error ***REMOVED***
	f, err := NewAtomicFileWriter(filename, perm)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	n, err := f.Write(data)
	if err == nil && n < len(data) ***REMOVED***
		err = io.ErrShortWrite
		f.(*atomicFileWriter).writeErr = err
	***REMOVED***
	if err1 := f.Close(); err == nil ***REMOVED***
		err = err1
	***REMOVED***
	return err
***REMOVED***

type atomicFileWriter struct ***REMOVED***
	f        *os.File
	fn       string
	writeErr error
	perm     os.FileMode
***REMOVED***

func (w *atomicFileWriter) Write(dt []byte) (int, error) ***REMOVED***
	n, err := w.f.Write(dt)
	if err != nil ***REMOVED***
		w.writeErr = err
	***REMOVED***
	return n, err
***REMOVED***

func (w *atomicFileWriter) Close() (retErr error) ***REMOVED***
	defer func() ***REMOVED***
		if retErr != nil || w.writeErr != nil ***REMOVED***
			os.Remove(w.f.Name())
		***REMOVED***
	***REMOVED***()
	if err := w.f.Sync(); err != nil ***REMOVED***
		w.f.Close()
		return err
	***REMOVED***
	if err := w.f.Close(); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := os.Chmod(w.f.Name(), w.perm); err != nil ***REMOVED***
		return err
	***REMOVED***
	if w.writeErr == nil ***REMOVED***
		return os.Rename(w.f.Name(), w.fn)
	***REMOVED***
	return nil
***REMOVED***

// AtomicWriteSet is used to atomically write a set
// of files and ensure they are visible at the same time.
// Must be committed to a new directory.
type AtomicWriteSet struct ***REMOVED***
	root string
***REMOVED***

// NewAtomicWriteSet creates a new atomic write set to
// atomically create a set of files. The given directory
// is used as the base directory for storing files before
// commit. If no temporary directory is given the system
// default is used.
func NewAtomicWriteSet(tmpDir string) (*AtomicWriteSet, error) ***REMOVED***
	td, err := ioutil.TempDir(tmpDir, "write-set-")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return &AtomicWriteSet***REMOVED***
		root: td,
	***REMOVED***, nil
***REMOVED***

// WriteFile writes a file to the set, guaranteeing the file
// has been synced.
func (ws *AtomicWriteSet) WriteFile(filename string, data []byte, perm os.FileMode) error ***REMOVED***
	f, err := ws.FileWriter(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	n, err := f.Write(data)
	if err == nil && n < len(data) ***REMOVED***
		err = io.ErrShortWrite
	***REMOVED***
	if err1 := f.Close(); err == nil ***REMOVED***
		err = err1
	***REMOVED***
	return err
***REMOVED***

type syncFileCloser struct ***REMOVED***
	*os.File
***REMOVED***

func (w syncFileCloser) Close() error ***REMOVED***
	err := w.File.Sync()
	if err1 := w.File.Close(); err == nil ***REMOVED***
		err = err1
	***REMOVED***
	return err
***REMOVED***

// FileWriter opens a file writer inside the set. The file
// should be synced and closed before calling commit.
func (ws *AtomicWriteSet) FileWriter(name string, flag int, perm os.FileMode) (io.WriteCloser, error) ***REMOVED***
	f, err := os.OpenFile(filepath.Join(ws.root, name), flag, perm)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return syncFileCloser***REMOVED***f***REMOVED***, nil
***REMOVED***

// Cancel cancels the set and removes all temporary data
// created in the set.
func (ws *AtomicWriteSet) Cancel() error ***REMOVED***
	return os.RemoveAll(ws.root)
***REMOVED***

// Commit moves all created files to the target directory. The
// target directory must not exist and the parent of the target
// directory must exist.
func (ws *AtomicWriteSet) Commit(target string) error ***REMOVED***
	return os.Rename(ws.root, target)
***REMOVED***

// String returns the location the set is writing to.
func (ws *AtomicWriteSet) String() string ***REMOVED***
	return ws.root
***REMOVED***
