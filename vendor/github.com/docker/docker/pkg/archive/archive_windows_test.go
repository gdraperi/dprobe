// +build windows

package archive

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestCopyFileWithInvalidDest(t *testing.T) ***REMOVED***
	// TODO Windows: This is currently failing. Not sure what has
	// recently changed in CopyWithTar as used to pass. Further investigation
	// is required.
	t.Skip("Currently fails")
	folder, err := ioutil.TempDir("", "docker-archive-test")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(folder)
	dest := "c:dest"
	srcFolder := filepath.Join(folder, "src")
	src := filepath.Join(folder, "src", "src")
	err = os.MkdirAll(srcFolder, 0740)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	ioutil.WriteFile(src, []byte("content"), 0777)
	err = defaultCopyWithTar(src, dest)
	if err == nil ***REMOVED***
		t.Fatalf("archiver.CopyWithTar should throw an error on invalid dest.")
	***REMOVED***
***REMOVED***

func TestCanonicalTarNameForPath(t *testing.T) ***REMOVED***
	cases := []struct ***REMOVED***
		in, expected string
		shouldFail   bool
	***REMOVED******REMOVED***
		***REMOVED***"foo", "foo", false***REMOVED***,
		***REMOVED***"foo/bar", "___", true***REMOVED***, // unix-styled windows path must fail
		***REMOVED***`foo\bar`, "foo/bar", false***REMOVED***,
	***REMOVED***
	for _, v := range cases ***REMOVED***
		if out, err := CanonicalTarNameForPath(v.in); err != nil && !v.shouldFail ***REMOVED***
			t.Fatalf("cannot get canonical name for path: %s: %v", v.in, err)
		***REMOVED*** else if v.shouldFail && err == nil ***REMOVED***
			t.Fatalf("canonical path call should have failed with error. in=%s out=%s", v.in, out)
		***REMOVED*** else if !v.shouldFail && out != v.expected ***REMOVED***
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
		***REMOVED***`foo\bar`, false, "foo/bar"***REMOVED***,
		***REMOVED***`foo\bar`, true, "foo/bar/"***REMOVED***,
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
		***REMOVED***0000, 0111***REMOVED***,
		***REMOVED***0777, 0755***REMOVED***,
		***REMOVED***0644, 0755***REMOVED***,
		***REMOVED***0755, 0755***REMOVED***,
		***REMOVED***0444, 0555***REMOVED***,
		***REMOVED***0755 | os.ModeDir, 0755 | os.ModeDir***REMOVED***,
		***REMOVED***0755 | os.ModeSymlink, 0755 | os.ModeSymlink***REMOVED***,
	***REMOVED***
	for _, v := range cases ***REMOVED***
		if out := chmodTarEntry(v.in); out != v.expected ***REMOVED***
			t.Fatalf("wrong chmod. expected:%v got:%v", v.expected, out)
		***REMOVED***
	***REMOVED***
***REMOVED***
