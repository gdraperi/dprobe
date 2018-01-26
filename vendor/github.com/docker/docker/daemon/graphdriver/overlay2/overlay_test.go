// +build linux

package overlay2

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/docker/docker/daemon/graphdriver"
	"github.com/docker/docker/daemon/graphdriver/graphtest"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/reexec"
	"golang.org/x/sys/unix"
)

func init() ***REMOVED***
	// Do not sure chroot to speed run time and allow archive
	// errors or hangs to be debugged directly from the test process.
	untar = archive.UntarUncompressed
	graphdriver.ApplyUncompressedLayer = archive.ApplyUncompressedLayer

	reexec.Init()
***REMOVED***

func cdMountFrom(dir, device, target, mType, label string) error ***REMOVED***
	wd, err := os.Getwd()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	os.Chdir(dir)
	defer os.Chdir(wd)

	return unix.Mount(device, target, mType, 0, label)
***REMOVED***

func skipIfNaive(t *testing.T) ***REMOVED***
	td, err := ioutil.TempDir("", "naive-check-")
	if err != nil ***REMOVED***
		t.Fatalf("Failed to create temp dir: %v", err)
	***REMOVED***
	defer os.RemoveAll(td)

	if useNaiveDiff(td) ***REMOVED***
		t.Skipf("Cannot run test with naive diff")
	***REMOVED***
***REMOVED***

// This avoids creating a new driver for each test if all tests are run
// Make sure to put new tests between TestOverlaySetup and TestOverlayTeardown
func TestOverlaySetup(t *testing.T) ***REMOVED***
	graphtest.GetDriver(t, driverName)
***REMOVED***

func TestOverlayCreateEmpty(t *testing.T) ***REMOVED***
	graphtest.DriverTestCreateEmpty(t, driverName)
***REMOVED***

func TestOverlayCreateBase(t *testing.T) ***REMOVED***
	graphtest.DriverTestCreateBase(t, driverName)
***REMOVED***

func TestOverlayCreateSnap(t *testing.T) ***REMOVED***
	graphtest.DriverTestCreateSnap(t, driverName)
***REMOVED***

func TestOverlay128LayerRead(t *testing.T) ***REMOVED***
	graphtest.DriverTestDeepLayerRead(t, 128, driverName)
***REMOVED***

func TestOverlayDiffApply10Files(t *testing.T) ***REMOVED***
	skipIfNaive(t)
	graphtest.DriverTestDiffApply(t, 10, driverName)
***REMOVED***

func TestOverlayChanges(t *testing.T) ***REMOVED***
	skipIfNaive(t)
	graphtest.DriverTestChanges(t, driverName)
***REMOVED***

func TestOverlayTeardown(t *testing.T) ***REMOVED***
	graphtest.PutDriver(t)
***REMOVED***

// Benchmarks should always setup new driver

func BenchmarkExists(b *testing.B) ***REMOVED***
	graphtest.DriverBenchExists(b, driverName)
***REMOVED***

func BenchmarkGetEmpty(b *testing.B) ***REMOVED***
	graphtest.DriverBenchGetEmpty(b, driverName)
***REMOVED***

func BenchmarkDiffBase(b *testing.B) ***REMOVED***
	graphtest.DriverBenchDiffBase(b, driverName)
***REMOVED***

func BenchmarkDiffSmallUpper(b *testing.B) ***REMOVED***
	graphtest.DriverBenchDiffN(b, 10, 10, driverName)
***REMOVED***

func BenchmarkDiff10KFileUpper(b *testing.B) ***REMOVED***
	graphtest.DriverBenchDiffN(b, 10, 10000, driverName)
***REMOVED***

func BenchmarkDiff10KFilesBottom(b *testing.B) ***REMOVED***
	graphtest.DriverBenchDiffN(b, 10000, 10, driverName)
***REMOVED***

func BenchmarkDiffApply100(b *testing.B) ***REMOVED***
	graphtest.DriverBenchDiffApplyN(b, 100, driverName)
***REMOVED***

func BenchmarkDiff20Layers(b *testing.B) ***REMOVED***
	graphtest.DriverBenchDeepLayerDiff(b, 20, driverName)
***REMOVED***

func BenchmarkRead20Layers(b *testing.B) ***REMOVED***
	graphtest.DriverBenchDeepLayerRead(b, 20, driverName)
***REMOVED***
