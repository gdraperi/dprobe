// +build linux

package vfs

import (
	"testing"

	"github.com/docker/docker/daemon/graphdriver/graphtest"

	"github.com/docker/docker/pkg/reexec"
)

func init() ***REMOVED***
	reexec.Init()
***REMOVED***

// This avoids creating a new driver for each test if all tests are run
// Make sure to put new tests between TestVfsSetup and TestVfsTeardown
func TestVfsSetup(t *testing.T) ***REMOVED***
	graphtest.GetDriver(t, "vfs")
***REMOVED***

func TestVfsCreateEmpty(t *testing.T) ***REMOVED***
	graphtest.DriverTestCreateEmpty(t, "vfs")
***REMOVED***

func TestVfsCreateBase(t *testing.T) ***REMOVED***
	graphtest.DriverTestCreateBase(t, "vfs")
***REMOVED***

func TestVfsCreateSnap(t *testing.T) ***REMOVED***
	graphtest.DriverTestCreateSnap(t, "vfs")
***REMOVED***

func TestVfsSetQuota(t *testing.T) ***REMOVED***
	graphtest.DriverTestSetQuota(t, "vfs", false)
***REMOVED***

func TestVfsTeardown(t *testing.T) ***REMOVED***
	graphtest.PutDriver(t)
***REMOVED***
