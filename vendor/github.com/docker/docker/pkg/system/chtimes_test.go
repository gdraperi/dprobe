package system

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// prepareTempFile creates a temporary file in a temporary directory.
func prepareTempFile(t *testing.T) (string, string) ***REMOVED***
	dir, err := ioutil.TempDir("", "docker-system-test")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	file := filepath.Join(dir, "exist")
	if err := ioutil.WriteFile(file, []byte("hello"), 0644); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	return file, dir
***REMOVED***

// TestChtimes tests Chtimes on a tempfile. Test only mTime, because aTime is OS dependent
func TestChtimes(t *testing.T) ***REMOVED***
	file, dir := prepareTempFile(t)
	defer os.RemoveAll(dir)

	beforeUnixEpochTime := time.Unix(0, 0).Add(-100 * time.Second)
	unixEpochTime := time.Unix(0, 0)
	afterUnixEpochTime := time.Unix(100, 0)
	unixMaxTime := maxTime

	// Test both aTime and mTime set to Unix Epoch
	Chtimes(file, unixEpochTime, unixEpochTime)

	f, err := os.Stat(file)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if f.ModTime() != unixEpochTime ***REMOVED***
		t.Fatalf("Expected: %s, got: %s", unixEpochTime, f.ModTime())
	***REMOVED***

	// Test aTime before Unix Epoch and mTime set to Unix Epoch
	Chtimes(file, beforeUnixEpochTime, unixEpochTime)

	f, err = os.Stat(file)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if f.ModTime() != unixEpochTime ***REMOVED***
		t.Fatalf("Expected: %s, got: %s", unixEpochTime, f.ModTime())
	***REMOVED***

	// Test aTime set to Unix Epoch and mTime before Unix Epoch
	Chtimes(file, unixEpochTime, beforeUnixEpochTime)

	f, err = os.Stat(file)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if f.ModTime() != unixEpochTime ***REMOVED***
		t.Fatalf("Expected: %s, got: %s", unixEpochTime, f.ModTime())
	***REMOVED***

	// Test both aTime and mTime set to after Unix Epoch (valid time)
	Chtimes(file, afterUnixEpochTime, afterUnixEpochTime)

	f, err = os.Stat(file)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if f.ModTime() != afterUnixEpochTime ***REMOVED***
		t.Fatalf("Expected: %s, got: %s", afterUnixEpochTime, f.ModTime())
	***REMOVED***

	// Test both aTime and mTime set to Unix max time
	Chtimes(file, unixMaxTime, unixMaxTime)

	f, err = os.Stat(file)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if f.ModTime().Truncate(time.Second) != unixMaxTime.Truncate(time.Second) ***REMOVED***
		t.Fatalf("Expected: %s, got: %s", unixMaxTime.Truncate(time.Second), f.ModTime().Truncate(time.Second))
	***REMOVED***
***REMOVED***
