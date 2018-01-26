package layer

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/containerd/continuity/driver"
	"github.com/docker/docker/daemon/graphdriver"
	"github.com/docker/docker/daemon/graphdriver/vfs"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/containerfs"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/stringid"
	"github.com/opencontainers/go-digest"
)

func init() ***REMOVED***
	graphdriver.ApplyUncompressedLayer = archive.UnpackLayer
	defaultArchiver := archive.NewDefaultArchiver()
	vfs.CopyDir = defaultArchiver.CopyWithTar
***REMOVED***

func newVFSGraphDriver(td string) (graphdriver.Driver, error) ***REMOVED***
	uidMap := []idtools.IDMap***REMOVED***
		***REMOVED***
			ContainerID: 0,
			HostID:      os.Getuid(),
			Size:        1,
		***REMOVED***,
	***REMOVED***
	gidMap := []idtools.IDMap***REMOVED***
		***REMOVED***
			ContainerID: 0,
			HostID:      os.Getgid(),
			Size:        1,
		***REMOVED***,
	***REMOVED***

	options := graphdriver.Options***REMOVED***Root: td, UIDMaps: uidMap, GIDMaps: gidMap***REMOVED***
	return graphdriver.GetDriver("vfs", nil, options)
***REMOVED***

func newTestGraphDriver(t *testing.T) (graphdriver.Driver, func()) ***REMOVED***
	td, err := ioutil.TempDir("", "graph-")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	driver, err := newVFSGraphDriver(td)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	return driver, func() ***REMOVED***
		os.RemoveAll(td)
	***REMOVED***
***REMOVED***

func newTestStore(t *testing.T) (Store, string, func()) ***REMOVED***
	td, err := ioutil.TempDir("", "layerstore-")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	graph, graphcleanup := newTestGraphDriver(t)
	fms, err := NewFSMetadataStore(td)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	ls, err := NewStoreFromGraphDriver(fms, graph, runtime.GOOS)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	return ls, td, func() ***REMOVED***
		graphcleanup()
		os.RemoveAll(td)
	***REMOVED***
***REMOVED***

type layerInit func(root containerfs.ContainerFS) error

func createLayer(ls Store, parent ChainID, layerFunc layerInit) (Layer, error) ***REMOVED***
	containerID := stringid.GenerateRandomID()
	mount, err := ls.CreateRWLayer(containerID, parent, nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	pathFS, err := mount.Mount("")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err := layerFunc(pathFS); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	ts, err := mount.TarStream()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer ts.Close()

	layer, err := ls.Register(ts, parent)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err := mount.Unmount(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if _, err := ls.ReleaseRWLayer(mount); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return layer, nil
***REMOVED***

type FileApplier interface ***REMOVED***
	ApplyFile(root containerfs.ContainerFS) error
***REMOVED***

type testFile struct ***REMOVED***
	name       string
	content    []byte
	permission os.FileMode
***REMOVED***

func newTestFile(name string, content []byte, perm os.FileMode) FileApplier ***REMOVED***
	return &testFile***REMOVED***
		name:       name,
		content:    content,
		permission: perm,
	***REMOVED***
***REMOVED***

func (tf *testFile) ApplyFile(root containerfs.ContainerFS) error ***REMOVED***
	fullPath := root.Join(root.Path(), tf.name)
	if err := root.MkdirAll(root.Dir(fullPath), 0755); err != nil ***REMOVED***
		return err
	***REMOVED***
	// Check if already exists
	if stat, err := root.Stat(fullPath); err == nil && stat.Mode().Perm() != tf.permission ***REMOVED***
		if err := root.Lchmod(fullPath, tf.permission); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return driver.WriteFile(root, fullPath, tf.content, tf.permission)
***REMOVED***

func initWithFiles(files ...FileApplier) layerInit ***REMOVED***
	return func(root containerfs.ContainerFS) error ***REMOVED***
		for _, f := range files ***REMOVED***
			if err := f.ApplyFile(root); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***
***REMOVED***

func getCachedLayer(l Layer) *roLayer ***REMOVED***
	if rl, ok := l.(*referencedCacheLayer); ok ***REMOVED***
		return rl.roLayer
	***REMOVED***
	return l.(*roLayer)
***REMOVED***

func getMountLayer(l RWLayer) *mountedLayer ***REMOVED***
	return l.(*referencedRWLayer).mountedLayer
***REMOVED***

func createMetadata(layers ...Layer) []Metadata ***REMOVED***
	metadata := make([]Metadata, len(layers))
	for i := range layers ***REMOVED***
		size, err := layers[i].Size()
		if err != nil ***REMOVED***
			panic(err)
		***REMOVED***

		metadata[i].ChainID = layers[i].ChainID()
		metadata[i].DiffID = layers[i].DiffID()
		metadata[i].Size = size
		metadata[i].DiffSize = getCachedLayer(layers[i]).size
	***REMOVED***

	return metadata
***REMOVED***

func assertMetadata(t *testing.T, metadata, expectedMetadata []Metadata) ***REMOVED***
	if len(metadata) != len(expectedMetadata) ***REMOVED***
		t.Fatalf("Unexpected number of deletes %d, expected %d", len(metadata), len(expectedMetadata))
	***REMOVED***

	for i := range metadata ***REMOVED***
		if metadata[i] != expectedMetadata[i] ***REMOVED***
			t.Errorf("Unexpected metadata\n\tExpected: %#v\n\tActual: %#v", expectedMetadata[i], metadata[i])
		***REMOVED***
	***REMOVED***
	if t.Failed() ***REMOVED***
		t.FailNow()
	***REMOVED***
***REMOVED***

func releaseAndCheckDeleted(t *testing.T, ls Store, layer Layer, removed ...Layer) ***REMOVED***
	layerCount := len(ls.(*layerStore).layerMap)
	expectedMetadata := createMetadata(removed...)
	metadata, err := ls.Release(layer)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	assertMetadata(t, metadata, expectedMetadata)

	if expected := layerCount - len(removed); len(ls.(*layerStore).layerMap) != expected ***REMOVED***
		t.Fatalf("Unexpected number of layers %d, expected %d", len(ls.(*layerStore).layerMap), expected)
	***REMOVED***
***REMOVED***

func cacheID(l Layer) string ***REMOVED***
	return getCachedLayer(l).cacheID
***REMOVED***

func assertLayerEqual(t *testing.T, l1, l2 Layer) ***REMOVED***
	if l1.ChainID() != l2.ChainID() ***REMOVED***
		t.Fatalf("Mismatched ChainID: %s vs %s", l1.ChainID(), l2.ChainID())
	***REMOVED***
	if l1.DiffID() != l2.DiffID() ***REMOVED***
		t.Fatalf("Mismatched DiffID: %s vs %s", l1.DiffID(), l2.DiffID())
	***REMOVED***

	size1, err := l1.Size()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	size2, err := l2.Size()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if size1 != size2 ***REMOVED***
		t.Fatalf("Mismatched size: %d vs %d", size1, size2)
	***REMOVED***

	if cacheID(l1) != cacheID(l2) ***REMOVED***
		t.Fatalf("Mismatched cache id: %s vs %s", cacheID(l1), cacheID(l2))
	***REMOVED***

	p1 := l1.Parent()
	p2 := l2.Parent()
	if p1 != nil && p2 != nil ***REMOVED***
		assertLayerEqual(t, p1, p2)
	***REMOVED*** else if p1 != nil || p2 != nil ***REMOVED***
		t.Fatalf("Mismatched parents: %v vs %v", p1, p2)
	***REMOVED***
***REMOVED***

func TestMountAndRegister(t *testing.T) ***REMOVED***
	ls, _, cleanup := newTestStore(t)
	defer cleanup()

	li := initWithFiles(newTestFile("testfile.txt", []byte("some test data"), 0644))
	layer, err := createLayer(ls, "", li)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	size, _ := layer.Size()
	t.Logf("Layer size: %d", size)

	mount2, err := ls.CreateRWLayer("new-test-mount", layer.ChainID(), nil)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	path2, err := mount2.Mount("")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	b, err := driver.ReadFile(path2, path2.Join(path2.Path(), "testfile.txt"))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if expected := "some test data"; string(b) != expected ***REMOVED***
		t.Fatalf("Wrong file data, expected %q, got %q", expected, string(b))
	***REMOVED***

	if err := mount2.Unmount(); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if _, err := ls.ReleaseRWLayer(mount2); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestLayerRelease(t *testing.T) ***REMOVED***
	// TODO Windows: Figure out why this is failing
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("Failing on Windows")
	***REMOVED***
	ls, _, cleanup := newTestStore(t)
	defer cleanup()

	layer1, err := createLayer(ls, "", initWithFiles(newTestFile("layer1.txt", []byte("layer 1 file"), 0644)))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	layer2, err := createLayer(ls, layer1.ChainID(), initWithFiles(newTestFile("layer2.txt", []byte("layer 2 file"), 0644)))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if _, err := ls.Release(layer1); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	layer3a, err := createLayer(ls, layer2.ChainID(), initWithFiles(newTestFile("layer3.txt", []byte("layer 3a file"), 0644)))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	layer3b, err := createLayer(ls, layer2.ChainID(), initWithFiles(newTestFile("layer3.txt", []byte("layer 3b file"), 0644)))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if _, err := ls.Release(layer2); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	t.Logf("Layer1:  %s", layer1.ChainID())
	t.Logf("Layer2:  %s", layer2.ChainID())
	t.Logf("Layer3a: %s", layer3a.ChainID())
	t.Logf("Layer3b: %s", layer3b.ChainID())

	if expected := 4; len(ls.(*layerStore).layerMap) != expected ***REMOVED***
		t.Fatalf("Unexpected number of layers %d, expected %d", len(ls.(*layerStore).layerMap), expected)
	***REMOVED***

	releaseAndCheckDeleted(t, ls, layer3b, layer3b)
	releaseAndCheckDeleted(t, ls, layer3a, layer3a, layer2, layer1)
***REMOVED***

func TestStoreRestore(t *testing.T) ***REMOVED***
	// TODO Windows: Figure out why this is failing
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("Failing on Windows")
	***REMOVED***
	ls, _, cleanup := newTestStore(t)
	defer cleanup()

	layer1, err := createLayer(ls, "", initWithFiles(newTestFile("layer1.txt", []byte("layer 1 file"), 0644)))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	layer2, err := createLayer(ls, layer1.ChainID(), initWithFiles(newTestFile("layer2.txt", []byte("layer 2 file"), 0644)))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if _, err := ls.Release(layer1); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	layer3, err := createLayer(ls, layer2.ChainID(), initWithFiles(newTestFile("layer3.txt", []byte("layer 3 file"), 0644)))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if _, err := ls.Release(layer2); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	m, err := ls.CreateRWLayer("some-mount_name", layer3.ChainID(), nil)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	pathFS, err := m.Mount("")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := driver.WriteFile(pathFS, pathFS.Join(pathFS.Path(), "testfile.txt"), []byte("nothing here"), 0644); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := m.Unmount(); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	ls2, err := NewStoreFromGraphDriver(ls.(*layerStore).store, ls.(*layerStore).driver, runtime.GOOS)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	layer3b, err := ls2.Get(layer3.ChainID())
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	assertLayerEqual(t, layer3b, layer3)

	// Create again with same name, should return error
	if _, err := ls2.CreateRWLayer("some-mount_name", layer3b.ChainID(), nil); err == nil ***REMOVED***
		t.Fatal("Expected error creating mount with same name")
	***REMOVED*** else if err != ErrMountNameConflict ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	m2, err := ls2.GetRWLayer("some-mount_name")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if mountPath, err := m2.Mount(""); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED*** else if pathFS.Path() != mountPath.Path() ***REMOVED***
		t.Fatalf("Unexpected path %s, expected %s", mountPath.Path(), pathFS.Path())
	***REMOVED***

	if mountPath, err := m2.Mount(""); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED*** else if pathFS.Path() != mountPath.Path() ***REMOVED***
		t.Fatalf("Unexpected path %s, expected %s", mountPath.Path(), pathFS.Path())
	***REMOVED***
	if err := m2.Unmount(); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	b, err := driver.ReadFile(pathFS, pathFS.Join(pathFS.Path(), "testfile.txt"))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if expected := "nothing here"; string(b) != expected ***REMOVED***
		t.Fatalf("Unexpected content %q, expected %q", string(b), expected)
	***REMOVED***

	if err := m2.Unmount(); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if metadata, err := ls2.ReleaseRWLayer(m2); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED*** else if len(metadata) != 0 ***REMOVED***
		t.Fatalf("Unexpectedly deleted layers: %#v", metadata)
	***REMOVED***

	if metadata, err := ls2.ReleaseRWLayer(m2); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED*** else if len(metadata) != 0 ***REMOVED***
		t.Fatalf("Unexpectedly deleted layers: %#v", metadata)
	***REMOVED***

	releaseAndCheckDeleted(t, ls2, layer3b, layer3, layer2, layer1)
***REMOVED***

func TestTarStreamStability(t *testing.T) ***REMOVED***
	// TODO Windows: Figure out why this is failing
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("Failing on Windows")
	***REMOVED***
	ls, _, cleanup := newTestStore(t)
	defer cleanup()

	files1 := []FileApplier***REMOVED***
		newTestFile("/etc/hosts", []byte("mydomain 10.0.0.1"), 0644),
		newTestFile("/etc/profile", []byte("PATH=/usr/bin"), 0644),
	***REMOVED***
	addedFile := newTestFile("/etc/shadow", []byte("root:::::::"), 0644)
	files2 := []FileApplier***REMOVED***
		newTestFile("/etc/hosts", []byte("mydomain 10.0.0.2"), 0644),
		newTestFile("/etc/profile", []byte("PATH=/usr/bin"), 0664),
		newTestFile("/root/.bashrc", []byte("PATH=/usr/sbin:/usr/bin"), 0644),
	***REMOVED***

	tar1, err := tarFromFiles(files1...)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	tar2, err := tarFromFiles(files2...)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	layer1, err := ls.Register(bytes.NewReader(tar1), "")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	// hack layer to add file
	p, err := ls.(*layerStore).driver.Get(layer1.(*referencedCacheLayer).cacheID, "")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := addedFile.ApplyFile(p); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := ls.(*layerStore).driver.Put(layer1.(*referencedCacheLayer).cacheID); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	layer2, err := ls.Register(bytes.NewReader(tar2), layer1.ChainID())
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	id1 := layer1.ChainID()
	t.Logf("Layer 1: %s", layer1.ChainID())
	t.Logf("Layer 2: %s", layer2.ChainID())

	if _, err := ls.Release(layer1); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	assertLayerDiff(t, tar2, layer2)

	layer1b, err := ls.Get(id1)
	if err != nil ***REMOVED***
		t.Logf("Content of layer map: %#v", ls.(*layerStore).layerMap)
		t.Fatal(err)
	***REMOVED***

	if _, err := ls.Release(layer2); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	assertLayerDiff(t, tar1, layer1b)

	if _, err := ls.Release(layer1b); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func assertLayerDiff(t *testing.T, expected []byte, layer Layer) ***REMOVED***
	expectedDigest := digest.FromBytes(expected)

	if digest.Digest(layer.DiffID()) != expectedDigest ***REMOVED***
		t.Fatalf("Mismatched diff id for %s, got %s, expected %s", layer.ChainID(), layer.DiffID(), expected)
	***REMOVED***

	ts, err := layer.TarStream()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer ts.Close()

	actual, err := ioutil.ReadAll(ts)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if len(actual) != len(expected) ***REMOVED***
		logByteDiff(t, actual, expected)
		t.Fatalf("Mismatched tar stream size for %s, got %d, expected %d", layer.ChainID(), len(actual), len(expected))
	***REMOVED***

	actualDigest := digest.FromBytes(actual)

	if actualDigest != expectedDigest ***REMOVED***
		logByteDiff(t, actual, expected)
		t.Fatalf("Wrong digest of tar stream, got %s, expected %s", actualDigest, expectedDigest)
	***REMOVED***
***REMOVED***

const maxByteLog = 4 * 1024

func logByteDiff(t *testing.T, actual, expected []byte) ***REMOVED***
	d1, d2 := byteDiff(actual, expected)
	if len(d1) == 0 && len(d2) == 0 ***REMOVED***
		return
	***REMOVED***

	prefix := len(actual) - len(d1)
	if len(d1) > maxByteLog || len(d2) > maxByteLog ***REMOVED***
		t.Logf("Byte diff after %d matching bytes", prefix)
	***REMOVED*** else ***REMOVED***
		t.Logf("Byte diff after %d matching bytes\nActual bytes after prefix:\n%x\nExpected bytes after prefix:\n%x", prefix, d1, d2)
	***REMOVED***
***REMOVED***

// byteDiff returns the differing bytes after the matching prefix
func byteDiff(b1, b2 []byte) ([]byte, []byte) ***REMOVED***
	i := 0
	for i < len(b1) && i < len(b2) ***REMOVED***
		if b1[i] != b2[i] ***REMOVED***
			break
		***REMOVED***
		i++
	***REMOVED***

	return b1[i:], b2[i:]
***REMOVED***

func tarFromFiles(files ...FileApplier) ([]byte, error) ***REMOVED***
	td, err := ioutil.TempDir("", "tar-")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer os.RemoveAll(td)

	for _, f := range files ***REMOVED***
		if err := f.ApplyFile(containerfs.NewLocalContainerFS(td)); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	r, err := archive.Tar(td, archive.Uncompressed)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, r); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return buf.Bytes(), nil
***REMOVED***

// assertReferences asserts that all the references are to the same
// image and represent the full set of references to that image.
func assertReferences(t *testing.T, references ...Layer) ***REMOVED***
	if len(references) == 0 ***REMOVED***
		return
	***REMOVED***
	base := references[0].(*referencedCacheLayer).roLayer
	seenReferences := map[Layer]struct***REMOVED******REMOVED******REMOVED***
		references[0]: ***REMOVED******REMOVED***,
	***REMOVED***
	for i := 1; i < len(references); i++ ***REMOVED***
		other := references[i].(*referencedCacheLayer).roLayer
		if base != other ***REMOVED***
			t.Fatalf("Unexpected referenced cache layer %s, expecting %s", other.ChainID(), base.ChainID())
		***REMOVED***
		if _, ok := base.references[references[i]]; !ok ***REMOVED***
			t.Fatalf("Reference not part of reference list: %v", references[i])
		***REMOVED***
		if _, ok := seenReferences[references[i]]; ok ***REMOVED***
			t.Fatalf("Duplicated reference %v", references[i])
		***REMOVED***
	***REMOVED***
	if rc := len(base.references); rc != len(references) ***REMOVED***
		t.Fatalf("Unexpected number of references %d, expecting %d", rc, len(references))
	***REMOVED***
***REMOVED***

func TestRegisterExistingLayer(t *testing.T) ***REMOVED***
	ls, _, cleanup := newTestStore(t)
	defer cleanup()

	baseFiles := []FileApplier***REMOVED***
		newTestFile("/etc/profile", []byte("# Base configuration"), 0644),
	***REMOVED***

	layerFiles := []FileApplier***REMOVED***
		newTestFile("/root/.bashrc", []byte("# Root configuration"), 0644),
	***REMOVED***

	li := initWithFiles(baseFiles...)
	layer1, err := createLayer(ls, "", li)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	tar1, err := tarFromFiles(layerFiles...)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	layer2a, err := ls.Register(bytes.NewReader(tar1), layer1.ChainID())
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	layer2b, err := ls.Register(bytes.NewReader(tar1), layer1.ChainID())
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	assertReferences(t, layer2a, layer2b)
***REMOVED***

func TestTarStreamVerification(t *testing.T) ***REMOVED***
	// TODO Windows: Figure out why this is failing
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("Failing on Windows")
	***REMOVED***
	ls, tmpdir, cleanup := newTestStore(t)
	defer cleanup()

	files1 := []FileApplier***REMOVED***
		newTestFile("/foo", []byte("abc"), 0644),
		newTestFile("/bar", []byte("def"), 0644),
	***REMOVED***
	files2 := []FileApplier***REMOVED***
		newTestFile("/foo", []byte("abc"), 0644),
		newTestFile("/bar", []byte("def"), 0600), // different perm
	***REMOVED***

	tar1, err := tarFromFiles(files1...)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	tar2, err := tarFromFiles(files2...)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	layer1, err := ls.Register(bytes.NewReader(tar1), "")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	layer2, err := ls.Register(bytes.NewReader(tar2), "")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	id1 := digest.Digest(layer1.ChainID())
	id2 := digest.Digest(layer2.ChainID())

	// Replace tar data files
	src, err := os.Open(filepath.Join(tmpdir, id1.Algorithm().String(), id1.Hex(), "tar-split.json.gz"))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer src.Close()

	dst, err := os.Create(filepath.Join(tmpdir, id2.Algorithm().String(), id2.Hex(), "tar-split.json.gz"))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	src.Sync()
	dst.Sync()

	ts, err := layer2.TarStream()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	_, err = io.Copy(ioutil.Discard, ts)
	if err == nil ***REMOVED***
		t.Fatal("expected data verification to fail")
	***REMOVED***
	if !strings.Contains(err.Error(), "could not verify layer data") ***REMOVED***
		t.Fatalf("wrong error returned from tarstream: %q", err)
	***REMOVED***
***REMOVED***
