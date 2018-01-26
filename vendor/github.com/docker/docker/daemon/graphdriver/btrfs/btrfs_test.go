// +build linux

package btrfs

import (
	"os"
	"path"
	"testing"

	"github.com/docker/docker/daemon/graphdriver/graphtest"
)

// This avoids creating a new driver for each test if all tests are run
// Make sure to put new tests between TestBtrfsSetup and TestBtrfsTeardown
func TestBtrfsSetup(t *testing.T) ***REMOVED***
	graphtest.GetDriver(t, "btrfs")
***REMOVED***

func TestBtrfsCreateEmpty(t *testing.T) ***REMOVED***
	graphtest.DriverTestCreateEmpty(t, "btrfs")
***REMOVED***

func TestBtrfsCreateBase(t *testing.T) ***REMOVED***
	graphtest.DriverTestCreateBase(t, "btrfs")
***REMOVED***

func TestBtrfsCreateSnap(t *testing.T) ***REMOVED***
	graphtest.DriverTestCreateSnap(t, "btrfs")
***REMOVED***

func TestBtrfsSubvolDelete(t *testing.T) ***REMOVED***
	d := graphtest.GetDriver(t, "btrfs")
	if err := d.CreateReadWrite("test", "", nil); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer graphtest.PutDriver(t)

	dirFS, err := d.Get("test", "")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer d.Put("test")

	dir := dirFS.Path()

	if err := subvolCreate(dir, "subvoltest"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if _, err := os.Stat(path.Join(dir, "subvoltest")); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := d.Remove("test"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if _, err := os.Stat(path.Join(dir, "subvoltest")); !os.IsNotExist(err) ***REMOVED***
		t.Fatalf("expected not exist error on nested subvol, got: %v", err)
	***REMOVED***
***REMOVED***

func TestBtrfsTeardown(t *testing.T) ***REMOVED***
	graphtest.PutDriver(t)
***REMOVED***
