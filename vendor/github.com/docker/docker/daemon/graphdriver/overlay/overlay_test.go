// +build linux

package overlay

import (
	"testing"

	"github.com/docker/docker/daemon/graphdriver"
	"github.com/docker/docker/daemon/graphdriver/graphtest"
	"github.com/docker/docker/pkg/archive"
)

func init() ***REMOVED***
	// Do not sure chroot to speed run time and allow archive
	// errors or hangs to be debugged directly from the test process.
	graphdriver.ApplyUncompressedLayer = archive.ApplyUncompressedLayer
***REMOVED***

// This avoids creating a new driver for each test if all tests are run
// Make sure to put new tests between TestOverlaySetup and TestOverlayTeardown
func TestOverlaySetup(t *testing.T) ***REMOVED***
	graphtest.GetDriver(t, "overlay")
***REMOVED***

func TestOverlayCreateEmpty(t *testing.T) ***REMOVED***
	graphtest.DriverTestCreateEmpty(t, "overlay")
***REMOVED***

func TestOverlayCreateBase(t *testing.T) ***REMOVED***
	graphtest.DriverTestCreateBase(t, "overlay")
***REMOVED***

func TestOverlayCreateSnap(t *testing.T) ***REMOVED***
	graphtest.DriverTestCreateSnap(t, "overlay")
***REMOVED***

func TestOverlay50LayerRead(t *testing.T) ***REMOVED***
	graphtest.DriverTestDeepLayerRead(t, 50, "overlay")
***REMOVED***

// Fails due to bug in calculating changes after apply
// likely related to https://github.com/docker/docker/issues/21555
func TestOverlayDiffApply10Files(t *testing.T) ***REMOVED***
	t.Skipf("Fails to compute changes after apply intermittently")
	graphtest.DriverTestDiffApply(t, 10, "overlay")
***REMOVED***

func TestOverlayChanges(t *testing.T) ***REMOVED***
	t.Skipf("Fails to compute changes intermittently")
	graphtest.DriverTestChanges(t, "overlay")
***REMOVED***

func TestOverlayTeardown(t *testing.T) ***REMOVED***
	graphtest.PutDriver(t)
***REMOVED***

// Benchmarks should always setup new driver

func BenchmarkExists(b *testing.B) ***REMOVED***
	graphtest.DriverBenchExists(b, "overlay")
***REMOVED***

func BenchmarkGetEmpty(b *testing.B) ***REMOVED***
	graphtest.DriverBenchGetEmpty(b, "overlay")
***REMOVED***

func BenchmarkDiffBase(b *testing.B) ***REMOVED***
	graphtest.DriverBenchDiffBase(b, "overlay")
***REMOVED***

func BenchmarkDiffSmallUpper(b *testing.B) ***REMOVED***
	graphtest.DriverBenchDiffN(b, 10, 10, "overlay")
***REMOVED***

func BenchmarkDiff10KFileUpper(b *testing.B) ***REMOVED***
	graphtest.DriverBenchDiffN(b, 10, 10000, "overlay")
***REMOVED***

func BenchmarkDiff10KFilesBottom(b *testing.B) ***REMOVED***
	graphtest.DriverBenchDiffN(b, 10000, 10, "overlay")
***REMOVED***

func BenchmarkDiffApply100(b *testing.B) ***REMOVED***
	graphtest.DriverBenchDiffApplyN(b, 100, "overlay")
***REMOVED***

func BenchmarkDiff20Layers(b *testing.B) ***REMOVED***
	graphtest.DriverBenchDeepLayerDiff(b, 20, "overlay")
***REMOVED***

func BenchmarkRead20Layers(b *testing.B) ***REMOVED***
	graphtest.DriverBenchDeepLayerRead(b, 20, "overlay")
***REMOVED***
