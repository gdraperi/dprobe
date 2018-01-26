package archive

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/docker/docker/pkg/system"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"
)

// setupOverlayTestDir creates files in a directory with overlay whiteouts
// Tree layout
// .
// ├── d1     # opaque, 0700
// │   └── f1 # empty file, 0600
// ├── d2     # opaque, 0750
// │   └── f1 # empty file, 0660
// └── d3     # 0700
//     └── f1 # whiteout, 0644
func setupOverlayTestDir(t *testing.T, src string) ***REMOVED***
	// Create opaque directory containing single file and permission 0700
	err := os.Mkdir(filepath.Join(src, "d1"), 0700)
	require.NoError(t, err)

	err = system.Lsetxattr(filepath.Join(src, "d1"), "trusted.overlay.opaque", []byte("y"), 0)
	require.NoError(t, err)

	err = ioutil.WriteFile(filepath.Join(src, "d1", "f1"), []byte***REMOVED******REMOVED***, 0600)
	require.NoError(t, err)

	// Create another opaque directory containing single file but with permission 0750
	err = os.Mkdir(filepath.Join(src, "d2"), 0750)
	require.NoError(t, err)

	err = system.Lsetxattr(filepath.Join(src, "d2"), "trusted.overlay.opaque", []byte("y"), 0)
	require.NoError(t, err)

	err = ioutil.WriteFile(filepath.Join(src, "d2", "f1"), []byte***REMOVED******REMOVED***, 0660)
	require.NoError(t, err)

	// Create regular directory with deleted file
	err = os.Mkdir(filepath.Join(src, "d3"), 0700)
	require.NoError(t, err)

	err = system.Mknod(filepath.Join(src, "d3", "f1"), unix.S_IFCHR, 0)
	require.NoError(t, err)
***REMOVED***

func checkOpaqueness(t *testing.T, path string, opaque string) ***REMOVED***
	xattrOpaque, err := system.Lgetxattr(path, "trusted.overlay.opaque")
	require.NoError(t, err)

	if string(xattrOpaque) != opaque ***REMOVED***
		t.Fatalf("Unexpected opaque value: %q, expected %q", string(xattrOpaque), opaque)
	***REMOVED***

***REMOVED***

func checkOverlayWhiteout(t *testing.T, path string) ***REMOVED***
	stat, err := os.Stat(path)
	require.NoError(t, err)

	statT, ok := stat.Sys().(*syscall.Stat_t)
	if !ok ***REMOVED***
		t.Fatalf("Unexpected type: %t, expected *syscall.Stat_t", stat.Sys())
	***REMOVED***
	if statT.Rdev != 0 ***REMOVED***
		t.Fatalf("Non-zero device number for whiteout")
	***REMOVED***
***REMOVED***

func checkFileMode(t *testing.T, path string, perm os.FileMode) ***REMOVED***
	stat, err := os.Stat(path)
	require.NoError(t, err)

	if stat.Mode() != perm ***REMOVED***
		t.Fatalf("Unexpected file mode for %s: %o, expected %o", path, stat.Mode(), perm)
	***REMOVED***
***REMOVED***

func TestOverlayTarUntar(t *testing.T) ***REMOVED***
	oldmask, err := system.Umask(0)
	require.NoError(t, err)
	defer system.Umask(oldmask)

	src, err := ioutil.TempDir("", "docker-test-overlay-tar-src")
	require.NoError(t, err)
	defer os.RemoveAll(src)

	setupOverlayTestDir(t, src)

	dst, err := ioutil.TempDir("", "docker-test-overlay-tar-dst")
	require.NoError(t, err)
	defer os.RemoveAll(dst)

	options := &TarOptions***REMOVED***
		Compression:    Uncompressed,
		WhiteoutFormat: OverlayWhiteoutFormat,
	***REMOVED***
	archive, err := TarWithOptions(src, options)
	require.NoError(t, err)
	defer archive.Close()

	err = Untar(archive, dst, options)
	require.NoError(t, err)

	checkFileMode(t, filepath.Join(dst, "d1"), 0700|os.ModeDir)
	checkFileMode(t, filepath.Join(dst, "d2"), 0750|os.ModeDir)
	checkFileMode(t, filepath.Join(dst, "d3"), 0700|os.ModeDir)
	checkFileMode(t, filepath.Join(dst, "d1", "f1"), 0600)
	checkFileMode(t, filepath.Join(dst, "d2", "f1"), 0660)
	checkFileMode(t, filepath.Join(dst, "d3", "f1"), os.ModeCharDevice|os.ModeDevice)

	checkOpaqueness(t, filepath.Join(dst, "d1"), "y")
	checkOpaqueness(t, filepath.Join(dst, "d2"), "y")
	checkOpaqueness(t, filepath.Join(dst, "d3"), "")
	checkOverlayWhiteout(t, filepath.Join(dst, "d3", "f1"))
***REMOVED***

func TestOverlayTarAUFSUntar(t *testing.T) ***REMOVED***
	oldmask, err := system.Umask(0)
	require.NoError(t, err)
	defer system.Umask(oldmask)

	src, err := ioutil.TempDir("", "docker-test-overlay-tar-src")
	require.NoError(t, err)
	defer os.RemoveAll(src)

	setupOverlayTestDir(t, src)

	dst, err := ioutil.TempDir("", "docker-test-overlay-tar-dst")
	require.NoError(t, err)
	defer os.RemoveAll(dst)

	archive, err := TarWithOptions(src, &TarOptions***REMOVED***
		Compression:    Uncompressed,
		WhiteoutFormat: OverlayWhiteoutFormat,
	***REMOVED***)
	require.NoError(t, err)
	defer archive.Close()

	err = Untar(archive, dst, &TarOptions***REMOVED***
		Compression:    Uncompressed,
		WhiteoutFormat: AUFSWhiteoutFormat,
	***REMOVED***)
	require.NoError(t, err)

	checkFileMode(t, filepath.Join(dst, "d1"), 0700|os.ModeDir)
	checkFileMode(t, filepath.Join(dst, "d1", WhiteoutOpaqueDir), 0700)
	checkFileMode(t, filepath.Join(dst, "d2"), 0750|os.ModeDir)
	checkFileMode(t, filepath.Join(dst, "d2", WhiteoutOpaqueDir), 0750)
	checkFileMode(t, filepath.Join(dst, "d3"), 0700|os.ModeDir)
	checkFileMode(t, filepath.Join(dst, "d1", "f1"), 0600)
	checkFileMode(t, filepath.Join(dst, "d2", "f1"), 0660)
	checkFileMode(t, filepath.Join(dst, "d3", WhiteoutPrefix+"f1"), 0600)
***REMOVED***
