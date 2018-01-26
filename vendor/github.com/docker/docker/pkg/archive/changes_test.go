package archive

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"runtime"
	"sort"
	"testing"
	"time"

	"github.com/docker/docker/pkg/system"
	"github.com/stretchr/testify/require"
)

func max(x, y int) int ***REMOVED***
	if x >= y ***REMOVED***
		return x
	***REMOVED***
	return y
***REMOVED***

func copyDir(src, dst string) error ***REMOVED***
	return exec.Command("cp", "-a", src, dst).Run()
***REMOVED***

type FileType uint32

const (
	Regular FileType = iota
	Dir
	Symlink
)

type FileData struct ***REMOVED***
	filetype    FileType
	path        string
	contents    string
	permissions os.FileMode
***REMOVED***

func createSampleDir(t *testing.T, root string) ***REMOVED***
	files := []FileData***REMOVED***
		***REMOVED***filetype: Regular, path: "file1", contents: "file1\n", permissions: 0600***REMOVED***,
		***REMOVED***filetype: Regular, path: "file2", contents: "file2\n", permissions: 0666***REMOVED***,
		***REMOVED***filetype: Regular, path: "file3", contents: "file3\n", permissions: 0404***REMOVED***,
		***REMOVED***filetype: Regular, path: "file4", contents: "file4\n", permissions: 0600***REMOVED***,
		***REMOVED***filetype: Regular, path: "file5", contents: "file5\n", permissions: 0600***REMOVED***,
		***REMOVED***filetype: Regular, path: "file6", contents: "file6\n", permissions: 0600***REMOVED***,
		***REMOVED***filetype: Regular, path: "file7", contents: "file7\n", permissions: 0600***REMOVED***,
		***REMOVED***filetype: Dir, path: "dir1", contents: "", permissions: 0740***REMOVED***,
		***REMOVED***filetype: Regular, path: "dir1/file1-1", contents: "file1-1\n", permissions: 01444***REMOVED***,
		***REMOVED***filetype: Regular, path: "dir1/file1-2", contents: "file1-2\n", permissions: 0666***REMOVED***,
		***REMOVED***filetype: Dir, path: "dir2", contents: "", permissions: 0700***REMOVED***,
		***REMOVED***filetype: Regular, path: "dir2/file2-1", contents: "file2-1\n", permissions: 0666***REMOVED***,
		***REMOVED***filetype: Regular, path: "dir2/file2-2", contents: "file2-2\n", permissions: 0666***REMOVED***,
		***REMOVED***filetype: Dir, path: "dir3", contents: "", permissions: 0700***REMOVED***,
		***REMOVED***filetype: Regular, path: "dir3/file3-1", contents: "file3-1\n", permissions: 0666***REMOVED***,
		***REMOVED***filetype: Regular, path: "dir3/file3-2", contents: "file3-2\n", permissions: 0666***REMOVED***,
		***REMOVED***filetype: Dir, path: "dir4", contents: "", permissions: 0700***REMOVED***,
		***REMOVED***filetype: Regular, path: "dir4/file3-1", contents: "file4-1\n", permissions: 0666***REMOVED***,
		***REMOVED***filetype: Regular, path: "dir4/file3-2", contents: "file4-2\n", permissions: 0666***REMOVED***,
		***REMOVED***filetype: Symlink, path: "symlink1", contents: "target1", permissions: 0666***REMOVED***,
		***REMOVED***filetype: Symlink, path: "symlink2", contents: "target2", permissions: 0666***REMOVED***,
		***REMOVED***filetype: Symlink, path: "symlink3", contents: root + "/file1", permissions: 0666***REMOVED***,
		***REMOVED***filetype: Symlink, path: "symlink4", contents: root + "/symlink3", permissions: 0666***REMOVED***,
		***REMOVED***filetype: Symlink, path: "dirSymlink", contents: root + "/dir1", permissions: 0740***REMOVED***,
	***REMOVED***
	provisionSampleDir(t, root, files)
***REMOVED***

func provisionSampleDir(t *testing.T, root string, files []FileData) ***REMOVED***
	now := time.Now()
	for _, info := range files ***REMOVED***
		p := path.Join(root, info.path)
		if info.filetype == Dir ***REMOVED***
			err := os.MkdirAll(p, info.permissions)
			require.NoError(t, err)
		***REMOVED*** else if info.filetype == Regular ***REMOVED***
			err := ioutil.WriteFile(p, []byte(info.contents), info.permissions)
			require.NoError(t, err)
		***REMOVED*** else if info.filetype == Symlink ***REMOVED***
			err := os.Symlink(info.contents, p)
			require.NoError(t, err)
		***REMOVED***

		if info.filetype != Symlink ***REMOVED***
			// Set a consistent ctime, atime for all files and dirs
			err := system.Chtimes(p, now, now)
			require.NoError(t, err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestChangeString(t *testing.T) ***REMOVED***
	modifyChange := Change***REMOVED***"change", ChangeModify***REMOVED***
	toString := modifyChange.String()
	if toString != "C change" ***REMOVED***
		t.Fatalf("String() of a change with ChangeModify Kind should have been %s but was %s", "C change", toString)
	***REMOVED***
	addChange := Change***REMOVED***"change", ChangeAdd***REMOVED***
	toString = addChange.String()
	if toString != "A change" ***REMOVED***
		t.Fatalf("String() of a change with ChangeAdd Kind should have been %s but was %s", "A change", toString)
	***REMOVED***
	deleteChange := Change***REMOVED***"change", ChangeDelete***REMOVED***
	toString = deleteChange.String()
	if toString != "D change" ***REMOVED***
		t.Fatalf("String() of a change with ChangeDelete Kind should have been %s but was %s", "D change", toString)
	***REMOVED***
***REMOVED***

func TestChangesWithNoChanges(t *testing.T) ***REMOVED***
	// TODO Windows. There may be a way of running this, but turning off for now
	// as createSampleDir uses symlinks.
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("symlinks on Windows")
	***REMOVED***
	rwLayer, err := ioutil.TempDir("", "docker-changes-test")
	require.NoError(t, err)
	defer os.RemoveAll(rwLayer)
	layer, err := ioutil.TempDir("", "docker-changes-test-layer")
	require.NoError(t, err)
	defer os.RemoveAll(layer)
	createSampleDir(t, layer)
	changes, err := Changes([]string***REMOVED***layer***REMOVED***, rwLayer)
	require.NoError(t, err)
	if len(changes) != 0 ***REMOVED***
		t.Fatalf("Changes with no difference should have detect no changes, but detected %d", len(changes))
	***REMOVED***
***REMOVED***

func TestChangesWithChanges(t *testing.T) ***REMOVED***
	// TODO Windows. There may be a way of running this, but turning off for now
	// as createSampleDir uses symlinks.
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("symlinks on Windows")
	***REMOVED***
	// Mock the readonly layer
	layer, err := ioutil.TempDir("", "docker-changes-test-layer")
	require.NoError(t, err)
	defer os.RemoveAll(layer)
	createSampleDir(t, layer)
	os.MkdirAll(path.Join(layer, "dir1/subfolder"), 0740)

	// Mock the RW layer
	rwLayer, err := ioutil.TempDir("", "docker-changes-test")
	require.NoError(t, err)
	defer os.RemoveAll(rwLayer)

	// Create a folder in RW layer
	dir1 := path.Join(rwLayer, "dir1")
	os.MkdirAll(dir1, 0740)
	deletedFile := path.Join(dir1, ".wh.file1-2")
	ioutil.WriteFile(deletedFile, []byte***REMOVED******REMOVED***, 0600)
	modifiedFile := path.Join(dir1, "file1-1")
	ioutil.WriteFile(modifiedFile, []byte***REMOVED***0x00***REMOVED***, 01444)
	// Let's add a subfolder for a newFile
	subfolder := path.Join(dir1, "subfolder")
	os.MkdirAll(subfolder, 0740)
	newFile := path.Join(subfolder, "newFile")
	ioutil.WriteFile(newFile, []byte***REMOVED******REMOVED***, 0740)

	changes, err := Changes([]string***REMOVED***layer***REMOVED***, rwLayer)
	require.NoError(t, err)

	expectedChanges := []Change***REMOVED***
		***REMOVED***"/dir1", ChangeModify***REMOVED***,
		***REMOVED***"/dir1/file1-1", ChangeModify***REMOVED***,
		***REMOVED***"/dir1/file1-2", ChangeDelete***REMOVED***,
		***REMOVED***"/dir1/subfolder", ChangeModify***REMOVED***,
		***REMOVED***"/dir1/subfolder/newFile", ChangeAdd***REMOVED***,
	***REMOVED***
	checkChanges(expectedChanges, changes, t)
***REMOVED***

// See https://github.com/docker/docker/pull/13590
func TestChangesWithChangesGH13590(t *testing.T) ***REMOVED***
	// TODO Windows. There may be a way of running this, but turning off for now
	// as createSampleDir uses symlinks.
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("symlinks on Windows")
	***REMOVED***
	baseLayer, err := ioutil.TempDir("", "docker-changes-test.")
	require.NoError(t, err)
	defer os.RemoveAll(baseLayer)

	dir3 := path.Join(baseLayer, "dir1/dir2/dir3")
	os.MkdirAll(dir3, 07400)

	file := path.Join(dir3, "file.txt")
	ioutil.WriteFile(file, []byte("hello"), 0666)

	layer, err := ioutil.TempDir("", "docker-changes-test2.")
	require.NoError(t, err)
	defer os.RemoveAll(layer)

	// Test creating a new file
	if err := copyDir(baseLayer+"/dir1", layer+"/"); err != nil ***REMOVED***
		t.Fatalf("Cmd failed: %q", err)
	***REMOVED***

	os.Remove(path.Join(layer, "dir1/dir2/dir3/file.txt"))
	file = path.Join(layer, "dir1/dir2/dir3/file1.txt")
	ioutil.WriteFile(file, []byte("bye"), 0666)

	changes, err := Changes([]string***REMOVED***baseLayer***REMOVED***, layer)
	require.NoError(t, err)

	expectedChanges := []Change***REMOVED***
		***REMOVED***"/dir1/dir2/dir3", ChangeModify***REMOVED***,
		***REMOVED***"/dir1/dir2/dir3/file1.txt", ChangeAdd***REMOVED***,
	***REMOVED***
	checkChanges(expectedChanges, changes, t)

	// Now test changing a file
	layer, err = ioutil.TempDir("", "docker-changes-test3.")
	require.NoError(t, err)
	defer os.RemoveAll(layer)

	if err := copyDir(baseLayer+"/dir1", layer+"/"); err != nil ***REMOVED***
		t.Fatalf("Cmd failed: %q", err)
	***REMOVED***

	file = path.Join(layer, "dir1/dir2/dir3/file.txt")
	ioutil.WriteFile(file, []byte("bye"), 0666)

	changes, err = Changes([]string***REMOVED***baseLayer***REMOVED***, layer)
	require.NoError(t, err)

	expectedChanges = []Change***REMOVED***
		***REMOVED***"/dir1/dir2/dir3/file.txt", ChangeModify***REMOVED***,
	***REMOVED***
	checkChanges(expectedChanges, changes, t)
***REMOVED***

// Create a directory, copy it, make sure we report no changes between the two
func TestChangesDirsEmpty(t *testing.T) ***REMOVED***
	// TODO Windows. There may be a way of running this, but turning off for now
	// as createSampleDir uses symlinks.
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("symlinks on Windows")
	***REMOVED***
	src, err := ioutil.TempDir("", "docker-changes-test")
	require.NoError(t, err)
	defer os.RemoveAll(src)
	createSampleDir(t, src)
	dst := src + "-copy"
	err = copyDir(src, dst)
	require.NoError(t, err)
	defer os.RemoveAll(dst)
	changes, err := ChangesDirs(dst, src)
	require.NoError(t, err)

	if len(changes) != 0 ***REMOVED***
		t.Fatalf("Reported changes for identical dirs: %v", changes)
	***REMOVED***
	os.RemoveAll(src)
	os.RemoveAll(dst)
***REMOVED***

func mutateSampleDir(t *testing.T, root string) ***REMOVED***
	// Remove a regular file
	err := os.RemoveAll(path.Join(root, "file1"))
	require.NoError(t, err)

	// Remove a directory
	err = os.RemoveAll(path.Join(root, "dir1"))
	require.NoError(t, err)

	// Remove a symlink
	err = os.RemoveAll(path.Join(root, "symlink1"))
	require.NoError(t, err)

	// Rewrite a file
	err = ioutil.WriteFile(path.Join(root, "file2"), []byte("fileNN\n"), 0777)
	require.NoError(t, err)

	// Replace a file
	err = os.RemoveAll(path.Join(root, "file3"))
	require.NoError(t, err)
	err = ioutil.WriteFile(path.Join(root, "file3"), []byte("fileMM\n"), 0404)
	require.NoError(t, err)

	// Touch file
	err = system.Chtimes(path.Join(root, "file4"), time.Now().Add(time.Second), time.Now().Add(time.Second))
	require.NoError(t, err)

	// Replace file with dir
	err = os.RemoveAll(path.Join(root, "file5"))
	require.NoError(t, err)
	err = os.MkdirAll(path.Join(root, "file5"), 0666)
	require.NoError(t, err)

	// Create new file
	err = ioutil.WriteFile(path.Join(root, "filenew"), []byte("filenew\n"), 0777)
	require.NoError(t, err)

	// Create new dir
	err = os.MkdirAll(path.Join(root, "dirnew"), 0766)
	require.NoError(t, err)

	// Create a new symlink
	err = os.Symlink("targetnew", path.Join(root, "symlinknew"))
	require.NoError(t, err)

	// Change a symlink
	err = os.RemoveAll(path.Join(root, "symlink2"))
	require.NoError(t, err)

	err = os.Symlink("target2change", path.Join(root, "symlink2"))
	require.NoError(t, err)

	// Replace dir with file
	err = os.RemoveAll(path.Join(root, "dir2"))
	require.NoError(t, err)
	err = ioutil.WriteFile(path.Join(root, "dir2"), []byte("dir2\n"), 0777)
	require.NoError(t, err)

	// Touch dir
	err = system.Chtimes(path.Join(root, "dir3"), time.Now().Add(time.Second), time.Now().Add(time.Second))
	require.NoError(t, err)
***REMOVED***

func TestChangesDirsMutated(t *testing.T) ***REMOVED***
	// TODO Windows. There may be a way of running this, but turning off for now
	// as createSampleDir uses symlinks.
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("symlinks on Windows")
	***REMOVED***
	src, err := ioutil.TempDir("", "docker-changes-test")
	require.NoError(t, err)
	createSampleDir(t, src)
	dst := src + "-copy"
	err = copyDir(src, dst)
	require.NoError(t, err)
	defer os.RemoveAll(src)
	defer os.RemoveAll(dst)

	mutateSampleDir(t, dst)

	changes, err := ChangesDirs(dst, src)
	require.NoError(t, err)

	sort.Sort(changesByPath(changes))

	expectedChanges := []Change***REMOVED***
		***REMOVED***"/dir1", ChangeDelete***REMOVED***,
		***REMOVED***"/dir2", ChangeModify***REMOVED***,
		***REMOVED***"/dirnew", ChangeAdd***REMOVED***,
		***REMOVED***"/file1", ChangeDelete***REMOVED***,
		***REMOVED***"/file2", ChangeModify***REMOVED***,
		***REMOVED***"/file3", ChangeModify***REMOVED***,
		***REMOVED***"/file4", ChangeModify***REMOVED***,
		***REMOVED***"/file5", ChangeModify***REMOVED***,
		***REMOVED***"/filenew", ChangeAdd***REMOVED***,
		***REMOVED***"/symlink1", ChangeDelete***REMOVED***,
		***REMOVED***"/symlink2", ChangeModify***REMOVED***,
		***REMOVED***"/symlinknew", ChangeAdd***REMOVED***,
	***REMOVED***

	for i := 0; i < max(len(changes), len(expectedChanges)); i++ ***REMOVED***
		if i >= len(expectedChanges) ***REMOVED***
			t.Fatalf("unexpected change %s\n", changes[i].String())
		***REMOVED***
		if i >= len(changes) ***REMOVED***
			t.Fatalf("no change for expected change %s\n", expectedChanges[i].String())
		***REMOVED***
		if changes[i].Path == expectedChanges[i].Path ***REMOVED***
			if changes[i] != expectedChanges[i] ***REMOVED***
				t.Fatalf("Wrong change for %s, expected %s, got %s\n", changes[i].Path, changes[i].String(), expectedChanges[i].String())
			***REMOVED***
		***REMOVED*** else if changes[i].Path < expectedChanges[i].Path ***REMOVED***
			t.Fatalf("unexpected change %s\n", changes[i].String())
		***REMOVED*** else ***REMOVED***
			t.Fatalf("no change for expected change %s != %s\n", expectedChanges[i].String(), changes[i].String())
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestApplyLayer(t *testing.T) ***REMOVED***
	// TODO Windows. There may be a way of running this, but turning off for now
	// as createSampleDir uses symlinks.
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("symlinks on Windows")
	***REMOVED***
	src, err := ioutil.TempDir("", "docker-changes-test")
	require.NoError(t, err)
	createSampleDir(t, src)
	defer os.RemoveAll(src)
	dst := src + "-copy"
	err = copyDir(src, dst)
	require.NoError(t, err)
	mutateSampleDir(t, dst)
	defer os.RemoveAll(dst)

	changes, err := ChangesDirs(dst, src)
	require.NoError(t, err)

	layer, err := ExportChanges(dst, changes, nil, nil)
	require.NoError(t, err)

	layerCopy, err := NewTempArchive(layer, "")
	require.NoError(t, err)

	_, err = ApplyLayer(src, layerCopy)
	require.NoError(t, err)

	changes2, err := ChangesDirs(src, dst)
	require.NoError(t, err)

	if len(changes2) != 0 ***REMOVED***
		t.Fatalf("Unexpected differences after reapplying mutation: %v", changes2)
	***REMOVED***
***REMOVED***

func TestChangesSizeWithHardlinks(t *testing.T) ***REMOVED***
	// TODO Windows. There may be a way of running this, but turning off for now
	// as createSampleDir uses symlinks.
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("hardlinks on Windows")
	***REMOVED***
	srcDir, err := ioutil.TempDir("", "docker-test-srcDir")
	require.NoError(t, err)
	defer os.RemoveAll(srcDir)

	destDir, err := ioutil.TempDir("", "docker-test-destDir")
	require.NoError(t, err)
	defer os.RemoveAll(destDir)

	creationSize, err := prepareUntarSourceDirectory(100, destDir, true)
	require.NoError(t, err)

	changes, err := ChangesDirs(destDir, srcDir)
	require.NoError(t, err)

	got := ChangesSize(destDir, changes)
	if got != int64(creationSize) ***REMOVED***
		t.Errorf("Expected %d bytes of changes, got %d", creationSize, got)
	***REMOVED***
***REMOVED***

func TestChangesSizeWithNoChanges(t *testing.T) ***REMOVED***
	size := ChangesSize("/tmp", nil)
	if size != 0 ***REMOVED***
		t.Fatalf("ChangesSizes with no changes should be 0, was %d", size)
	***REMOVED***
***REMOVED***

func TestChangesSizeWithOnlyDeleteChanges(t *testing.T) ***REMOVED***
	changes := []Change***REMOVED***
		***REMOVED***Path: "deletedPath", Kind: ChangeDelete***REMOVED***,
	***REMOVED***
	size := ChangesSize("/tmp", changes)
	if size != 0 ***REMOVED***
		t.Fatalf("ChangesSizes with only delete changes should be 0, was %d", size)
	***REMOVED***
***REMOVED***

func TestChangesSize(t *testing.T) ***REMOVED***
	parentPath, err := ioutil.TempDir("", "docker-changes-test")
	require.NoError(t, err)
	defer os.RemoveAll(parentPath)
	addition := path.Join(parentPath, "addition")
	err = ioutil.WriteFile(addition, []byte***REMOVED***0x01, 0x01, 0x01***REMOVED***, 0744)
	require.NoError(t, err)
	modification := path.Join(parentPath, "modification")
	err = ioutil.WriteFile(modification, []byte***REMOVED***0x01, 0x01, 0x01***REMOVED***, 0744)
	require.NoError(t, err)

	changes := []Change***REMOVED***
		***REMOVED***Path: "addition", Kind: ChangeAdd***REMOVED***,
		***REMOVED***Path: "modification", Kind: ChangeModify***REMOVED***,
	***REMOVED***
	size := ChangesSize(parentPath, changes)
	if size != 6 ***REMOVED***
		t.Fatalf("Expected 6 bytes of changes, got %d", size)
	***REMOVED***
***REMOVED***

func checkChanges(expectedChanges, changes []Change, t *testing.T) ***REMOVED***
	sort.Sort(changesByPath(expectedChanges))
	sort.Sort(changesByPath(changes))
	for i := 0; i < max(len(changes), len(expectedChanges)); i++ ***REMOVED***
		if i >= len(expectedChanges) ***REMOVED***
			t.Fatalf("unexpected change %s\n", changes[i].String())
		***REMOVED***
		if i >= len(changes) ***REMOVED***
			t.Fatalf("no change for expected change %s\n", expectedChanges[i].String())
		***REMOVED***
		if changes[i].Path == expectedChanges[i].Path ***REMOVED***
			if changes[i] != expectedChanges[i] ***REMOVED***
				t.Fatalf("Wrong change for %s, expected %s, got %s\n", changes[i].Path, changes[i].String(), expectedChanges[i].String())
			***REMOVED***
		***REMOVED*** else if changes[i].Path < expectedChanges[i].Path ***REMOVED***
			t.Fatalf("unexpected change %s\n", changes[i].String())
		***REMOVED*** else ***REMOVED***
			t.Fatalf("no change for expected change %s != %s\n", expectedChanges[i].String(), changes[i].String())
		***REMOVED***
	***REMOVED***
***REMOVED***
