// +build !windows

// Licensed under the Apache License, Version 2.0; See LICENSE.APACHE

package symlink

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// TODO Windows: This needs some serious work to port to Windows. For now,
// turning off testing in this package.

type dirOrLink struct ***REMOVED***
	path   string
	target string
***REMOVED***

func makeFs(tmpdir string, fs []dirOrLink) error ***REMOVED***
	for _, s := range fs ***REMOVED***
		s.path = filepath.Join(tmpdir, s.path)
		if s.target == "" ***REMOVED***
			os.MkdirAll(s.path, 0755)
			continue
		***REMOVED***
		if err := os.MkdirAll(filepath.Dir(s.path), 0755); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := os.Symlink(s.target, s.path); err != nil && !os.IsExist(err) ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func testSymlink(tmpdir, path, expected, scope string) error ***REMOVED***
	rewrite, err := FollowSymlinkInScope(filepath.Join(tmpdir, path), filepath.Join(tmpdir, scope))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	expected, err = filepath.Abs(filepath.Join(tmpdir, expected))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if expected != rewrite ***REMOVED***
		return fmt.Errorf("Expected %q got %q", expected, rewrite)
	***REMOVED***
	return nil
***REMOVED***

func TestFollowSymlinkAbsolute(t *testing.T) ***REMOVED***
	tmpdir, err := ioutil.TempDir("", "TestFollowSymlinkAbsolute")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmpdir)
	if err := makeFs(tmpdir, []dirOrLink***REMOVED******REMOVED***path: "testdata/fs/a/d", target: "/b"***REMOVED******REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := testSymlink(tmpdir, "testdata/fs/a/d/c/data", "testdata/b/c/data", "testdata"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestFollowSymlinkRelativePath(t *testing.T) ***REMOVED***
	tmpdir, err := ioutil.TempDir("", "TestFollowSymlinkRelativePath")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmpdir)
	if err := makeFs(tmpdir, []dirOrLink***REMOVED******REMOVED***path: "testdata/fs/i", target: "a"***REMOVED******REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := testSymlink(tmpdir, "testdata/fs/i", "testdata/fs/a", "testdata"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestFollowSymlinkSkipSymlinksOutsideScope(t *testing.T) ***REMOVED***
	tmpdir, err := ioutil.TempDir("", "TestFollowSymlinkSkipSymlinksOutsideScope")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmpdir)
	if err := makeFs(tmpdir, []dirOrLink***REMOVED***
		***REMOVED***path: "linkdir", target: "realdir"***REMOVED***,
		***REMOVED***path: "linkdir/foo/bar"***REMOVED***,
	***REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := testSymlink(tmpdir, "linkdir/foo/bar", "linkdir/foo/bar", "linkdir/foo"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestFollowSymlinkInvalidScopePathPair(t *testing.T) ***REMOVED***
	if _, err := FollowSymlinkInScope("toto", "testdata"); err == nil ***REMOVED***
		t.Fatal("expected an error")
	***REMOVED***
***REMOVED***

func TestFollowSymlinkLastLink(t *testing.T) ***REMOVED***
	tmpdir, err := ioutil.TempDir("", "TestFollowSymlinkLastLink")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmpdir)
	if err := makeFs(tmpdir, []dirOrLink***REMOVED******REMOVED***path: "testdata/fs/a/d", target: "/b"***REMOVED******REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := testSymlink(tmpdir, "testdata/fs/a/d", "testdata/b", "testdata"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestFollowSymlinkRelativeLinkChangeScope(t *testing.T) ***REMOVED***
	tmpdir, err := ioutil.TempDir("", "TestFollowSymlinkRelativeLinkChangeScope")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmpdir)
	if err := makeFs(tmpdir, []dirOrLink***REMOVED******REMOVED***path: "testdata/fs/a/e", target: "../b"***REMOVED******REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := testSymlink(tmpdir, "testdata/fs/a/e/c/data", "testdata/fs/b/c/data", "testdata"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	// avoid letting allowing symlink e lead us to ../b
	// normalize to the "testdata/fs/a"
	if err := testSymlink(tmpdir, "testdata/fs/a/e", "testdata/fs/a/b", "testdata/fs/a"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestFollowSymlinkDeepRelativeLinkChangeScope(t *testing.T) ***REMOVED***
	tmpdir, err := ioutil.TempDir("", "TestFollowSymlinkDeepRelativeLinkChangeScope")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmpdir)

	if err := makeFs(tmpdir, []dirOrLink***REMOVED******REMOVED***path: "testdata/fs/a/f", target: "../../../../test"***REMOVED******REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	// avoid letting symlink f lead us out of the "testdata" scope
	// we don't normalize because symlink f is in scope and there is no
	// information leak
	if err := testSymlink(tmpdir, "testdata/fs/a/f", "testdata/test", "testdata"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	// avoid letting symlink f lead us out of the "testdata/fs" scope
	// we don't normalize because symlink f is in scope and there is no
	// information leak
	if err := testSymlink(tmpdir, "testdata/fs/a/f", "testdata/fs/test", "testdata/fs"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestFollowSymlinkRelativeLinkChain(t *testing.T) ***REMOVED***
	tmpdir, err := ioutil.TempDir("", "TestFollowSymlinkRelativeLinkChain")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmpdir)

	// avoid letting symlink g (pointed at by symlink h) take out of scope
	// TODO: we should probably normalize to scope here because ../[....]/root
	// is out of scope and we leak information
	if err := makeFs(tmpdir, []dirOrLink***REMOVED***
		***REMOVED***path: "testdata/fs/b/h", target: "../g"***REMOVED***,
		***REMOVED***path: "testdata/fs/g", target: "../../../../../../../../../../../../root"***REMOVED***,
	***REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := testSymlink(tmpdir, "testdata/fs/b/h", "testdata/root", "testdata"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestFollowSymlinkBreakoutPath(t *testing.T) ***REMOVED***
	tmpdir, err := ioutil.TempDir("", "TestFollowSymlinkBreakoutPath")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmpdir)

	// avoid letting symlink -> ../directory/file escape from scope
	// normalize to "testdata/fs/j"
	if err := makeFs(tmpdir, []dirOrLink***REMOVED******REMOVED***path: "testdata/fs/j/k", target: "../i/a"***REMOVED******REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := testSymlink(tmpdir, "testdata/fs/j/k", "testdata/fs/j/i/a", "testdata/fs/j"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestFollowSymlinkToRoot(t *testing.T) ***REMOVED***
	tmpdir, err := ioutil.TempDir("", "TestFollowSymlinkToRoot")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmpdir)

	// make sure we don't allow escaping to /
	// normalize to dir
	if err := makeFs(tmpdir, []dirOrLink***REMOVED******REMOVED***path: "foo", target: "/"***REMOVED******REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := testSymlink(tmpdir, "foo", "", ""); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestFollowSymlinkSlashDotdot(t *testing.T) ***REMOVED***
	tmpdir, err := ioutil.TempDir("", "TestFollowSymlinkSlashDotdot")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmpdir)
	tmpdir = filepath.Join(tmpdir, "dir", "subdir")

	// make sure we don't allow escaping to /
	// normalize to dir
	if err := makeFs(tmpdir, []dirOrLink***REMOVED******REMOVED***path: "foo", target: "/../../"***REMOVED******REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := testSymlink(tmpdir, "foo", "", ""); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestFollowSymlinkDotdot(t *testing.T) ***REMOVED***
	tmpdir, err := ioutil.TempDir("", "TestFollowSymlinkDotdot")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmpdir)
	tmpdir = filepath.Join(tmpdir, "dir", "subdir")

	// make sure we stay in scope without leaking information
	// this also checks for escaping to /
	// normalize to dir
	if err := makeFs(tmpdir, []dirOrLink***REMOVED******REMOVED***path: "foo", target: "../../"***REMOVED******REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := testSymlink(tmpdir, "foo", "", ""); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestFollowSymlinkRelativePath2(t *testing.T) ***REMOVED***
	tmpdir, err := ioutil.TempDir("", "TestFollowSymlinkRelativePath2")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmpdir)

	if err := makeFs(tmpdir, []dirOrLink***REMOVED******REMOVED***path: "bar/foo", target: "baz/target"***REMOVED******REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := testSymlink(tmpdir, "bar/foo", "bar/baz/target", ""); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestFollowSymlinkScopeLink(t *testing.T) ***REMOVED***
	tmpdir, err := ioutil.TempDir("", "TestFollowSymlinkScopeLink")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmpdir)

	if err := makeFs(tmpdir, []dirOrLink***REMOVED***
		***REMOVED***path: "root2"***REMOVED***,
		***REMOVED***path: "root", target: "root2"***REMOVED***,
		***REMOVED***path: "root2/foo", target: "../bar"***REMOVED***,
	***REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := testSymlink(tmpdir, "root/foo", "root/bar", "root"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestFollowSymlinkRootScope(t *testing.T) ***REMOVED***
	tmpdir, err := ioutil.TempDir("", "TestFollowSymlinkRootScope")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmpdir)

	expected, err := filepath.EvalSymlinks(tmpdir)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	rewrite, err := FollowSymlinkInScope(tmpdir, "/")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if rewrite != expected ***REMOVED***
		t.Fatalf("expected %q got %q", expected, rewrite)
	***REMOVED***
***REMOVED***

func TestFollowSymlinkEmpty(t *testing.T) ***REMOVED***
	res, err := FollowSymlinkInScope("", "")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	wd, err := os.Getwd()
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if res != wd ***REMOVED***
		t.Fatalf("expected %q got %q", wd, res)
	***REMOVED***
***REMOVED***

func TestFollowSymlinkCircular(t *testing.T) ***REMOVED***
	tmpdir, err := ioutil.TempDir("", "TestFollowSymlinkCircular")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmpdir)

	if err := makeFs(tmpdir, []dirOrLink***REMOVED******REMOVED***path: "root/foo", target: "foo"***REMOVED******REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := testSymlink(tmpdir, "root/foo", "", "root"); err == nil ***REMOVED***
		t.Fatal("expected an error for foo -> foo")
	***REMOVED***

	if err := makeFs(tmpdir, []dirOrLink***REMOVED***
		***REMOVED***path: "root/bar", target: "baz"***REMOVED***,
		***REMOVED***path: "root/baz", target: "../bak"***REMOVED***,
		***REMOVED***path: "root/bak", target: "/bar"***REMOVED***,
	***REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := testSymlink(tmpdir, "root/foo", "", "root"); err == nil ***REMOVED***
		t.Fatal("expected an error for bar -> baz -> bak -> bar")
	***REMOVED***
***REMOVED***

func TestFollowSymlinkComplexChainWithTargetPathsContainingLinks(t *testing.T) ***REMOVED***
	tmpdir, err := ioutil.TempDir("", "TestFollowSymlinkComplexChainWithTargetPathsContainingLinks")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmpdir)

	if err := makeFs(tmpdir, []dirOrLink***REMOVED***
		***REMOVED***path: "root2"***REMOVED***,
		***REMOVED***path: "root", target: "root2"***REMOVED***,
		***REMOVED***path: "root/a", target: "r/s"***REMOVED***,
		***REMOVED***path: "root/r", target: "../root/t"***REMOVED***,
		***REMOVED***path: "root/root/t/s/b", target: "/../u"***REMOVED***,
		***REMOVED***path: "root/u/c", target: "."***REMOVED***,
		***REMOVED***path: "root/u/x/y", target: "../v"***REMOVED***,
		***REMOVED***path: "root/u/v", target: "/../w"***REMOVED***,
	***REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := testSymlink(tmpdir, "root/a/b/c/x/y/z", "root/w/z", "root"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestFollowSymlinkBreakoutNonExistent(t *testing.T) ***REMOVED***
	tmpdir, err := ioutil.TempDir("", "TestFollowSymlinkBreakoutNonExistent")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmpdir)

	if err := makeFs(tmpdir, []dirOrLink***REMOVED***
		***REMOVED***path: "root/slash", target: "/"***REMOVED***,
		***REMOVED***path: "root/sym", target: "/idontexist/../slash"***REMOVED***,
	***REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := testSymlink(tmpdir, "root/sym/file", "root/file", "root"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***

func TestFollowSymlinkNoLexicalCleaning(t *testing.T) ***REMOVED***
	tmpdir, err := ioutil.TempDir("", "TestFollowSymlinkNoLexicalCleaning")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(tmpdir)

	if err := makeFs(tmpdir, []dirOrLink***REMOVED***
		***REMOVED***path: "root/sym", target: "/foo/bar"***REMOVED***,
		***REMOVED***path: "root/hello", target: "/sym/../baz"***REMOVED***,
	***REMOVED***); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if err := testSymlink(tmpdir, "root/hello", "root/foo/baz", "root"); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
***REMOVED***
