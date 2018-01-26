package dockerfile

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// createTestTempDir creates a temporary directory for testing.
// It returns the created path and a cleanup function which is meant to be used as deferred call.
// When an error occurs, it terminates the test.
func createTestTempDir(t *testing.T, dir, prefix string) (string, func()) ***REMOVED***
	path, err := ioutil.TempDir(dir, prefix)

	if err != nil ***REMOVED***
		t.Fatalf("Error when creating directory %s with prefix %s: %s", dir, prefix, err)
	***REMOVED***

	return path, func() ***REMOVED***
		err = os.RemoveAll(path)

		if err != nil ***REMOVED***
			t.Fatalf("Error when removing directory %s: %s", path, err)
		***REMOVED***
	***REMOVED***
***REMOVED***

// createTestTempFile creates a temporary file within dir with specific contents and permissions.
// When an error occurs, it terminates the test
func createTestTempFile(t *testing.T, dir, filename, contents string, perm os.FileMode) string ***REMOVED***
	filePath := filepath.Join(dir, filename)
	err := ioutil.WriteFile(filePath, []byte(contents), perm)

	if err != nil ***REMOVED***
		t.Fatalf("Error when creating %s file: %s", filename, err)
	***REMOVED***

	return filePath
***REMOVED***

// createTestSymlink creates a symlink file within dir which points to oldname
func createTestSymlink(t *testing.T, dir, filename, oldname string) string ***REMOVED***
	filePath := filepath.Join(dir, filename)
	if err := os.Symlink(oldname, filePath); err != nil ***REMOVED***
		t.Fatalf("Error when creating %s symlink to %s: %s", filename, oldname, err)
	***REMOVED***

	return filePath
***REMOVED***
