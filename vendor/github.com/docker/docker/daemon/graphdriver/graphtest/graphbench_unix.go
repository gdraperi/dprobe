// +build linux freebsd

package graphtest

import (
	"io"
	"io/ioutil"
	"testing"

	contdriver "github.com/containerd/continuity/driver"
	"github.com/docker/docker/pkg/stringid"
	"github.com/stretchr/testify/require"
)

// DriverBenchExists benchmarks calls to exist
func DriverBenchExists(b *testing.B, drivername string, driveroptions ...string) ***REMOVED***
	driver := GetDriver(b, drivername, driveroptions...)
	defer PutDriver(b)

	base := stringid.GenerateRandomID()

	if err := driver.Create(base, "", nil); err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***

	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		if !driver.Exists(base) ***REMOVED***
			b.Fatal("Newly created image doesn't exist")
		***REMOVED***
	***REMOVED***
***REMOVED***

// DriverBenchGetEmpty benchmarks calls to get on an empty layer
func DriverBenchGetEmpty(b *testing.B, drivername string, driveroptions ...string) ***REMOVED***
	driver := GetDriver(b, drivername, driveroptions...)
	defer PutDriver(b)

	base := stringid.GenerateRandomID()

	if err := driver.Create(base, "", nil); err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***

	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		_, err := driver.Get(base, "")
		b.StopTimer()
		if err != nil ***REMOVED***
			b.Fatalf("Error getting mount: %s", err)
		***REMOVED***
		if err := driver.Put(base); err != nil ***REMOVED***
			b.Fatalf("Error putting mount: %s", err)
		***REMOVED***
		b.StartTimer()
	***REMOVED***
***REMOVED***

// DriverBenchDiffBase benchmarks calls to diff on a root layer
func DriverBenchDiffBase(b *testing.B, drivername string, driveroptions ...string) ***REMOVED***
	driver := GetDriver(b, drivername, driveroptions...)
	defer PutDriver(b)

	base := stringid.GenerateRandomID()
	if err := driver.Create(base, "", nil); err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***

	if err := addFiles(driver, base, 3); err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***

	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		arch, err := driver.Diff(base, "")
		if err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
		_, err = io.Copy(ioutil.Discard, arch)
		if err != nil ***REMOVED***
			b.Fatalf("Error copying archive: %s", err)
		***REMOVED***
		arch.Close()
	***REMOVED***
***REMOVED***

// DriverBenchDiffN benchmarks calls to diff on two layers with
// a provided number of files on the lower and upper layers.
func DriverBenchDiffN(b *testing.B, bottom, top int, drivername string, driveroptions ...string) ***REMOVED***
	driver := GetDriver(b, drivername, driveroptions...)
	defer PutDriver(b)
	base := stringid.GenerateRandomID()
	upper := stringid.GenerateRandomID()
	if err := driver.Create(base, "", nil); err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***

	if err := addManyFiles(driver, base, bottom, 3); err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***

	if err := driver.Create(upper, base, nil); err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***

	if err := addManyFiles(driver, upper, top, 6); err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***
	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		arch, err := driver.Diff(upper, "")
		if err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
		_, err = io.Copy(ioutil.Discard, arch)
		if err != nil ***REMOVED***
			b.Fatalf("Error copying archive: %s", err)
		***REMOVED***
		arch.Close()
	***REMOVED***
***REMOVED***

// DriverBenchDiffApplyN benchmarks calls to diff and apply together
func DriverBenchDiffApplyN(b *testing.B, fileCount int, drivername string, driveroptions ...string) ***REMOVED***
	driver := GetDriver(b, drivername, driveroptions...)
	defer PutDriver(b)
	base := stringid.GenerateRandomID()
	upper := stringid.GenerateRandomID()
	if err := driver.Create(base, "", nil); err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***

	if err := addManyFiles(driver, base, fileCount, 3); err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***

	if err := driver.Create(upper, base, nil); err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***

	if err := addManyFiles(driver, upper, fileCount, 6); err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***
	diffSize, err := driver.DiffSize(upper, "")
	if err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***
	b.ResetTimer()
	b.StopTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		diff := stringid.GenerateRandomID()
		if err := driver.Create(diff, base, nil); err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***

		if err := checkManyFiles(driver, diff, fileCount, 3); err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***

		b.StartTimer()

		arch, err := driver.Diff(upper, "")
		if err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***

		applyDiffSize, err := driver.ApplyDiff(diff, "", arch)
		if err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***

		b.StopTimer()
		arch.Close()

		if applyDiffSize != diffSize ***REMOVED***
			// TODO: enforce this
			//b.Fatalf("Apply diff size different, got %d, expected %s", applyDiffSize, diffSize)
		***REMOVED***
		if err := checkManyFiles(driver, diff, fileCount, 6); err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
	***REMOVED***
***REMOVED***

// DriverBenchDeepLayerDiff benchmarks calls to diff on top of a given number of layers.
func DriverBenchDeepLayerDiff(b *testing.B, layerCount int, drivername string, driveroptions ...string) ***REMOVED***
	driver := GetDriver(b, drivername, driveroptions...)
	defer PutDriver(b)

	base := stringid.GenerateRandomID()
	if err := driver.Create(base, "", nil); err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***

	if err := addFiles(driver, base, 50); err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***

	topLayer, err := addManyLayers(driver, base, layerCount)
	if err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***

	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		arch, err := driver.Diff(topLayer, "")
		if err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
		_, err = io.Copy(ioutil.Discard, arch)
		if err != nil ***REMOVED***
			b.Fatalf("Error copying archive: %s", err)
		***REMOVED***
		arch.Close()
	***REMOVED***
***REMOVED***

// DriverBenchDeepLayerRead benchmarks calls to read a file under a given number of layers.
func DriverBenchDeepLayerRead(b *testing.B, layerCount int, drivername string, driveroptions ...string) ***REMOVED***
	driver := GetDriver(b, drivername, driveroptions...)
	defer PutDriver(b)

	base := stringid.GenerateRandomID()
	if err := driver.Create(base, "", nil); err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***

	content := []byte("test content")
	if err := addFile(driver, base, "testfile.txt", content); err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***

	topLayer, err := addManyLayers(driver, base, layerCount)
	if err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***

	root, err := driver.Get(topLayer, "")
	if err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***
	defer driver.Put(topLayer)

	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***

		// Read content
		c, err := contdriver.ReadFile(root, root.Join(root.Path(), "testfile.txt"))
		if err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***

		b.StopTimer()
		require.Equal(b, content, c)
		b.StartTimer()
	***REMOVED***
***REMOVED***
