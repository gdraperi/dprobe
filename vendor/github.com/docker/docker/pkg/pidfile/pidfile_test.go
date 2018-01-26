package pidfile

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestNewAndRemove(t *testing.T) ***REMOVED***
	dir, err := ioutil.TempDir(os.TempDir(), "test-pidfile")
	if err != nil ***REMOVED***
		t.Fatal("Could not create test directory")
	***REMOVED***

	path := filepath.Join(dir, "testfile")
	file, err := New(path)
	if err != nil ***REMOVED***
		t.Fatal("Could not create test file", err)
	***REMOVED***

	_, err = New(path)
	if err == nil ***REMOVED***
		t.Fatal("Test file creation not blocked")
	***REMOVED***

	if err := file.Remove(); err != nil ***REMOVED***
		t.Fatal("Could not delete created test file")
	***REMOVED***
***REMOVED***

func TestRemoveInvalidPath(t *testing.T) ***REMOVED***
	file := PIDFile***REMOVED***path: filepath.Join("foo", "bar")***REMOVED***

	if err := file.Remove(); err == nil ***REMOVED***
		t.Fatal("Non-existing file doesn't give an error on delete")
	***REMOVED***
***REMOVED***
