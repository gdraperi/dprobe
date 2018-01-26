package ioutils

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

var (
	testMode os.FileMode = 0640
)

func init() ***REMOVED***
	// Windows does not support full Linux file mode
	if runtime.GOOS == "windows" ***REMOVED***
		testMode = 0666
	***REMOVED***
***REMOVED***

func TestAtomicWriteToFile(t *testing.T) ***REMOVED***
	tmpDir, err := ioutil.TempDir("", "atomic-writers-test")
	if err != nil ***REMOVED***
		t.Fatalf("Error when creating temporary directory: %s", err)
	***REMOVED***
	defer os.RemoveAll(tmpDir)

	expected := []byte("barbaz")
	if err := AtomicWriteFile(filepath.Join(tmpDir, "foo"), expected, testMode); err != nil ***REMOVED***
		t.Fatalf("Error writing to file: %v", err)
	***REMOVED***

	actual, err := ioutil.ReadFile(filepath.Join(tmpDir, "foo"))
	if err != nil ***REMOVED***
		t.Fatalf("Error reading from file: %v", err)
	***REMOVED***

	if !bytes.Equal(actual, expected) ***REMOVED***
		t.Fatalf("Data mismatch, expected %q, got %q", expected, actual)
	***REMOVED***

	st, err := os.Stat(filepath.Join(tmpDir, "foo"))
	if err != nil ***REMOVED***
		t.Fatalf("Error statting file: %v", err)
	***REMOVED***
	if expected := os.FileMode(testMode); st.Mode() != expected ***REMOVED***
		t.Fatalf("Mode mismatched, expected %o, got %o", expected, st.Mode())
	***REMOVED***
***REMOVED***

func TestAtomicWriteSetCommit(t *testing.T) ***REMOVED***
	tmpDir, err := ioutil.TempDir("", "atomic-writerset-test")
	if err != nil ***REMOVED***
		t.Fatalf("Error when creating temporary directory: %s", err)
	***REMOVED***
	defer os.RemoveAll(tmpDir)

	if err := os.Mkdir(filepath.Join(tmpDir, "tmp"), 0700); err != nil ***REMOVED***
		t.Fatalf("Error creating tmp directory: %s", err)
	***REMOVED***

	targetDir := filepath.Join(tmpDir, "target")
	ws, err := NewAtomicWriteSet(filepath.Join(tmpDir, "tmp"))
	if err != nil ***REMOVED***
		t.Fatalf("Error creating atomic write set: %s", err)
	***REMOVED***

	expected := []byte("barbaz")
	if err := ws.WriteFile("foo", expected, testMode); err != nil ***REMOVED***
		t.Fatalf("Error writing to file: %v", err)
	***REMOVED***

	if _, err := ioutil.ReadFile(filepath.Join(targetDir, "foo")); err == nil ***REMOVED***
		t.Fatalf("Expected error reading file where should not exist")
	***REMOVED***

	if err := ws.Commit(targetDir); err != nil ***REMOVED***
		t.Fatalf("Error committing file: %s", err)
	***REMOVED***

	actual, err := ioutil.ReadFile(filepath.Join(targetDir, "foo"))
	if err != nil ***REMOVED***
		t.Fatalf("Error reading from file: %v", err)
	***REMOVED***

	if !bytes.Equal(actual, expected) ***REMOVED***
		t.Fatalf("Data mismatch, expected %q, got %q", expected, actual)
	***REMOVED***

	st, err := os.Stat(filepath.Join(targetDir, "foo"))
	if err != nil ***REMOVED***
		t.Fatalf("Error statting file: %v", err)
	***REMOVED***
	if expected := os.FileMode(testMode); st.Mode() != expected ***REMOVED***
		t.Fatalf("Mode mismatched, expected %o, got %o", expected, st.Mode())
	***REMOVED***

***REMOVED***

func TestAtomicWriteSetCancel(t *testing.T) ***REMOVED***
	tmpDir, err := ioutil.TempDir("", "atomic-writerset-test")
	if err != nil ***REMOVED***
		t.Fatalf("Error when creating temporary directory: %s", err)
	***REMOVED***
	defer os.RemoveAll(tmpDir)

	if err := os.Mkdir(filepath.Join(tmpDir, "tmp"), 0700); err != nil ***REMOVED***
		t.Fatalf("Error creating tmp directory: %s", err)
	***REMOVED***

	ws, err := NewAtomicWriteSet(filepath.Join(tmpDir, "tmp"))
	if err != nil ***REMOVED***
		t.Fatalf("Error creating atomic write set: %s", err)
	***REMOVED***

	expected := []byte("barbaz")
	if err := ws.WriteFile("foo", expected, testMode); err != nil ***REMOVED***
		t.Fatalf("Error writing to file: %v", err)
	***REMOVED***

	if err := ws.Cancel(); err != nil ***REMOVED***
		t.Fatalf("Error committing file: %s", err)
	***REMOVED***

	if _, err := ioutil.ReadFile(filepath.Join(tmpDir, "target", "foo")); err == nil ***REMOVED***
		t.Fatalf("Expected error reading file where should not exist")
	***REMOVED*** else if !os.IsNotExist(err) ***REMOVED***
		t.Fatalf("Unexpected error reading file: %s", err)
	***REMOVED***
***REMOVED***
