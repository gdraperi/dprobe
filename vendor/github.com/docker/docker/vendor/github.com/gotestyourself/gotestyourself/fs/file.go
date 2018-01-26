/*Package fs provides tools for creating and working with temporary files and
directories.
*/
package fs

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gotestyourself/gotestyourself/assert"
)

// Path objects return their filesystem path. Both File and Dir implement Path.
type Path interface ***REMOVED***
	Path() string
	Remove()
***REMOVED***

var (
	_ Path = &Dir***REMOVED******REMOVED***
	_ Path = &File***REMOVED******REMOVED***
)

// File is a temporary file on the filesystem
type File struct ***REMOVED***
	path string
***REMOVED***

type helperT interface ***REMOVED***
	Helper()
***REMOVED***

// NewFile creates a new file in a temporary directory using prefix as part of
// the filename. The PathOps are applied to the before returning the File.
func NewFile(t assert.TestingT, prefix string, ops ...PathOp) *File ***REMOVED***
	if ht, ok := t.(helperT); ok ***REMOVED***
		ht.Helper()
	***REMOVED***
	tempfile, err := ioutil.TempFile("", prefix+"-")
	assert.NilError(t, err)
	file := &File***REMOVED***path: tempfile.Name()***REMOVED***
	assert.NilError(t, tempfile.Close())

	for _, op := range ops ***REMOVED***
		assert.NilError(t, op(file))
	***REMOVED***
	return file
***REMOVED***

// Path returns the full path to the file
func (f *File) Path() string ***REMOVED***
	return f.path
***REMOVED***

// Remove the file
func (f *File) Remove() ***REMOVED***
	// nolint: errcheck
	os.Remove(f.path)
***REMOVED***

// Dir is a temporary directory
type Dir struct ***REMOVED***
	path string
***REMOVED***

// NewDir returns a new temporary directory using prefix as part of the directory
// name. The PathOps are applied before returning the Dir.
func NewDir(t assert.TestingT, prefix string, ops ...PathOp) *Dir ***REMOVED***
	if ht, ok := t.(helperT); ok ***REMOVED***
		ht.Helper()
	***REMOVED***
	path, err := ioutil.TempDir("", prefix+"-")
	assert.NilError(t, err)
	dir := &Dir***REMOVED***path: path***REMOVED***

	for _, op := range ops ***REMOVED***
		assert.NilError(t, op(dir))
	***REMOVED***
	return dir
***REMOVED***

// Path returns the full path to the directory
func (d *Dir) Path() string ***REMOVED***
	return d.path
***REMOVED***

// Remove the directory
func (d *Dir) Remove() ***REMOVED***
	// nolint: errcheck
	os.RemoveAll(d.path)
***REMOVED***

// Join returns a new path with this directory as the base of the path
func (d *Dir) Join(parts ...string) string ***REMOVED***
	return filepath.Join(append([]string***REMOVED***d.Path()***REMOVED***, parts...)...)
***REMOVED***
