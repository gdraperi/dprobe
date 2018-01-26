package layer

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/docker/docker/daemon/graphdriver"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/stringid"
	"github.com/vbatts/tar-split/tar/asm"
	"github.com/vbatts/tar-split/tar/storage"
)

func writeTarSplitFile(name string, tarContent []byte) error ***REMOVED***
	f, err := os.OpenFile(name, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer f.Close()

	fz := gzip.NewWriter(f)

	metaPacker := storage.NewJSONPacker(fz)
	defer fz.Close()

	rdr, err := asm.NewInputTarStream(bytes.NewReader(tarContent), metaPacker, nil)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if _, err := io.Copy(ioutil.Discard, rdr); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func TestLayerMigration(t *testing.T) ***REMOVED***
	// TODO Windows: Figure out why this is failing
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("Failing on Windows")
	***REMOVED***
	td, err := ioutil.TempDir("", "migration-test-")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(td)

	layer1Files := []FileApplier***REMOVED***
		newTestFile("/root/.bashrc", []byte("# Boring configuration"), 0644),
		newTestFile("/etc/profile", []byte("# Base configuration"), 0644),
	***REMOVED***

	layer2Files := []FileApplier***REMOVED***
		newTestFile("/root/.bashrc", []byte("# Updated configuration"), 0644),
	***REMOVED***

	tar1, err := tarFromFiles(layer1Files...)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	tar2, err := tarFromFiles(layer2Files...)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	graph, err := newVFSGraphDriver(filepath.Join(td, "graphdriver-"))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	graphID1 := stringid.GenerateRandomID()
	if err := graph.Create(graphID1, "", nil); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if _, err := graph.ApplyDiff(graphID1, "", bytes.NewReader(tar1)); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	tf1 := filepath.Join(td, "tar1.json.gz")
	if err := writeTarSplitFile(tf1, tar1); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	fms, err := NewFSMetadataStore(filepath.Join(td, "layers"))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	ls, err := NewStoreFromGraphDriver(fms, graph, runtime.GOOS)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	newTarDataPath := filepath.Join(td, ".migration-tardata")
	diffID, size, err := ls.(*layerStore).ChecksumForGraphID(graphID1, "", tf1, newTarDataPath)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	layer1a, err := ls.(*layerStore).RegisterByGraphID(graphID1, "", diffID, newTarDataPath, size)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	layer1b, err := ls.Register(bytes.NewReader(tar1), "")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	assertReferences(t, layer1a, layer1b)
	// Attempt register, should be same
	layer2a, err := ls.Register(bytes.NewReader(tar2), layer1a.ChainID())
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	graphID2 := stringid.GenerateRandomID()
	if err := graph.Create(graphID2, graphID1, nil); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if _, err := graph.ApplyDiff(graphID2, graphID1, bytes.NewReader(tar2)); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	tf2 := filepath.Join(td, "tar2.json.gz")
	if err := writeTarSplitFile(tf2, tar2); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	diffID, size, err = ls.(*layerStore).ChecksumForGraphID(graphID2, graphID1, tf2, newTarDataPath)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	layer2b, err := ls.(*layerStore).RegisterByGraphID(graphID2, layer1a.ChainID(), diffID, tf2, size)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	assertReferences(t, layer2a, layer2b)

	if metadata, err := ls.Release(layer2a); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED*** else if len(metadata) > 0 ***REMOVED***
		t.Fatalf("Unexpected layer removal after first release: %#v", metadata)
	***REMOVED***

	metadata, err := ls.Release(layer2b)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	assertMetadata(t, metadata, createMetadata(layer2a))
***REMOVED***

func tarFromFilesInGraph(graph graphdriver.Driver, graphID, parentID string, files ...FileApplier) ([]byte, error) ***REMOVED***
	t, err := tarFromFiles(files...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if err := graph.Create(graphID, parentID, nil); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if _, err := graph.ApplyDiff(graphID, parentID, bytes.NewReader(t)); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	ar, err := graph.Diff(graphID, parentID)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer ar.Close()

	return ioutil.ReadAll(ar)
***REMOVED***

func TestLayerMigrationNoTarsplit(t *testing.T) ***REMOVED***
	// TODO Windows: Figure out why this is failing
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("Failing on Windows")
	***REMOVED***
	td, err := ioutil.TempDir("", "migration-test-")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(td)

	layer1Files := []FileApplier***REMOVED***
		newTestFile("/root/.bashrc", []byte("# Boring configuration"), 0644),
		newTestFile("/etc/profile", []byte("# Base configuration"), 0644),
	***REMOVED***

	layer2Files := []FileApplier***REMOVED***
		newTestFile("/root/.bashrc", []byte("# Updated configuration"), 0644),
	***REMOVED***

	graph, err := newVFSGraphDriver(filepath.Join(td, "graphdriver-"))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	graphID1 := stringid.GenerateRandomID()
	graphID2 := stringid.GenerateRandomID()

	tar1, err := tarFromFilesInGraph(graph, graphID1, "", layer1Files...)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	tar2, err := tarFromFilesInGraph(graph, graphID2, graphID1, layer2Files...)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	fms, err := NewFSMetadataStore(filepath.Join(td, "layers"))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	ls, err := NewStoreFromGraphDriver(fms, graph, runtime.GOOS)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	newTarDataPath := filepath.Join(td, ".migration-tardata")
	diffID, size, err := ls.(*layerStore).ChecksumForGraphID(graphID1, "", "", newTarDataPath)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	layer1a, err := ls.(*layerStore).RegisterByGraphID(graphID1, "", diffID, newTarDataPath, size)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	layer1b, err := ls.Register(bytes.NewReader(tar1), "")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	assertReferences(t, layer1a, layer1b)

	// Attempt register, should be same
	layer2a, err := ls.Register(bytes.NewReader(tar2), layer1a.ChainID())
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	diffID, size, err = ls.(*layerStore).ChecksumForGraphID(graphID2, graphID1, "", newTarDataPath)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	layer2b, err := ls.(*layerStore).RegisterByGraphID(graphID2, layer1a.ChainID(), diffID, newTarDataPath, size)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	assertReferences(t, layer2a, layer2b)

	if metadata, err := ls.Release(layer2a); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED*** else if len(metadata) > 0 ***REMOVED***
		t.Fatalf("Unexpected layer removal after first release: %#v", metadata)
	***REMOVED***

	metadata, err := ls.Release(layer2b)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	assertMetadata(t, metadata, createMetadata(layer2a))
***REMOVED***

func TestMountMigration(t *testing.T) ***REMOVED***
	// TODO Windows: Figure out why this is failing (obvious - paths... needs porting)
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("Failing on Windows")
	***REMOVED***
	ls, _, cleanup := newTestStore(t)
	defer cleanup()

	baseFiles := []FileApplier***REMOVED***
		newTestFile("/root/.bashrc", []byte("# Boring configuration"), 0644),
		newTestFile("/etc/profile", []byte("# Base configuration"), 0644),
	***REMOVED***
	initFiles := []FileApplier***REMOVED***
		newTestFile("/etc/hosts", []byte***REMOVED******REMOVED***, 0644),
		newTestFile("/etc/resolv.conf", []byte***REMOVED******REMOVED***, 0644),
	***REMOVED***
	mountFiles := []FileApplier***REMOVED***
		newTestFile("/etc/hosts", []byte("localhost 127.0.0.1"), 0644),
		newTestFile("/root/.bashrc", []byte("# Updated configuration"), 0644),
		newTestFile("/root/testfile1.txt", []byte("nothing valuable"), 0644),
	***REMOVED***

	initTar, err := tarFromFiles(initFiles...)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	mountTar, err := tarFromFiles(mountFiles...)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	graph := ls.(*layerStore).driver

	layer1, err := createLayer(ls, "", initWithFiles(baseFiles...))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	graphID1 := layer1.(*referencedCacheLayer).cacheID

	containerID := stringid.GenerateRandomID()
	containerInit := fmt.Sprintf("%s-init", containerID)

	if err := graph.Create(containerInit, graphID1, nil); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if _, err := graph.ApplyDiff(containerInit, graphID1, bytes.NewReader(initTar)); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := graph.Create(containerID, containerInit, nil); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if _, err := graph.ApplyDiff(containerID, containerInit, bytes.NewReader(mountTar)); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := ls.(*layerStore).CreateRWLayerByGraphID("migration-mount", containerID, layer1.ChainID()); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	rwLayer1, err := ls.GetRWLayer("migration-mount")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if _, err := rwLayer1.Mount(""); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	changes, err := rwLayer1.Changes()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if expected := 5; len(changes) != expected ***REMOVED***
		t.Logf("Changes %#v", changes)
		t.Fatalf("Wrong number of changes %d, expected %d", len(changes), expected)
	***REMOVED***

	sortChanges(changes)

	assertChange(t, changes[0], archive.Change***REMOVED***
		Path: "/etc",
		Kind: archive.ChangeModify,
	***REMOVED***)
	assertChange(t, changes[1], archive.Change***REMOVED***
		Path: "/etc/hosts",
		Kind: archive.ChangeModify,
	***REMOVED***)
	assertChange(t, changes[2], archive.Change***REMOVED***
		Path: "/root",
		Kind: archive.ChangeModify,
	***REMOVED***)
	assertChange(t, changes[3], archive.Change***REMOVED***
		Path: "/root/.bashrc",
		Kind: archive.ChangeModify,
	***REMOVED***)
	assertChange(t, changes[4], archive.Change***REMOVED***
		Path: "/root/testfile1.txt",
		Kind: archive.ChangeAdd,
	***REMOVED***)

	if _, err := ls.CreateRWLayer("migration-mount", layer1.ChainID(), nil); err == nil ***REMOVED***
		t.Fatal("Expected error creating mount with same name")
	***REMOVED*** else if err != ErrMountNameConflict ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	rwLayer2, err := ls.GetRWLayer("migration-mount")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if getMountLayer(rwLayer1) != getMountLayer(rwLayer2) ***REMOVED***
		t.Fatal("Expected same layer from get with same name as from migrate")
	***REMOVED***

	if _, err := rwLayer2.Mount(""); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if _, err := rwLayer2.Mount(""); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if metadata, err := ls.Release(layer1); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED*** else if len(metadata) > 0 ***REMOVED***
		t.Fatalf("Expected no layers to be deleted, deleted %#v", metadata)
	***REMOVED***

	if err := rwLayer1.Unmount(); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if _, err := ls.ReleaseRWLayer(rwLayer1); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	if err := rwLayer2.Unmount(); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := rwLayer2.Unmount(); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	metadata, err := ls.ReleaseRWLayer(rwLayer2)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if len(metadata) == 0 ***REMOVED***
		t.Fatal("Expected base layer to be deleted when deleting mount")
	***REMOVED***

	assertMetadata(t, metadata, createMetadata(layer1))
***REMOVED***
