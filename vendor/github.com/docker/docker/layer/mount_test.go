package layer

import (
	"io/ioutil"
	"runtime"
	"sort"
	"testing"

	"github.com/containerd/continuity/driver"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/containerfs"
)

func TestMountInit(t *testing.T) ***REMOVED***
	// TODO Windows: Figure out why this is failing
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("Failing on Windows")
	***REMOVED***
	ls, _, cleanup := newTestStore(t)
	defer cleanup()

	basefile := newTestFile("testfile.txt", []byte("base data!"), 0644)
	initfile := newTestFile("testfile.txt", []byte("init data!"), 0777)

	li := initWithFiles(basefile)
	layer, err := createLayer(ls, "", li)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	mountInit := func(root containerfs.ContainerFS) error ***REMOVED***
		return initfile.ApplyFile(root)
	***REMOVED***

	rwLayerOpts := &CreateRWLayerOpts***REMOVED***
		InitFunc: mountInit,
	***REMOVED***
	m, err := ls.CreateRWLayer("fun-mount", layer.ChainID(), rwLayerOpts)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	pathFS, err := m.Mount("")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	fi, err := pathFS.Stat(pathFS.Join(pathFS.Path(), "testfile.txt"))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	f, err := pathFS.Open(pathFS.Join(pathFS.Path(), "testfile.txt"))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if expected := "init data!"; string(b) != expected ***REMOVED***
		t.Fatalf("Unexpected test file contents %q, expected %q", string(b), expected)
	***REMOVED***

	if fi.Mode().Perm() != 0777 ***REMOVED***
		t.Fatalf("Unexpected filemode %o, expecting %o", fi.Mode().Perm(), 0777)
	***REMOVED***
***REMOVED***

func TestMountSize(t *testing.T) ***REMOVED***
	// TODO Windows: Figure out why this is failing
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("Failing on Windows")
	***REMOVED***
	ls, _, cleanup := newTestStore(t)
	defer cleanup()

	content1 := []byte("Base contents")
	content2 := []byte("Mutable contents")
	contentInit := []byte("why am I excluded from the size ☹")

	li := initWithFiles(newTestFile("file1", content1, 0644))
	layer, err := createLayer(ls, "", li)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	mountInit := func(root containerfs.ContainerFS) error ***REMOVED***
		return newTestFile("file-init", contentInit, 0777).ApplyFile(root)
	***REMOVED***
	rwLayerOpts := &CreateRWLayerOpts***REMOVED***
		InitFunc: mountInit,
	***REMOVED***

	m, err := ls.CreateRWLayer("mount-size", layer.ChainID(), rwLayerOpts)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	pathFS, err := m.Mount("")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := driver.WriteFile(pathFS, pathFS.Join(pathFS.Path(), "file2"), content2, 0755); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	mountSize, err := m.Size()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if expected := len(content2); int(mountSize) != expected ***REMOVED***
		t.Fatalf("Unexpected mount size %d, expected %d", int(mountSize), expected)
	***REMOVED***
***REMOVED***

func TestMountChanges(t *testing.T) ***REMOVED***
	// TODO Windows: Figure out why this is failing
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("Failing on Windows")
	***REMOVED***
	ls, _, cleanup := newTestStore(t)
	defer cleanup()

	basefiles := []FileApplier***REMOVED***
		newTestFile("testfile1.txt", []byte("base data!"), 0644),
		newTestFile("testfile2.txt", []byte("base data!"), 0644),
		newTestFile("testfile3.txt", []byte("base data!"), 0644),
	***REMOVED***
	initfile := newTestFile("testfile1.txt", []byte("init data!"), 0777)

	li := initWithFiles(basefiles...)
	layer, err := createLayer(ls, "", li)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	mountInit := func(root containerfs.ContainerFS) error ***REMOVED***
		return initfile.ApplyFile(root)
	***REMOVED***
	rwLayerOpts := &CreateRWLayerOpts***REMOVED***
		InitFunc: mountInit,
	***REMOVED***

	m, err := ls.CreateRWLayer("mount-changes", layer.ChainID(), rwLayerOpts)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	pathFS, err := m.Mount("")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := pathFS.Lchmod(pathFS.Join(pathFS.Path(), "testfile1.txt"), 0755); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := driver.WriteFile(pathFS, pathFS.Join(pathFS.Path(), "testfile1.txt"), []byte("mount data!"), 0755); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := pathFS.Remove(pathFS.Join(pathFS.Path(), "testfile2.txt")); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := pathFS.Lchmod(pathFS.Join(pathFS.Path(), "testfile3.txt"), 0755); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := driver.WriteFile(pathFS, pathFS.Join(pathFS.Path(), "testfile4.txt"), []byte("mount data!"), 0644); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	changes, err := m.Changes()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if expected := 4; len(changes) != expected ***REMOVED***
		t.Fatalf("Wrong number of changes %d, expected %d", len(changes), expected)
	***REMOVED***

	sortChanges(changes)

	assertChange(t, changes[0], archive.Change***REMOVED***
		Path: "/testfile1.txt",
		Kind: archive.ChangeModify,
	***REMOVED***)
	assertChange(t, changes[1], archive.Change***REMOVED***
		Path: "/testfile2.txt",
		Kind: archive.ChangeDelete,
	***REMOVED***)
	assertChange(t, changes[2], archive.Change***REMOVED***
		Path: "/testfile3.txt",
		Kind: archive.ChangeModify,
	***REMOVED***)
	assertChange(t, changes[3], archive.Change***REMOVED***
		Path: "/testfile4.txt",
		Kind: archive.ChangeAdd,
	***REMOVED***)
***REMOVED***

func assertChange(t *testing.T, actual, expected archive.Change) ***REMOVED***
	if actual.Path != expected.Path ***REMOVED***
		t.Fatalf("Unexpected change path %s, expected %s", actual.Path, expected.Path)
	***REMOVED***
	if actual.Kind != expected.Kind ***REMOVED***
		t.Fatalf("Unexpected change type %s, expected %s", actual.Kind, expected.Kind)
	***REMOVED***
***REMOVED***

func sortChanges(changes []archive.Change) ***REMOVED***
	cs := &changeSorter***REMOVED***
		changes: changes,
	***REMOVED***
	sort.Sort(cs)
***REMOVED***

type changeSorter struct ***REMOVED***
	changes []archive.Change
***REMOVED***

func (cs *changeSorter) Len() int ***REMOVED***
	return len(cs.changes)
***REMOVED***

func (cs *changeSorter) Swap(i, j int) ***REMOVED***
	cs.changes[i], cs.changes[j] = cs.changes[j], cs.changes[i]
***REMOVED***

func (cs *changeSorter) Less(i, j int) bool ***REMOVED***
	return cs.changes[i].Path < cs.changes[j].Path
***REMOVED***
