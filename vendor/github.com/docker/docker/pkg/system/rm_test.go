package system

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/docker/docker/pkg/mount"
)

func TestEnsureRemoveAllNotExist(t *testing.T) ***REMOVED***
	// should never return an error for a non-existent path
	if err := EnsureRemoveAll("/non/existent/path"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestEnsureRemoveAllWithDir(t *testing.T) ***REMOVED***
	dir, err := ioutil.TempDir("", "test-ensure-removeall-with-dir")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := EnsureRemoveAll(dir); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestEnsureRemoveAllWithFile(t *testing.T) ***REMOVED***
	tmp, err := ioutil.TempFile("", "test-ensure-removeall-with-dir")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	tmp.Close()
	if err := EnsureRemoveAll(tmp.Name()); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestEnsureRemoveAllWithMount(t *testing.T) ***REMOVED***
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("mount not supported on Windows")
	***REMOVED***

	dir1, err := ioutil.TempDir("", "test-ensure-removeall-with-dir1")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	dir2, err := ioutil.TempDir("", "test-ensure-removeall-with-dir2")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(dir2)

	bindDir := filepath.Join(dir1, "bind")
	if err := os.MkdirAll(bindDir, 0755); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := mount.Mount(dir2, bindDir, "none", "bind"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	done := make(chan struct***REMOVED******REMOVED***)
	go func() ***REMOVED***
		err = EnsureRemoveAll(dir1)
		close(done)
	***REMOVED***()

	select ***REMOVED***
	case <-done:
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for EnsureRemoveAll to finish")
	***REMOVED***

	if _, err := os.Stat(dir1); !os.IsNotExist(err) ***REMOVED***
		t.Fatalf("expected %q to not exist", dir1)
	***REMOVED***
***REMOVED***
