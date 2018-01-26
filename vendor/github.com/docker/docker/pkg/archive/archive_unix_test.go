// +build !windows

package archive

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/docker/docker/pkg/system"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"
)

func TestCanonicalTarNameForPath(t *testing.T) ***REMOVED***
	cases := []struct***REMOVED*** in, expected string ***REMOVED******REMOVED***
		***REMOVED***"foo", "foo"***REMOVED***,
		***REMOVED***"foo/bar", "foo/bar"***REMOVED***,
		***REMOVED***"foo/dir/", "foo/dir/"***REMOVED***,
	***REMOVED***
	for _, v := range cases ***REMOVED***
		if out, err := CanonicalTarNameForPath(v.in); err != nil ***REMOVED***
			t.Fatalf("cannot get canonical name for path: %s: %v", v.in, err)
		***REMOVED*** else if out != v.expected ***REMOVED***
			t.Fatalf("wrong canonical tar name. expected:%s got:%s", v.expected, out)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestCanonicalTarName(t *testing.T) ***REMOVED***
	cases := []struct ***REMOVED***
		in       string
		isDir    bool
		expected string
	***REMOVED******REMOVED***
		***REMOVED***"foo", false, "foo"***REMOVED***,
		***REMOVED***"foo", true, "foo/"***REMOVED***,
		***REMOVED***"foo/bar", false, "foo/bar"***REMOVED***,
		***REMOVED***"foo/bar", true, "foo/bar/"***REMOVED***,
	***REMOVED***
	for _, v := range cases ***REMOVED***
		if out, err := canonicalTarName(v.in, v.isDir); err != nil ***REMOVED***
			t.Fatalf("cannot get canonical name for path: %s: %v", v.in, err)
		***REMOVED*** else if out != v.expected ***REMOVED***
			t.Fatalf("wrong canonical tar name. expected:%s got:%s", v.expected, out)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestChmodTarEntry(t *testing.T) ***REMOVED***
	cases := []struct ***REMOVED***
		in, expected os.FileMode
	***REMOVED******REMOVED***
		***REMOVED***0000, 0000***REMOVED***,
		***REMOVED***0777, 0777***REMOVED***,
		***REMOVED***0644, 0644***REMOVED***,
		***REMOVED***0755, 0755***REMOVED***,
		***REMOVED***0444, 0444***REMOVED***,
	***REMOVED***
	for _, v := range cases ***REMOVED***
		if out := chmodTarEntry(v.in); out != v.expected ***REMOVED***
			t.Fatalf("wrong chmod. expected:%v got:%v", v.expected, out)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestTarWithHardLink(t *testing.T) ***REMOVED***
	origin, err := ioutil.TempDir("", "docker-test-tar-hardlink")
	require.NoError(t, err)
	defer os.RemoveAll(origin)

	err = ioutil.WriteFile(filepath.Join(origin, "1"), []byte("hello world"), 0700)
	require.NoError(t, err)

	err = os.Link(filepath.Join(origin, "1"), filepath.Join(origin, "2"))
	require.NoError(t, err)

	var i1, i2 uint64
	i1, err = getNlink(filepath.Join(origin, "1"))
	require.NoError(t, err)

	// sanity check that we can hardlink
	if i1 != 2 ***REMOVED***
		t.Skipf("skipping since hardlinks don't work here; expected 2 links, got %d", i1)
	***REMOVED***

	dest, err := ioutil.TempDir("", "docker-test-tar-hardlink-dest")
	require.NoError(t, err)
	defer os.RemoveAll(dest)

	// we'll do this in two steps to separate failure
	fh, err := Tar(origin, Uncompressed)
	require.NoError(t, err)

	// ensure we can read the whole thing with no error, before writing back out
	buf, err := ioutil.ReadAll(fh)
	require.NoError(t, err)

	bRdr := bytes.NewReader(buf)
	err = Untar(bRdr, dest, &TarOptions***REMOVED***Compression: Uncompressed***REMOVED***)
	require.NoError(t, err)

	i1, err = getInode(filepath.Join(dest, "1"))
	require.NoError(t, err)

	i2, err = getInode(filepath.Join(dest, "2"))
	require.NoError(t, err)

	assert.Equal(t, i1, i2)
***REMOVED***

func TestTarWithHardLinkAndRebase(t *testing.T) ***REMOVED***
	tmpDir, err := ioutil.TempDir("", "docker-test-tar-hardlink-rebase")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	origin := filepath.Join(tmpDir, "origin")
	err = os.Mkdir(origin, 0700)
	require.NoError(t, err)

	err = ioutil.WriteFile(filepath.Join(origin, "1"), []byte("hello world"), 0700)
	require.NoError(t, err)

	err = os.Link(filepath.Join(origin, "1"), filepath.Join(origin, "2"))
	require.NoError(t, err)

	var i1, i2 uint64
	i1, err = getNlink(filepath.Join(origin, "1"))
	require.NoError(t, err)

	// sanity check that we can hardlink
	if i1 != 2 ***REMOVED***
		t.Skipf("skipping since hardlinks don't work here; expected 2 links, got %d", i1)
	***REMOVED***

	dest := filepath.Join(tmpDir, "dest")
	bRdr, err := TarResourceRebase(origin, "origin")
	require.NoError(t, err)

	dstDir, srcBase := SplitPathDirEntry(origin)
	_, dstBase := SplitPathDirEntry(dest)
	content := RebaseArchiveEntries(bRdr, srcBase, dstBase)
	err = Untar(content, dstDir, &TarOptions***REMOVED***Compression: Uncompressed, NoLchown: true, NoOverwriteDirNonDir: true***REMOVED***)
	require.NoError(t, err)

	i1, err = getInode(filepath.Join(dest, "1"))
	require.NoError(t, err)
	i2, err = getInode(filepath.Join(dest, "2"))
	require.NoError(t, err)

	assert.Equal(t, i1, i2)
***REMOVED***

func getNlink(path string) (uint64, error) ***REMOVED***
	stat, err := os.Stat(path)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	statT, ok := stat.Sys().(*syscall.Stat_t)
	if !ok ***REMOVED***
		return 0, fmt.Errorf("expected type *syscall.Stat_t, got %t", stat.Sys())
	***REMOVED***
	// We need this conversion on ARM64
	return uint64(statT.Nlink), nil
***REMOVED***

func getInode(path string) (uint64, error) ***REMOVED***
	stat, err := os.Stat(path)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	statT, ok := stat.Sys().(*syscall.Stat_t)
	if !ok ***REMOVED***
		return 0, fmt.Errorf("expected type *syscall.Stat_t, got %t", stat.Sys())
	***REMOVED***
	return statT.Ino, nil
***REMOVED***

func TestTarWithBlockCharFifo(t *testing.T) ***REMOVED***
	origin, err := ioutil.TempDir("", "docker-test-tar-hardlink")
	require.NoError(t, err)

	defer os.RemoveAll(origin)
	err = ioutil.WriteFile(filepath.Join(origin, "1"), []byte("hello world"), 0700)
	require.NoError(t, err)

	err = system.Mknod(filepath.Join(origin, "2"), unix.S_IFBLK, int(system.Mkdev(int64(12), int64(5))))
	require.NoError(t, err)
	err = system.Mknod(filepath.Join(origin, "3"), unix.S_IFCHR, int(system.Mkdev(int64(12), int64(5))))
	require.NoError(t, err)
	err = system.Mknod(filepath.Join(origin, "4"), unix.S_IFIFO, int(system.Mkdev(int64(12), int64(5))))
	require.NoError(t, err)

	dest, err := ioutil.TempDir("", "docker-test-tar-hardlink-dest")
	require.NoError(t, err)
	defer os.RemoveAll(dest)

	// we'll do this in two steps to separate failure
	fh, err := Tar(origin, Uncompressed)
	require.NoError(t, err)

	// ensure we can read the whole thing with no error, before writing back out
	buf, err := ioutil.ReadAll(fh)
	require.NoError(t, err)

	bRdr := bytes.NewReader(buf)
	err = Untar(bRdr, dest, &TarOptions***REMOVED***Compression: Uncompressed***REMOVED***)
	require.NoError(t, err)

	changes, err := ChangesDirs(origin, dest)
	require.NoError(t, err)

	if len(changes) > 0 ***REMOVED***
		t.Fatalf("Tar with special device (block, char, fifo) should keep them (recreate them when untar) : %v", changes)
	***REMOVED***
***REMOVED***

// TestTarUntarWithXattr is Unix as Lsetxattr is not supported on Windows
func TestTarUntarWithXattr(t *testing.T) ***REMOVED***
	origin, err := ioutil.TempDir("", "docker-test-untar-origin")
	require.NoError(t, err)
	defer os.RemoveAll(origin)
	err = ioutil.WriteFile(filepath.Join(origin, "1"), []byte("hello world"), 0700)
	require.NoError(t, err)

	err = ioutil.WriteFile(filepath.Join(origin, "2"), []byte("welcome!"), 0700)
	require.NoError(t, err)
	err = ioutil.WriteFile(filepath.Join(origin, "3"), []byte("will be ignored"), 0700)
	require.NoError(t, err)
	err = system.Lsetxattr(filepath.Join(origin, "2"), "security.capability", []byte***REMOVED***0x00***REMOVED***, 0)
	require.NoError(t, err)

	for _, c := range []Compression***REMOVED***
		Uncompressed,
		Gzip,
	***REMOVED*** ***REMOVED***
		changes, err := tarUntar(t, origin, &TarOptions***REMOVED***
			Compression:     c,
			ExcludePatterns: []string***REMOVED***"3"***REMOVED***,
		***REMOVED***)

		if err != nil ***REMOVED***
			t.Fatalf("Error tar/untar for compression %s: %s", c.Extension(), err)
		***REMOVED***

		if len(changes) != 1 || changes[0].Path != "/3" ***REMOVED***
			t.Fatalf("Unexpected differences after tarUntar: %v", changes)
		***REMOVED***
		capability, _ := system.Lgetxattr(filepath.Join(origin, "2"), "security.capability")
		if capability == nil && capability[0] != 0x00 ***REMOVED***
			t.Fatalf("Untar should have kept the 'security.capability' xattr.")
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestCopyInfoDestinationPathSymlink(t *testing.T) ***REMOVED***
	tmpDir, _ := getTestTempDirs(t)
	defer removeAllPaths(tmpDir)

	root := strings.TrimRight(tmpDir, "/") + "/"

	type FileTestData struct ***REMOVED***
		resource FileData
		file     string
		expected CopyInfo
	***REMOVED***

	testData := []FileTestData***REMOVED***
		//Create a directory: /tmp/archive-copy-test*/dir1
		//Test will "copy" file1 to dir1
		***REMOVED***resource: FileData***REMOVED***filetype: Dir, path: "dir1", permissions: 0740***REMOVED***, file: "file1", expected: CopyInfo***REMOVED***Path: root + "dir1/file1", Exists: false, IsDir: false***REMOVED******REMOVED***,

		//Create a symlink directory to dir1: /tmp/archive-copy-test*/dirSymlink -> dir1
		//Test will "copy" file2 to dirSymlink
		***REMOVED***resource: FileData***REMOVED***filetype: Symlink, path: "dirSymlink", contents: root + "dir1", permissions: 0600***REMOVED***, file: "file2", expected: CopyInfo***REMOVED***Path: root + "dirSymlink/file2", Exists: false, IsDir: false***REMOVED******REMOVED***,

		//Create a file in tmp directory: /tmp/archive-copy-test*/file1
		//Test to cover when the full file path already exists.
		***REMOVED***resource: FileData***REMOVED***filetype: Regular, path: "file1", permissions: 0600***REMOVED***, file: "", expected: CopyInfo***REMOVED***Path: root + "file1", Exists: true***REMOVED******REMOVED***,

		//Create a directory: /tmp/archive-copy*/dir2
		//Test to cover when the full directory path already exists
		***REMOVED***resource: FileData***REMOVED***filetype: Dir, path: "dir2", permissions: 0740***REMOVED***, file: "", expected: CopyInfo***REMOVED***Path: root + "dir2", Exists: true, IsDir: true***REMOVED******REMOVED***,

		//Create a symlink to a non-existent target: /tmp/archive-copy*/symlink1 -> noSuchTarget
		//Negative test to cover symlinking to a target that does not exit
		***REMOVED***resource: FileData***REMOVED***filetype: Symlink, path: "symlink1", contents: "noSuchTarget", permissions: 0600***REMOVED***, file: "", expected: CopyInfo***REMOVED***Path: root + "noSuchTarget", Exists: false***REMOVED******REMOVED***,

		//Create a file in tmp directory for next test: /tmp/existingfile
		***REMOVED***resource: FileData***REMOVED***filetype: Regular, path: "existingfile", permissions: 0600***REMOVED***, file: "", expected: CopyInfo***REMOVED***Path: root + "existingfile", Exists: true***REMOVED******REMOVED***,

		//Create a symlink to an existing file: /tmp/archive-copy*/symlink2 -> /tmp/existingfile
		//Test to cover when the parent directory of a new file is a symlink
		***REMOVED***resource: FileData***REMOVED***filetype: Symlink, path: "symlink2", contents: "existingfile", permissions: 0600***REMOVED***, file: "", expected: CopyInfo***REMOVED***Path: root + "existingfile", Exists: true***REMOVED******REMOVED***,
	***REMOVED***

	var dirs []FileData
	for _, data := range testData ***REMOVED***
		dirs = append(dirs, data.resource)
	***REMOVED***
	provisionSampleDir(t, tmpDir, dirs)

	for _, info := range testData ***REMOVED***
		p := filepath.Join(tmpDir, info.resource.path, info.file)
		ci, err := CopyInfoDestinationPath(p)
		assert.NoError(t, err)
		assert.Equal(t, info.expected, ci)
	***REMOVED***
***REMOVED***
