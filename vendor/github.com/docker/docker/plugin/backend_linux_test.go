package plugin

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestAtomicRemoveAllNormal(t *testing.T) ***REMOVED***
	dir, err := ioutil.TempDir("", "atomic-remove-with-normal")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(dir) // just try to make sure this gets cleaned up

	if err := atomicRemoveAll(dir); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if _, err := os.Stat(dir); !os.IsNotExist(err) ***REMOVED***
		t.Fatalf("dir should be gone: %v", err)
	***REMOVED***
	if _, err := os.Stat(dir + "-removing"); !os.IsNotExist(err) ***REMOVED***
		t.Fatalf("dir should be gone: %v", err)
	***REMOVED***
***REMOVED***

func TestAtomicRemoveAllAlreadyExists(t *testing.T) ***REMOVED***
	dir, err := ioutil.TempDir("", "atomic-remove-already-exists")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(dir) // just try to make sure this gets cleaned up

	if err := os.MkdirAll(dir+"-removing", 0755); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(dir + "-removing")

	if err := atomicRemoveAll(dir); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if _, err := os.Stat(dir); !os.IsNotExist(err) ***REMOVED***
		t.Fatalf("dir should be gone: %v", err)
	***REMOVED***
	if _, err := os.Stat(dir + "-removing"); !os.IsNotExist(err) ***REMOVED***
		t.Fatalf("dir should be gone: %v", err)
	***REMOVED***
***REMOVED***

func TestAtomicRemoveAllNotExist(t *testing.T) ***REMOVED***
	if err := atomicRemoveAll("/not-exist"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	dir, err := ioutil.TempDir("", "atomic-remove-already-exists")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(dir) // just try to make sure this gets cleaned up

	// create the removing dir, but not the "real" one
	foo := filepath.Join(dir, "foo")
	removing := dir + "-removing"
	if err := os.MkdirAll(removing, 0755); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := atomicRemoveAll(dir); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if _, err := os.Stat(foo); !os.IsNotExist(err) ***REMOVED***
		t.Fatalf("dir should be gone: %v", err)
	***REMOVED***
	if _, err := os.Stat(removing); !os.IsNotExist(err) ***REMOVED***
		t.Fatalf("dir should be gone: %v", err)
	***REMOVED***
***REMOVED***
