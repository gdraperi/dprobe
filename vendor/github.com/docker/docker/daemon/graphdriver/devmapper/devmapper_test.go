// +build linux

package devmapper

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"

	"github.com/docker/docker/daemon/graphdriver"
	"github.com/docker/docker/daemon/graphdriver/graphtest"
	"github.com/docker/docker/pkg/parsers/kernel"
	"golang.org/x/sys/unix"
)

func init() ***REMOVED***
	// Reduce the size of the base fs and loopback for the tests
	defaultDataLoopbackSize = 300 * 1024 * 1024
	defaultMetaDataLoopbackSize = 200 * 1024 * 1024
	defaultBaseFsSize = 300 * 1024 * 1024
	defaultUdevSyncOverride = true
	if err := initLoopbacks(); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
***REMOVED***

// initLoopbacks ensures that the loopback devices are properly created within
// the system running the device mapper tests.
func initLoopbacks() error ***REMOVED***
	statT, err := getBaseLoopStats()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	// create at least 8 loopback files, ya, that is a good number
	for i := 0; i < 8; i++ ***REMOVED***
		loopPath := fmt.Sprintf("/dev/loop%d", i)
		// only create new loopback files if they don't exist
		if _, err := os.Stat(loopPath); err != nil ***REMOVED***
			if mkerr := syscall.Mknod(loopPath,
				uint32(statT.Mode|syscall.S_IFBLK), int((7<<8)|(i&0xff)|((i&0xfff00)<<12))); mkerr != nil ***REMOVED***
				return mkerr
			***REMOVED***
			os.Chown(loopPath, int(statT.Uid), int(statT.Gid))
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// getBaseLoopStats inspects /dev/loop0 to collect uid,gid, and mode for the
// loop0 device on the system.  If it does not exist we assume 0,0,0660 for the
// stat data
func getBaseLoopStats() (*syscall.Stat_t, error) ***REMOVED***
	loop0, err := os.Stat("/dev/loop0")
	if err != nil ***REMOVED***
		if os.IsNotExist(err) ***REMOVED***
			return &syscall.Stat_t***REMOVED***
				Uid:  0,
				Gid:  0,
				Mode: 0660,
			***REMOVED***, nil
		***REMOVED***
		return nil, err
	***REMOVED***
	return loop0.Sys().(*syscall.Stat_t), nil
***REMOVED***

// This avoids creating a new driver for each test if all tests are run
// Make sure to put new tests between TestDevmapperSetup and TestDevmapperTeardown
func TestDevmapperSetup(t *testing.T) ***REMOVED***
	graphtest.GetDriver(t, "devicemapper")
***REMOVED***

func TestDevmapperCreateEmpty(t *testing.T) ***REMOVED***
	graphtest.DriverTestCreateEmpty(t, "devicemapper")
***REMOVED***

func TestDevmapperCreateBase(t *testing.T) ***REMOVED***
	graphtest.DriverTestCreateBase(t, "devicemapper")
***REMOVED***

func TestDevmapperCreateSnap(t *testing.T) ***REMOVED***
	graphtest.DriverTestCreateSnap(t, "devicemapper")
***REMOVED***

func TestDevmapperTeardown(t *testing.T) ***REMOVED***
	graphtest.PutDriver(t)
***REMOVED***

func TestDevmapperReduceLoopBackSize(t *testing.T) ***REMOVED***
	tenMB := int64(10 * 1024 * 1024)
	testChangeLoopBackSize(t, -tenMB, defaultDataLoopbackSize, defaultMetaDataLoopbackSize)
***REMOVED***

func TestDevmapperIncreaseLoopBackSize(t *testing.T) ***REMOVED***
	tenMB := int64(10 * 1024 * 1024)
	testChangeLoopBackSize(t, tenMB, defaultDataLoopbackSize+tenMB, defaultMetaDataLoopbackSize+tenMB)
***REMOVED***

func testChangeLoopBackSize(t *testing.T, delta, expectDataSize, expectMetaDataSize int64) ***REMOVED***
	driver := graphtest.GetDriver(t, "devicemapper").(*graphtest.Driver).Driver.(*graphdriver.NaiveDiffDriver).ProtoDriver.(*Driver)
	defer graphtest.PutDriver(t)
	// make sure data or metadata loopback size are the default size
	if s := driver.DeviceSet.Status(); s.Data.Total != uint64(defaultDataLoopbackSize) || s.Metadata.Total != uint64(defaultMetaDataLoopbackSize) ***REMOVED***
		t.Fatal("data or metadata loop back size is incorrect")
	***REMOVED***
	if err := driver.Cleanup(); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	//Reload
	d, err := Init(driver.home, []string***REMOVED***
		fmt.Sprintf("dm.loopdatasize=%d", defaultDataLoopbackSize+delta),
		fmt.Sprintf("dm.loopmetadatasize=%d", defaultMetaDataLoopbackSize+delta),
	***REMOVED***, nil, nil)
	if err != nil ***REMOVED***
		t.Fatalf("error creating devicemapper driver: %v", err)
	***REMOVED***
	driver = d.(*graphdriver.NaiveDiffDriver).ProtoDriver.(*Driver)
	if s := driver.DeviceSet.Status(); s.Data.Total != uint64(expectDataSize) || s.Metadata.Total != uint64(expectMetaDataSize) ***REMOVED***
		t.Fatal("data or metadata loop back size is incorrect")
	***REMOVED***
	if err := driver.Cleanup(); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

// Make sure devices.Lock() has been release upon return from cleanupDeletedDevices() function
func TestDevmapperLockReleasedDeviceDeletion(t *testing.T) ***REMOVED***
	driver := graphtest.GetDriver(t, "devicemapper").(*graphtest.Driver).Driver.(*graphdriver.NaiveDiffDriver).ProtoDriver.(*Driver)
	defer graphtest.PutDriver(t)

	// Call cleanupDeletedDevices() and after the call take and release
	// DeviceSet Lock. If lock has not been released, this will hang.
	driver.DeviceSet.cleanupDeletedDevices()

	doneChan := make(chan bool)

	go func() ***REMOVED***
		driver.DeviceSet.Lock()
		defer driver.DeviceSet.Unlock()
		doneChan <- true
	***REMOVED***()

	select ***REMOVED***
	case <-time.After(time.Second * 5):
		// Timer expired. That means lock was not released upon
		// function return and we are deadlocked. Release lock
		// here so that cleanup could succeed and fail the test.
		driver.DeviceSet.Unlock()
		t.Fatal("Could not acquire devices lock after call to cleanupDeletedDevices()")
	case <-doneChan:
	***REMOVED***
***REMOVED***

// Ensure that mounts aren't leakedriver. It's non-trivial for us to test the full
// reproducer of #34573 in a unit test, but we can at least make sure that a
// simple command run in a new namespace doesn't break things horribly.
func TestDevmapperMountLeaks(t *testing.T) ***REMOVED***
	if !kernel.CheckKernelVersion(3, 18, 0) ***REMOVED***
		t.Skipf("kernel version <3.18.0 and so is missing torvalds/linux@8ed936b5671bfb33d89bc60bdcc7cf0470ba52fe.")
	***REMOVED***

	driver := graphtest.GetDriver(t, "devicemapper", "dm.use_deferred_removal=false", "dm.use_deferred_deletion=false").(*graphtest.Driver).Driver.(*graphdriver.NaiveDiffDriver).ProtoDriver.(*Driver)
	defer graphtest.PutDriver(t)

	// We need to create a new (dummy) device.
	if err := driver.Create("some-layer", "", nil); err != nil ***REMOVED***
		t.Fatalf("setting up some-layer: %v", err)
	***REMOVED***

	// Mount the device.
	_, err := driver.Get("some-layer", "")
	if err != nil ***REMOVED***
		t.Fatalf("mounting some-layer: %v", err)
	***REMOVED***

	// Create a new subprocess which will inherit our mountpoint, then
	// intentionally leak it and stick around. We can't do this entirely within
	// Go because forking and namespaces in Go are really not handled well at
	// all.
	cmd := exec.Cmd***REMOVED***
		Path: "/bin/sh",
		Args: []string***REMOVED***
			"/bin/sh", "-c",
			"mount --make-rprivate / && sleep 1000s",
		***REMOVED***,
		SysProcAttr: &syscall.SysProcAttr***REMOVED***
			Unshareflags: syscall.CLONE_NEWNS,
		***REMOVED***,
	***REMOVED***
	if err := cmd.Start(); err != nil ***REMOVED***
		t.Fatalf("starting sub-command: %v", err)
	***REMOVED***
	defer func() ***REMOVED***
		unix.Kill(cmd.Process.Pid, unix.SIGKILL)
		cmd.Wait()
	***REMOVED***()

	// Now try to "drop" the device.
	if err := driver.Put("some-layer"); err != nil ***REMOVED***
		t.Fatalf("unmounting some-layer: %v", err)
	***REMOVED***
***REMOVED***
