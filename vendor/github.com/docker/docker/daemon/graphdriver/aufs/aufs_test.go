// +build linux

package aufs

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sync"
	"testing"

	"github.com/docker/docker/daemon/graphdriver"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/reexec"
	"github.com/docker/docker/pkg/stringid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	tmpOuter = path.Join(os.TempDir(), "aufs-tests")
	tmp      = path.Join(tmpOuter, "aufs")
)

func init() ***REMOVED***
	reexec.Init()
***REMOVED***

func testInit(dir string, t testing.TB) graphdriver.Driver ***REMOVED***
	d, err := Init(dir, nil, nil, nil)
	if err != nil ***REMOVED***
		if err == graphdriver.ErrNotSupported ***REMOVED***
			t.Skip(err)
		***REMOVED*** else ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	***REMOVED***
	return d
***REMOVED***

func driverGet(d *Driver, id string, mntLabel string) (string, error) ***REMOVED***
	mnt, err := d.Get(id, mntLabel)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return mnt.Path(), nil
***REMOVED***

func newDriver(t testing.TB) *Driver ***REMOVED***
	if err := os.MkdirAll(tmp, 0755); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	d := testInit(tmp, t)
	return d.(*Driver)
***REMOVED***

func TestNewDriver(t *testing.T) ***REMOVED***
	if err := os.MkdirAll(tmp, 0755); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	d := testInit(tmp, t)
	defer os.RemoveAll(tmp)
	if d == nil ***REMOVED***
		t.Fatal("Driver should not be nil")
	***REMOVED***
***REMOVED***

func TestAufsString(t *testing.T) ***REMOVED***
	d := newDriver(t)
	defer os.RemoveAll(tmp)

	if d.String() != "aufs" ***REMOVED***
		t.Fatalf("Expected aufs got %s", d.String())
	***REMOVED***
***REMOVED***

func TestCreateDirStructure(t *testing.T) ***REMOVED***
	newDriver(t)
	defer os.RemoveAll(tmp)

	paths := []string***REMOVED***
		"mnt",
		"layers",
		"diff",
	***REMOVED***

	for _, p := range paths ***REMOVED***
		if _, err := os.Stat(path.Join(tmp, p)); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	***REMOVED***
***REMOVED***

// We should be able to create two drivers with the same dir structure
func TestNewDriverFromExistingDir(t *testing.T) ***REMOVED***
	if err := os.MkdirAll(tmp, 0755); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	testInit(tmp, t)
	testInit(tmp, t)
	os.RemoveAll(tmp)
***REMOVED***

func TestCreateNewDir(t *testing.T) ***REMOVED***
	d := newDriver(t)
	defer os.RemoveAll(tmp)

	if err := d.Create("1", "", nil); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestCreateNewDirStructure(t *testing.T) ***REMOVED***
	d := newDriver(t)
	defer os.RemoveAll(tmp)

	if err := d.Create("1", "", nil); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	paths := []string***REMOVED***
		"mnt",
		"diff",
		"layers",
	***REMOVED***

	for _, p := range paths ***REMOVED***
		if _, err := os.Stat(path.Join(tmp, p, "1")); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestRemoveImage(t *testing.T) ***REMOVED***
	d := newDriver(t)
	defer os.RemoveAll(tmp)

	if err := d.Create("1", "", nil); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := d.Remove("1"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	paths := []string***REMOVED***
		"mnt",
		"diff",
		"layers",
	***REMOVED***

	for _, p := range paths ***REMOVED***
		if _, err := os.Stat(path.Join(tmp, p, "1")); err == nil ***REMOVED***
			t.Fatalf("Error should not be nil because dirs with id 1 should be deleted: %s", p)
		***REMOVED***
		if _, err := os.Stat(path.Join(tmp, p, "1-removing")); err == nil ***REMOVED***
			t.Fatalf("Error should not be nil because dirs with id 1-removing should be deleted: %s", p)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestGetWithoutParent(t *testing.T) ***REMOVED***
	d := newDriver(t)
	defer os.RemoveAll(tmp)

	if err := d.Create("1", "", nil); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	diffPath, err := d.Get("1", "")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	expected := path.Join(tmp, "diff", "1")
	if diffPath.Path() != expected ***REMOVED***
		t.Fatalf("Expected path %s got %s", expected, diffPath)
	***REMOVED***
***REMOVED***

func TestCleanupWithNoDirs(t *testing.T) ***REMOVED***
	d := newDriver(t)
	defer os.RemoveAll(tmp)

	err := d.Cleanup()
	assert.NoError(t, err)
***REMOVED***

func TestCleanupWithDir(t *testing.T) ***REMOVED***
	d := newDriver(t)
	defer os.RemoveAll(tmp)

	if err := d.Create("1", "", nil); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := d.Cleanup(); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestMountedFalseResponse(t *testing.T) ***REMOVED***
	d := newDriver(t)
	defer os.RemoveAll(tmp)

	err := d.Create("1", "", nil)
	require.NoError(t, err)

	response, err := d.mounted(d.getDiffPath("1"))
	require.NoError(t, err)
	assert.False(t, response)
***REMOVED***

func TestMountedTrueResponse(t *testing.T) ***REMOVED***
	d := newDriver(t)
	defer os.RemoveAll(tmp)
	defer d.Cleanup()

	err := d.Create("1", "", nil)
	require.NoError(t, err)
	err = d.Create("2", "1", nil)
	require.NoError(t, err)

	_, err = d.Get("2", "")
	require.NoError(t, err)

	response, err := d.mounted(d.pathCache["2"])
	require.NoError(t, err)
	assert.True(t, response)
***REMOVED***

func TestMountWithParent(t *testing.T) ***REMOVED***
	d := newDriver(t)
	defer os.RemoveAll(tmp)

	if err := d.Create("1", "", nil); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := d.Create("2", "1", nil); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	defer func() ***REMOVED***
		if err := d.Cleanup(); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	***REMOVED***()

	mntPath, err := d.Get("2", "")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if mntPath == nil ***REMOVED***
		t.Fatal("mntPath should not be nil")
	***REMOVED***

	expected := path.Join(tmp, "mnt", "2")
	if mntPath.Path() != expected ***REMOVED***
		t.Fatalf("Expected %s got %s", expected, mntPath.Path())
	***REMOVED***
***REMOVED***

func TestRemoveMountedDir(t *testing.T) ***REMOVED***
	d := newDriver(t)
	defer os.RemoveAll(tmp)

	if err := d.Create("1", "", nil); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := d.Create("2", "1", nil); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	defer func() ***REMOVED***
		if err := d.Cleanup(); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	***REMOVED***()

	mntPath, err := d.Get("2", "")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if mntPath == nil ***REMOVED***
		t.Fatal("mntPath should not be nil")
	***REMOVED***

	mounted, err := d.mounted(d.pathCache["2"])
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if !mounted ***REMOVED***
		t.Fatal("Dir id 2 should be mounted")
	***REMOVED***

	if err := d.Remove("2"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestCreateWithInvalidParent(t *testing.T) ***REMOVED***
	d := newDriver(t)
	defer os.RemoveAll(tmp)

	if err := d.Create("1", "docker", nil); err == nil ***REMOVED***
		t.Fatal("Error should not be nil with parent does not exist")
	***REMOVED***
***REMOVED***

func TestGetDiff(t *testing.T) ***REMOVED***
	d := newDriver(t)
	defer os.RemoveAll(tmp)

	if err := d.CreateReadWrite("1", "", nil); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	diffPath, err := driverGet(d, "1", "")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// Add a file to the diff path with a fixed size
	size := int64(1024)

	f, err := os.Create(path.Join(diffPath, "test_file"))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := f.Truncate(size); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	f.Close()

	a, err := d.Diff("1", "")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if a == nil ***REMOVED***
		t.Fatal("Archive should not be nil")
	***REMOVED***
***REMOVED***

func TestChanges(t *testing.T) ***REMOVED***
	d := newDriver(t)
	defer os.RemoveAll(tmp)

	if err := d.Create("1", "", nil); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := d.CreateReadWrite("2", "1", nil); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	defer func() ***REMOVED***
		if err := d.Cleanup(); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	***REMOVED***()

	mntPoint, err := driverGet(d, "2", "")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// Create a file to save in the mountpoint
	f, err := os.Create(path.Join(mntPoint, "test.txt"))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if _, err := f.WriteString("testline"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := f.Close(); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	changes, err := d.Changes("2", "")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if len(changes) != 1 ***REMOVED***
		t.Fatalf("Dir 2 should have one change from parent got %d", len(changes))
	***REMOVED***
	change := changes[0]

	expectedPath := "/test.txt"
	if change.Path != expectedPath ***REMOVED***
		t.Fatalf("Expected path %s got %s", expectedPath, change.Path)
	***REMOVED***

	if change.Kind != archive.ChangeAdd ***REMOVED***
		t.Fatalf("Change kind should be ChangeAdd got %s", change.Kind)
	***REMOVED***

	if err := d.CreateReadWrite("3", "2", nil); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	mntPoint, err = driverGet(d, "3", "")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// Create a file to save in the mountpoint
	f, err = os.Create(path.Join(mntPoint, "test2.txt"))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if _, err := f.WriteString("testline"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := f.Close(); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	changes, err = d.Changes("3", "2")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if len(changes) != 1 ***REMOVED***
		t.Fatalf("Dir 2 should have one change from parent got %d", len(changes))
	***REMOVED***
	change = changes[0]

	expectedPath = "/test2.txt"
	if change.Path != expectedPath ***REMOVED***
		t.Fatalf("Expected path %s got %s", expectedPath, change.Path)
	***REMOVED***

	if change.Kind != archive.ChangeAdd ***REMOVED***
		t.Fatalf("Change kind should be ChangeAdd got %s", change.Kind)
	***REMOVED***
***REMOVED***

func TestDiffSize(t *testing.T) ***REMOVED***
	d := newDriver(t)
	defer os.RemoveAll(tmp)

	if err := d.CreateReadWrite("1", "", nil); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	diffPath, err := driverGet(d, "1", "")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// Add a file to the diff path with a fixed size
	size := int64(1024)

	f, err := os.Create(path.Join(diffPath, "test_file"))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := f.Truncate(size); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	s, err := f.Stat()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	size = s.Size()
	if err := f.Close(); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	diffSize, err := d.DiffSize("1", "")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if diffSize != size ***REMOVED***
		t.Fatalf("Expected size to be %d got %d", size, diffSize)
	***REMOVED***
***REMOVED***

func TestChildDiffSize(t *testing.T) ***REMOVED***
	d := newDriver(t)
	defer os.RemoveAll(tmp)
	defer d.Cleanup()

	if err := d.CreateReadWrite("1", "", nil); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	diffPath, err := driverGet(d, "1", "")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// Add a file to the diff path with a fixed size
	size := int64(1024)

	f, err := os.Create(path.Join(diffPath, "test_file"))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := f.Truncate(size); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	s, err := f.Stat()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	size = s.Size()
	if err := f.Close(); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	diffSize, err := d.DiffSize("1", "")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if diffSize != size ***REMOVED***
		t.Fatalf("Expected size to be %d got %d", size, diffSize)
	***REMOVED***

	if err := d.Create("2", "1", nil); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	diffSize, err = d.DiffSize("2", "1")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	// The diff size for the child should be zero
	if diffSize != 0 ***REMOVED***
		t.Fatalf("Expected size to be %d got %d", 0, diffSize)
	***REMOVED***
***REMOVED***

func TestExists(t *testing.T) ***REMOVED***
	d := newDriver(t)
	defer os.RemoveAll(tmp)
	defer d.Cleanup()

	if err := d.Create("1", "", nil); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if d.Exists("none") ***REMOVED***
		t.Fatal("id none should not exist in the driver")
	***REMOVED***

	if !d.Exists("1") ***REMOVED***
		t.Fatal("id 1 should exist in the driver")
	***REMOVED***
***REMOVED***

func TestStatus(t *testing.T) ***REMOVED***
	d := newDriver(t)
	defer os.RemoveAll(tmp)
	defer d.Cleanup()

	if err := d.Create("1", "", nil); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	status := d.Status()
	assert.Len(t, status, 4)

	rootDir := status[0]
	dirs := status[2]
	if rootDir[0] != "Root Dir" ***REMOVED***
		t.Fatalf("Expected Root Dir got %s", rootDir[0])
	***REMOVED***
	if rootDir[1] != d.rootPath() ***REMOVED***
		t.Fatalf("Expected %s got %s", d.rootPath(), rootDir[1])
	***REMOVED***
	if dirs[0] != "Dirs" ***REMOVED***
		t.Fatalf("Expected Dirs got %s", dirs[0])
	***REMOVED***
	if dirs[1] != "1" ***REMOVED***
		t.Fatalf("Expected 1 got %s", dirs[1])
	***REMOVED***
***REMOVED***

func TestApplyDiff(t *testing.T) ***REMOVED***
	d := newDriver(t)
	defer os.RemoveAll(tmp)
	defer d.Cleanup()

	if err := d.CreateReadWrite("1", "", nil); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	diffPath, err := driverGet(d, "1", "")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// Add a file to the diff path with a fixed size
	size := int64(1024)

	f, err := os.Create(path.Join(diffPath, "test_file"))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := f.Truncate(size); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	f.Close()

	diff, err := d.Diff("1", "")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := d.Create("2", "", nil); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := d.Create("3", "2", nil); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := d.applyDiff("3", diff); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// Ensure that the file is in the mount point for id 3

	mountPoint, err := driverGet(d, "3", "")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if _, err := os.Stat(path.Join(mountPoint, "test_file")); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func hash(c string) string ***REMOVED***
	h := sha256.New()
	fmt.Fprint(h, c)
	return hex.EncodeToString(h.Sum(nil))
***REMOVED***

func testMountMoreThan42Layers(t *testing.T, mountPath string) ***REMOVED***
	if err := os.MkdirAll(mountPath, 0755); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	defer os.RemoveAll(mountPath)
	d := testInit(mountPath, t).(*Driver)
	defer d.Cleanup()
	var last string
	var expected int

	for i := 1; i < 127; i++ ***REMOVED***
		expected++
		var (
			parent  = fmt.Sprintf("%d", i-1)
			current = fmt.Sprintf("%d", i)
		)

		if parent == "0" ***REMOVED***
			parent = ""
		***REMOVED*** else ***REMOVED***
			parent = hash(parent)
		***REMOVED***
		current = hash(current)

		err := d.CreateReadWrite(current, parent, nil)
		require.NoError(t, err, "current layer %d", i)

		point, err := driverGet(d, current, "")
		require.NoError(t, err, "current layer %d", i)

		f, err := os.Create(path.Join(point, current))
		require.NoError(t, err, "current layer %d", i)
		f.Close()

		if i%10 == 0 ***REMOVED***
			err := os.Remove(path.Join(point, parent))
			require.NoError(t, err, "current layer %d", i)
			expected--
		***REMOVED***
		last = current
	***REMOVED***

	// Perform the actual mount for the top most image
	point, err := driverGet(d, last, "")
	require.NoError(t, err)
	files, err := ioutil.ReadDir(point)
	require.NoError(t, err)
	assert.Len(t, files, expected)
***REMOVED***

func TestMountMoreThan42Layers(t *testing.T) ***REMOVED***
	defer os.RemoveAll(tmpOuter)
	testMountMoreThan42Layers(t, tmp)
***REMOVED***

func TestMountMoreThan42LayersMatchingPathLength(t *testing.T) ***REMOVED***
	defer os.RemoveAll(tmpOuter)
	zeroes := "0"
	for ***REMOVED***
		// This finds a mount path so that when combined into aufs mount options
		// 4096 byte boundary would be in between the paths or in permission
		// section. For '/tmp' it will use '/tmp/aufs-tests/00000000/aufs'
		mountPath := path.Join(tmpOuter, zeroes, "aufs")
		pathLength := 77 + len(mountPath)

		if mod := 4095 % pathLength; mod == 0 || mod > pathLength-2 ***REMOVED***
			t.Logf("Using path: %s", mountPath)
			testMountMoreThan42Layers(t, mountPath)
			return
		***REMOVED***
		zeroes += "0"
	***REMOVED***
***REMOVED***

func BenchmarkConcurrentAccess(b *testing.B) ***REMOVED***
	b.StopTimer()
	b.ResetTimer()

	d := newDriver(b)
	defer os.RemoveAll(tmp)
	defer d.Cleanup()

	numConcurrent := 256
	// create a bunch of ids
	var ids []string
	for i := 0; i < numConcurrent; i++ ***REMOVED***
		ids = append(ids, stringid.GenerateNonCryptoID())
	***REMOVED***

	if err := d.Create(ids[0], "", nil); err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***

	if err := d.Create(ids[1], ids[0], nil); err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***

	parent := ids[1]
	ids = append(ids[2:])

	chErr := make(chan error, numConcurrent)
	var outerGroup sync.WaitGroup
	outerGroup.Add(len(ids))
	b.StartTimer()

	// here's the actual bench
	for _, id := range ids ***REMOVED***
		go func(id string) ***REMOVED***
			defer outerGroup.Done()
			if err := d.Create(id, parent, nil); err != nil ***REMOVED***
				b.Logf("Create %s failed", id)
				chErr <- err
				return
			***REMOVED***
			var innerGroup sync.WaitGroup
			for i := 0; i < b.N; i++ ***REMOVED***
				innerGroup.Add(1)
				go func() ***REMOVED***
					d.Get(id, "")
					d.Put(id)
					innerGroup.Done()
				***REMOVED***()
			***REMOVED***
			innerGroup.Wait()
			d.Remove(id)
		***REMOVED***(id)
	***REMOVED***

	outerGroup.Wait()
	b.StopTimer()
	close(chErr)
	for err := range chErr ***REMOVED***
		if err != nil ***REMOVED***
			b.Log(err)
			b.Fail()
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestInitStaleCleanup(t *testing.T) ***REMOVED***
	if err := os.MkdirAll(tmp, 0755); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmp)

	for _, d := range []string***REMOVED***"diff", "mnt"***REMOVED*** ***REMOVED***
		if err := os.MkdirAll(filepath.Join(tmp, d, "123-removing"), 0755); err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
	***REMOVED***

	testInit(tmp, t)
	for _, d := range []string***REMOVED***"diff", "mnt"***REMOVED*** ***REMOVED***
		if _, err := os.Stat(filepath.Join(tmp, d, "123-removing")); err == nil ***REMOVED***
			t.Fatal("cleanup failed")
		***REMOVED***
	***REMOVED***
***REMOVED***
