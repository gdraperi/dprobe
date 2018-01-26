package archive

import (
	"archive/tar"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/docker/docker/pkg/ioutils"
)

func TestApplyLayerInvalidFilenames(t *testing.T) ***REMOVED***
	// TODO Windows: Figure out how to fix this test.
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("Passes but hits breakoutError: platform and architecture is not supported")
	***REMOVED***
	for i, headers := range [][]*tar.Header***REMOVED***
		***REMOVED***
			***REMOVED***
				Name:     "../victim/dotdot",
				Typeflag: tar.TypeReg,
				Mode:     0644,
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			***REMOVED***
				// Note the leading slash
				Name:     "/../victim/slash-dotdot",
				Typeflag: tar.TypeReg,
				Mode:     0644,
			***REMOVED***,
		***REMOVED***,
	***REMOVED*** ***REMOVED***
		if err := testBreakout("applylayer", "docker-TestApplyLayerInvalidFilenames", headers); err != nil ***REMOVED***
			t.Fatalf("i=%d. %v", i, err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestApplyLayerInvalidHardlink(t *testing.T) ***REMOVED***
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("TypeLink support on Windows")
	***REMOVED***
	for i, headers := range [][]*tar.Header***REMOVED***
		***REMOVED*** // try reading victim/hello (../)
			***REMOVED***
				Name:     "dotdot",
				Typeflag: tar.TypeLink,
				Linkname: "../victim/hello",
				Mode:     0644,
			***REMOVED***,
		***REMOVED***,
		***REMOVED*** // try reading victim/hello (/../)
			***REMOVED***
				Name:     "slash-dotdot",
				Typeflag: tar.TypeLink,
				// Note the leading slash
				Linkname: "/../victim/hello",
				Mode:     0644,
			***REMOVED***,
		***REMOVED***,
		***REMOVED*** // try writing victim/file
			***REMOVED***
				Name:     "loophole-victim",
				Typeflag: tar.TypeLink,
				Linkname: "../victim",
				Mode:     0755,
			***REMOVED***,
			***REMOVED***
				Name:     "loophole-victim/file",
				Typeflag: tar.TypeReg,
				Mode:     0644,
			***REMOVED***,
		***REMOVED***,
		***REMOVED*** // try reading victim/hello (hardlink, symlink)
			***REMOVED***
				Name:     "loophole-victim",
				Typeflag: tar.TypeLink,
				Linkname: "../victim",
				Mode:     0755,
			***REMOVED***,
			***REMOVED***
				Name:     "symlink",
				Typeflag: tar.TypeSymlink,
				Linkname: "loophole-victim/hello",
				Mode:     0644,
			***REMOVED***,
		***REMOVED***,
		***REMOVED*** // Try reading victim/hello (hardlink, hardlink)
			***REMOVED***
				Name:     "loophole-victim",
				Typeflag: tar.TypeLink,
				Linkname: "../victim",
				Mode:     0755,
			***REMOVED***,
			***REMOVED***
				Name:     "hardlink",
				Typeflag: tar.TypeLink,
				Linkname: "loophole-victim/hello",
				Mode:     0644,
			***REMOVED***,
		***REMOVED***,
		***REMOVED*** // Try removing victim directory (hardlink)
			***REMOVED***
				Name:     "loophole-victim",
				Typeflag: tar.TypeLink,
				Linkname: "../victim",
				Mode:     0755,
			***REMOVED***,
			***REMOVED***
				Name:     "loophole-victim",
				Typeflag: tar.TypeReg,
				Mode:     0644,
			***REMOVED***,
		***REMOVED***,
	***REMOVED*** ***REMOVED***
		if err := testBreakout("applylayer", "docker-TestApplyLayerInvalidHardlink", headers); err != nil ***REMOVED***
			t.Fatalf("i=%d. %v", i, err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestApplyLayerInvalidSymlink(t *testing.T) ***REMOVED***
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("TypeSymLink support on Windows")
	***REMOVED***
	for i, headers := range [][]*tar.Header***REMOVED***
		***REMOVED*** // try reading victim/hello (../)
			***REMOVED***
				Name:     "dotdot",
				Typeflag: tar.TypeSymlink,
				Linkname: "../victim/hello",
				Mode:     0644,
			***REMOVED***,
		***REMOVED***,
		***REMOVED*** // try reading victim/hello (/../)
			***REMOVED***
				Name:     "slash-dotdot",
				Typeflag: tar.TypeSymlink,
				// Note the leading slash
				Linkname: "/../victim/hello",
				Mode:     0644,
			***REMOVED***,
		***REMOVED***,
		***REMOVED*** // try writing victim/file
			***REMOVED***
				Name:     "loophole-victim",
				Typeflag: tar.TypeSymlink,
				Linkname: "../victim",
				Mode:     0755,
			***REMOVED***,
			***REMOVED***
				Name:     "loophole-victim/file",
				Typeflag: tar.TypeReg,
				Mode:     0644,
			***REMOVED***,
		***REMOVED***,
		***REMOVED*** // try reading victim/hello (symlink, symlink)
			***REMOVED***
				Name:     "loophole-victim",
				Typeflag: tar.TypeSymlink,
				Linkname: "../victim",
				Mode:     0755,
			***REMOVED***,
			***REMOVED***
				Name:     "symlink",
				Typeflag: tar.TypeSymlink,
				Linkname: "loophole-victim/hello",
				Mode:     0644,
			***REMOVED***,
		***REMOVED***,
		***REMOVED*** // try reading victim/hello (symlink, hardlink)
			***REMOVED***
				Name:     "loophole-victim",
				Typeflag: tar.TypeSymlink,
				Linkname: "../victim",
				Mode:     0755,
			***REMOVED***,
			***REMOVED***
				Name:     "hardlink",
				Typeflag: tar.TypeLink,
				Linkname: "loophole-victim/hello",
				Mode:     0644,
			***REMOVED***,
		***REMOVED***,
		***REMOVED*** // try removing victim directory (symlink)
			***REMOVED***
				Name:     "loophole-victim",
				Typeflag: tar.TypeSymlink,
				Linkname: "../victim",
				Mode:     0755,
			***REMOVED***,
			***REMOVED***
				Name:     "loophole-victim",
				Typeflag: tar.TypeReg,
				Mode:     0644,
			***REMOVED***,
		***REMOVED***,
	***REMOVED*** ***REMOVED***
		if err := testBreakout("applylayer", "docker-TestApplyLayerInvalidSymlink", headers); err != nil ***REMOVED***
			t.Fatalf("i=%d. %v", i, err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestApplyLayerWhiteouts(t *testing.T) ***REMOVED***
	// TODO Windows: Figure out why this test fails
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip("Failing on Windows")
	***REMOVED***

	wd, err := ioutil.TempDir("", "graphdriver-test-whiteouts")
	if err != nil ***REMOVED***
		return
	***REMOVED***
	defer os.RemoveAll(wd)

	base := []string***REMOVED***
		".baz",
		"bar/",
		"bar/bax",
		"bar/bay/",
		"baz",
		"foo/",
		"foo/.abc",
		"foo/.bcd/",
		"foo/.bcd/a",
		"foo/cde/",
		"foo/cde/def",
		"foo/cde/efg",
		"foo/fgh",
		"foobar",
	***REMOVED***

	type tcase struct ***REMOVED***
		change, expected []string
	***REMOVED***

	tcases := []tcase***REMOVED***
		***REMOVED***
			base,
			base,
		***REMOVED***,
		***REMOVED***
			[]string***REMOVED***
				".bay",
				".wh.baz",
				"foo/",
				"foo/.bce",
				"foo/.wh..wh..opq",
				"foo/cde/",
				"foo/cde/efg",
			***REMOVED***,
			[]string***REMOVED***
				".bay",
				".baz",
				"bar/",
				"bar/bax",
				"bar/bay/",
				"foo/",
				"foo/.bce",
				"foo/cde/",
				"foo/cde/efg",
				"foobar",
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			[]string***REMOVED***
				".bay",
				".wh..baz",
				".wh.foobar",
				"foo/",
				"foo/.abc",
				"foo/.wh.cde",
				"bar/",
			***REMOVED***,
			[]string***REMOVED***
				".bay",
				"bar/",
				"bar/bax",
				"bar/bay/",
				"foo/",
				"foo/.abc",
				"foo/.bce",
			***REMOVED***,
		***REMOVED***,
		***REMOVED***
			[]string***REMOVED***
				".abc",
				".wh..wh..opq",
				"foobar",
			***REMOVED***,
			[]string***REMOVED***
				".abc",
				"foobar",
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	for i, tc := range tcases ***REMOVED***
		l, err := makeTestLayer(tc.change)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***

		_, err = UnpackLayer(wd, l, nil)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		err = l.Close()
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***

		paths, err := readDirContents(wd)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***

		if !reflect.DeepEqual(tc.expected, paths) ***REMOVED***
			t.Fatalf("invalid files for layer %d: expected %q, got %q", i, tc.expected, paths)
		***REMOVED***
	***REMOVED***

***REMOVED***

func makeTestLayer(paths []string) (rc io.ReadCloser, err error) ***REMOVED***
	tmpDir, err := ioutil.TempDir("", "graphdriver-test-mklayer")
	if err != nil ***REMOVED***
		return
	***REMOVED***
	defer func() ***REMOVED***
		if err != nil ***REMOVED***
			os.RemoveAll(tmpDir)
		***REMOVED***
	***REMOVED***()
	for _, p := range paths ***REMOVED***
		if p[len(p)-1] == filepath.Separator ***REMOVED***
			if err = os.MkdirAll(filepath.Join(tmpDir, p), 0700); err != nil ***REMOVED***
				return
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if err = ioutil.WriteFile(filepath.Join(tmpDir, p), nil, 0600); err != nil ***REMOVED***
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***
	archive, err := Tar(tmpDir, Uncompressed)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	return ioutils.NewReadCloserWrapper(archive, func() error ***REMOVED***
		err := archive.Close()
		os.RemoveAll(tmpDir)
		return err
	***REMOVED***), nil
***REMOVED***

func readDirContents(root string) ([]string, error) ***REMOVED***
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error ***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if path == root ***REMOVED***
			return nil
		***REMOVED***
		rel, err := filepath.Rel(root, path)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if info.IsDir() ***REMOVED***
			rel = rel + "/"
		***REMOVED***
		files = append(files, rel)
		return nil
	***REMOVED***)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return files, nil
***REMOVED***
