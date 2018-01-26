// +build linux freebsd

package system

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"testing"
)

// prepareFiles creates files for testing in the temp directory
func prepareFiles(t *testing.T) (string, string, string, string) ***REMOVED***
	dir, err := ioutil.TempDir("", "docker-system-test")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	file := filepath.Join(dir, "exist")
	if err := ioutil.WriteFile(file, []byte("hello"), 0644); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	invalid := filepath.Join(dir, "doesnt-exist")

	symlink := filepath.Join(dir, "symlink")
	if err := os.Symlink(file, symlink); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	return file, invalid, symlink, dir
***REMOVED***

func TestLUtimesNano(t *testing.T) ***REMOVED***
	file, invalid, symlink, dir := prepareFiles(t)
	defer os.RemoveAll(dir)

	before, err := os.Stat(file)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	ts := []syscall.Timespec***REMOVED******REMOVED***Sec: 0, Nsec: 0***REMOVED***, ***REMOVED***Sec: 0, Nsec: 0***REMOVED******REMOVED***
	if err := LUtimesNano(symlink, ts); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	symlinkInfo, err := os.Lstat(symlink)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if before.ModTime().Unix() == symlinkInfo.ModTime().Unix() ***REMOVED***
		t.Fatal("The modification time of the symlink should be different")
	***REMOVED***

	fileInfo, err := os.Stat(file)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if before.ModTime().Unix() != fileInfo.ModTime().Unix() ***REMOVED***
		t.Fatal("The modification time of the file should be same")
	***REMOVED***

	if err := LUtimesNano(invalid, ts); err == nil ***REMOVED***
		t.Fatal("Doesn't return an error on a non-existing file")
	***REMOVED***
***REMOVED***
