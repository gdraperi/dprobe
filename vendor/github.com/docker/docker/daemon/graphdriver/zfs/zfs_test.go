// +build linux

package zfs

import (
	"testing"

	"github.com/docker/docker/daemon/graphdriver/graphtest"
)

// This avoids creating a new driver for each test if all tests are run
// Make sure to put new tests between TestZfsSetup and TestZfsTeardown
func TestZfsSetup(t *testing.T) ***REMOVED***
	graphtest.GetDriver(t, "zfs")
***REMOVED***

func TestZfsCreateEmpty(t *testing.T) ***REMOVED***
	graphtest.DriverTestCreateEmpty(t, "zfs")
***REMOVED***

func TestZfsCreateBase(t *testing.T) ***REMOVED***
	graphtest.DriverTestCreateBase(t, "zfs")
***REMOVED***

func TestZfsCreateSnap(t *testing.T) ***REMOVED***
	graphtest.DriverTestCreateSnap(t, "zfs")
***REMOVED***

func TestZfsSetQuota(t *testing.T) ***REMOVED***
	graphtest.DriverTestSetQuota(t, "zfs", true)
***REMOVED***

func TestZfsTeardown(t *testing.T) ***REMOVED***
	graphtest.PutDriver(t)
***REMOVED***
