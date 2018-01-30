// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webdav

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"testing"

	"golang.org/x/net/context"
)

func TestSlashClean(t *testing.T) ***REMOVED***
	testCases := []string***REMOVED***
		"",
		".",
		"/",
		"/./",
		"//",
		"//.",
		"//a",
		"/a",
		"/a/b/c",
		"/a//b/./../c/d/",
		"a",
		"a/b/c",
	***REMOVED***
	for _, tc := range testCases ***REMOVED***
		got := slashClean(tc)
		want := path.Clean("/" + tc)
		if got != want ***REMOVED***
			t.Errorf("tc=%q: got %q, want %q", tc, got, want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestDirResolve(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		dir, name, want string
	***REMOVED******REMOVED***
		***REMOVED***"/", "", "/"***REMOVED***,
		***REMOVED***"/", "/", "/"***REMOVED***,
		***REMOVED***"/", ".", "/"***REMOVED***,
		***REMOVED***"/", "./a", "/a"***REMOVED***,
		***REMOVED***"/", "..", "/"***REMOVED***,
		***REMOVED***"/", "..", "/"***REMOVED***,
		***REMOVED***"/", "../", "/"***REMOVED***,
		***REMOVED***"/", "../.", "/"***REMOVED***,
		***REMOVED***"/", "../a", "/a"***REMOVED***,
		***REMOVED***"/", "../..", "/"***REMOVED***,
		***REMOVED***"/", "../bar/a", "/bar/a"***REMOVED***,
		***REMOVED***"/", "../baz/a", "/baz/a"***REMOVED***,
		***REMOVED***"/", "...", "/..."***REMOVED***,
		***REMOVED***"/", ".../a", "/.../a"***REMOVED***,
		***REMOVED***"/", ".../..", "/"***REMOVED***,
		***REMOVED***"/", "a", "/a"***REMOVED***,
		***REMOVED***"/", "a/./b", "/a/b"***REMOVED***,
		***REMOVED***"/", "a/../../b", "/b"***REMOVED***,
		***REMOVED***"/", "a/../b", "/b"***REMOVED***,
		***REMOVED***"/", "a/b", "/a/b"***REMOVED***,
		***REMOVED***"/", "a/b/c/../../d", "/a/d"***REMOVED***,
		***REMOVED***"/", "a/b/c/../../../d", "/d"***REMOVED***,
		***REMOVED***"/", "a/b/c/../../../../d", "/d"***REMOVED***,
		***REMOVED***"/", "a/b/c/d", "/a/b/c/d"***REMOVED***,

		***REMOVED***"/foo/bar", "", "/foo/bar"***REMOVED***,
		***REMOVED***"/foo/bar", "/", "/foo/bar"***REMOVED***,
		***REMOVED***"/foo/bar", ".", "/foo/bar"***REMOVED***,
		***REMOVED***"/foo/bar", "./a", "/foo/bar/a"***REMOVED***,
		***REMOVED***"/foo/bar", "..", "/foo/bar"***REMOVED***,
		***REMOVED***"/foo/bar", "../", "/foo/bar"***REMOVED***,
		***REMOVED***"/foo/bar", "../.", "/foo/bar"***REMOVED***,
		***REMOVED***"/foo/bar", "../a", "/foo/bar/a"***REMOVED***,
		***REMOVED***"/foo/bar", "../..", "/foo/bar"***REMOVED***,
		***REMOVED***"/foo/bar", "../bar/a", "/foo/bar/bar/a"***REMOVED***,
		***REMOVED***"/foo/bar", "../baz/a", "/foo/bar/baz/a"***REMOVED***,
		***REMOVED***"/foo/bar", "...", "/foo/bar/..."***REMOVED***,
		***REMOVED***"/foo/bar", ".../a", "/foo/bar/.../a"***REMOVED***,
		***REMOVED***"/foo/bar", ".../..", "/foo/bar"***REMOVED***,
		***REMOVED***"/foo/bar", "a", "/foo/bar/a"***REMOVED***,
		***REMOVED***"/foo/bar", "a/./b", "/foo/bar/a/b"***REMOVED***,
		***REMOVED***"/foo/bar", "a/../../b", "/foo/bar/b"***REMOVED***,
		***REMOVED***"/foo/bar", "a/../b", "/foo/bar/b"***REMOVED***,
		***REMOVED***"/foo/bar", "a/b", "/foo/bar/a/b"***REMOVED***,
		***REMOVED***"/foo/bar", "a/b/c/../../d", "/foo/bar/a/d"***REMOVED***,
		***REMOVED***"/foo/bar", "a/b/c/../../../d", "/foo/bar/d"***REMOVED***,
		***REMOVED***"/foo/bar", "a/b/c/../../../../d", "/foo/bar/d"***REMOVED***,
		***REMOVED***"/foo/bar", "a/b/c/d", "/foo/bar/a/b/c/d"***REMOVED***,

		***REMOVED***"/foo/bar/", "", "/foo/bar"***REMOVED***,
		***REMOVED***"/foo/bar/", "/", "/foo/bar"***REMOVED***,
		***REMOVED***"/foo/bar/", ".", "/foo/bar"***REMOVED***,
		***REMOVED***"/foo/bar/", "./a", "/foo/bar/a"***REMOVED***,
		***REMOVED***"/foo/bar/", "..", "/foo/bar"***REMOVED***,

		***REMOVED***"/foo//bar///", "", "/foo/bar"***REMOVED***,
		***REMOVED***"/foo//bar///", "/", "/foo/bar"***REMOVED***,
		***REMOVED***"/foo//bar///", ".", "/foo/bar"***REMOVED***,
		***REMOVED***"/foo//bar///", "./a", "/foo/bar/a"***REMOVED***,
		***REMOVED***"/foo//bar///", "..", "/foo/bar"***REMOVED***,

		***REMOVED***"/x/y/z", "ab/c\x00d/ef", ""***REMOVED***,

		***REMOVED***".", "", "."***REMOVED***,
		***REMOVED***".", "/", "."***REMOVED***,
		***REMOVED***".", ".", "."***REMOVED***,
		***REMOVED***".", "./a", "a"***REMOVED***,
		***REMOVED***".", "..", "."***REMOVED***,
		***REMOVED***".", "..", "."***REMOVED***,
		***REMOVED***".", "../", "."***REMOVED***,
		***REMOVED***".", "../.", "."***REMOVED***,
		***REMOVED***".", "../a", "a"***REMOVED***,
		***REMOVED***".", "../..", "."***REMOVED***,
		***REMOVED***".", "../bar/a", "bar/a"***REMOVED***,
		***REMOVED***".", "../baz/a", "baz/a"***REMOVED***,
		***REMOVED***".", "...", "..."***REMOVED***,
		***REMOVED***".", ".../a", ".../a"***REMOVED***,
		***REMOVED***".", ".../..", "."***REMOVED***,
		***REMOVED***".", "a", "a"***REMOVED***,
		***REMOVED***".", "a/./b", "a/b"***REMOVED***,
		***REMOVED***".", "a/../../b", "b"***REMOVED***,
		***REMOVED***".", "a/../b", "b"***REMOVED***,
		***REMOVED***".", "a/b", "a/b"***REMOVED***,
		***REMOVED***".", "a/b/c/../../d", "a/d"***REMOVED***,
		***REMOVED***".", "a/b/c/../../../d", "d"***REMOVED***,
		***REMOVED***".", "a/b/c/../../../../d", "d"***REMOVED***,
		***REMOVED***".", "a/b/c/d", "a/b/c/d"***REMOVED***,

		***REMOVED***"", "", "."***REMOVED***,
		***REMOVED***"", "/", "."***REMOVED***,
		***REMOVED***"", ".", "."***REMOVED***,
		***REMOVED***"", "./a", "a"***REMOVED***,
		***REMOVED***"", "..", "."***REMOVED***,
	***REMOVED***

	for _, tc := range testCases ***REMOVED***
		d := Dir(filepath.FromSlash(tc.dir))
		if got := filepath.ToSlash(d.resolve(tc.name)); got != tc.want ***REMOVED***
			t.Errorf("dir=%q, name=%q: got %q, want %q", tc.dir, tc.name, got, tc.want)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestWalk(t *testing.T) ***REMOVED***
	type walkStep struct ***REMOVED***
		name, frag string
		final      bool
	***REMOVED***

	testCases := []struct ***REMOVED***
		dir  string
		want []walkStep
	***REMOVED******REMOVED***
		***REMOVED***"", []walkStep***REMOVED***
			***REMOVED***"", "", true***REMOVED***,
		***REMOVED******REMOVED***,
		***REMOVED***"/", []walkStep***REMOVED***
			***REMOVED***"", "", true***REMOVED***,
		***REMOVED******REMOVED***,
		***REMOVED***"/a", []walkStep***REMOVED***
			***REMOVED***"", "a", true***REMOVED***,
		***REMOVED******REMOVED***,
		***REMOVED***"/a/", []walkStep***REMOVED***
			***REMOVED***"", "a", true***REMOVED***,
		***REMOVED******REMOVED***,
		***REMOVED***"/a/b", []walkStep***REMOVED***
			***REMOVED***"", "a", false***REMOVED***,
			***REMOVED***"a", "b", true***REMOVED***,
		***REMOVED******REMOVED***,
		***REMOVED***"/a/b/", []walkStep***REMOVED***
			***REMOVED***"", "a", false***REMOVED***,
			***REMOVED***"a", "b", true***REMOVED***,
		***REMOVED******REMOVED***,
		***REMOVED***"/a/b/c", []walkStep***REMOVED***
			***REMOVED***"", "a", false***REMOVED***,
			***REMOVED***"a", "b", false***REMOVED***,
			***REMOVED***"b", "c", true***REMOVED***,
		***REMOVED******REMOVED***,
		// The following test case is the one mentioned explicitly
		// in the method description.
		***REMOVED***"/foo/bar/x", []walkStep***REMOVED***
			***REMOVED***"", "foo", false***REMOVED***,
			***REMOVED***"foo", "bar", false***REMOVED***,
			***REMOVED***"bar", "x", true***REMOVED***,
		***REMOVED******REMOVED***,
	***REMOVED***

	ctx := context.Background()

	for _, tc := range testCases ***REMOVED***
		fs := NewMemFS().(*memFS)

		parts := strings.Split(tc.dir, "/")
		for p := 2; p < len(parts); p++ ***REMOVED***
			d := strings.Join(parts[:p], "/")
			if err := fs.Mkdir(ctx, d, 0666); err != nil ***REMOVED***
				t.Errorf("tc.dir=%q: mkdir: %q: %v", tc.dir, d, err)
			***REMOVED***
		***REMOVED***

		i, prevFrag := 0, ""
		err := fs.walk("test", tc.dir, func(dir *memFSNode, frag string, final bool) error ***REMOVED***
			got := walkStep***REMOVED***
				name:  prevFrag,
				frag:  frag,
				final: final,
			***REMOVED***
			want := tc.want[i]

			if got != want ***REMOVED***
				return fmt.Errorf("got %+v, want %+v", got, want)
			***REMOVED***
			i, prevFrag = i+1, frag
			return nil
		***REMOVED***)
		if err != nil ***REMOVED***
			t.Errorf("tc.dir=%q: %v", tc.dir, err)
		***REMOVED***
	***REMOVED***
***REMOVED***

// find appends to ss the names of the named file and its children. It is
// analogous to the Unix find command.
//
// The returned strings are not guaranteed to be in any particular order.
func find(ctx context.Context, ss []string, fs FileSystem, name string) ([]string, error) ***REMOVED***
	stat, err := fs.Stat(ctx, name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	ss = append(ss, name)
	if stat.IsDir() ***REMOVED***
		f, err := fs.OpenFile(ctx, name, os.O_RDONLY, 0)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		defer f.Close()
		children, err := f.Readdir(-1)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		for _, c := range children ***REMOVED***
			ss, err = find(ctx, ss, fs, path.Join(name, c.Name()))
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return ss, nil
***REMOVED***

func testFS(t *testing.T, fs FileSystem) ***REMOVED***
	errStr := func(err error) string ***REMOVED***
		switch ***REMOVED***
		case os.IsExist(err):
			return "errExist"
		case os.IsNotExist(err):
			return "errNotExist"
		case err != nil:
			return "err"
		***REMOVED***
		return "ok"
	***REMOVED***

	// The non-"find" non-"stat" test cases should change the file system state. The
	// indentation of the "find"s and "stat"s helps distinguish such test cases.
	testCases := []string***REMOVED***
		"  stat / want dir",
		"  stat /a want errNotExist",
		"  stat /d want errNotExist",
		"  stat /d/e want errNotExist",
		"create /a A want ok",
		"  stat /a want 1",
		"create /d/e EEE want errNotExist",
		"mk-dir /a want errExist",
		"mk-dir /d/m want errNotExist",
		"mk-dir /d want ok",
		"  stat /d want dir",
		"create /d/e EEE want ok",
		"  stat /d/e want 3",
		"  find / /a /d /d/e",
		"create /d/f FFFF want ok",
		"create /d/g GGGGGGG want ok",
		"mk-dir /d/m want ok",
		"mk-dir /d/m want errExist",
		"create /d/m/p PPPPP want ok",
		"  stat /d/e want 3",
		"  stat /d/f want 4",
		"  stat /d/g want 7",
		"  stat /d/h want errNotExist",
		"  stat /d/m want dir",
		"  stat /d/m/p want 5",
		"  find / /a /d /d/e /d/f /d/g /d/m /d/m/p",
		"rm-all /d want ok",
		"  stat /a want 1",
		"  stat /d want errNotExist",
		"  stat /d/e want errNotExist",
		"  stat /d/f want errNotExist",
		"  stat /d/g want errNotExist",
		"  stat /d/m want errNotExist",
		"  stat /d/m/p want errNotExist",
		"  find / /a",
		"mk-dir /d/m want errNotExist",
		"mk-dir /d want ok",
		"create /d/f FFFF want ok",
		"rm-all /d/f want ok",
		"mk-dir /d/m want ok",
		"rm-all /z want ok",
		"rm-all / want err",
		"create /b BB want ok",
		"  stat / want dir",
		"  stat /a want 1",
		"  stat /b want 2",
		"  stat /c want errNotExist",
		"  stat /d want dir",
		"  stat /d/m want dir",
		"  find / /a /b /d /d/m",
		"move__ o=F /b /c want ok",
		"  stat /b want errNotExist",
		"  stat /c want 2",
		"  stat /d/m want dir",
		"  stat /d/n want errNotExist",
		"  find / /a /c /d /d/m",
		"move__ o=F /d/m /d/n want ok",
		"create /d/n/q QQQQ want ok",
		"  stat /d/m want errNotExist",
		"  stat /d/n want dir",
		"  stat /d/n/q want 4",
		"move__ o=F /d /d/n/z want err",
		"move__ o=T /c /d/n/q want ok",
		"  stat /c want errNotExist",
		"  stat /d/n/q want 2",
		"  find / /a /d /d/n /d/n/q",
		"create /d/n/r RRRRR want ok",
		"mk-dir /u want ok",
		"mk-dir /u/v want ok",
		"move__ o=F /d/n /u want errExist",
		"create /t TTTTTT want ok",
		"move__ o=F /d/n /t want errExist",
		"rm-all /t want ok",
		"move__ o=F /d/n /t want ok",
		"  stat /d want dir",
		"  stat /d/n want errNotExist",
		"  stat /d/n/r want errNotExist",
		"  stat /t want dir",
		"  stat /t/q want 2",
		"  stat /t/r want 5",
		"  find / /a /d /t /t/q /t/r /u /u/v",
		"move__ o=F /t / want errExist",
		"move__ o=T /t /u/v want ok",
		"  stat /u/v/r want 5",
		"move__ o=F / /z want err",
		"  find / /a /d /u /u/v /u/v/q /u/v/r",
		"  stat /a want 1",
		"  stat /b want errNotExist",
		"  stat /c want errNotExist",
		"  stat /u/v/r want 5",
		"copy__ o=F d=0 /a /b want ok",
		"copy__ o=T d=0 /a /c want ok",
		"  stat /a want 1",
		"  stat /b want 1",
		"  stat /c want 1",
		"  stat /u/v/r want 5",
		"copy__ o=F d=0 /u/v/r /b want errExist",
		"  stat /b want 1",
		"copy__ o=T d=0 /u/v/r /b want ok",
		"  stat /a want 1",
		"  stat /b want 5",
		"  stat /u/v/r want 5",
		"rm-all /a want ok",
		"rm-all /b want ok",
		"mk-dir /u/v/w want ok",
		"create /u/v/w/s SSSSSSSS want ok",
		"  stat /d want dir",
		"  stat /d/x want errNotExist",
		"  stat /d/y want errNotExist",
		"  stat /u/v/r want 5",
		"  stat /u/v/w/s want 8",
		"  find / /c /d /u /u/v /u/v/q /u/v/r /u/v/w /u/v/w/s",
		"copy__ o=T d=0 /u/v /d/x want ok",
		"copy__ o=T d=∞ /u/v /d/y want ok",
		"rm-all /u want ok",
		"  stat /d/x want dir",
		"  stat /d/x/q want errNotExist",
		"  stat /d/x/r want errNotExist",
		"  stat /d/x/w want errNotExist",
		"  stat /d/x/w/s want errNotExist",
		"  stat /d/y want dir",
		"  stat /d/y/q want 2",
		"  stat /d/y/r want 5",
		"  stat /d/y/w want dir",
		"  stat /d/y/w/s want 8",
		"  stat /u want errNotExist",
		"  find / /c /d /d/x /d/y /d/y/q /d/y/r /d/y/w /d/y/w/s",
		"copy__ o=F d=∞ /d/y /d/x want errExist",
	***REMOVED***

	ctx := context.Background()

	for i, tc := range testCases ***REMOVED***
		tc = strings.TrimSpace(tc)
		j := strings.IndexByte(tc, ' ')
		if j < 0 ***REMOVED***
			t.Fatalf("test case #%d %q: invalid command", i, tc)
		***REMOVED***
		op, arg := tc[:j], tc[j+1:]

		switch op ***REMOVED***
		default:
			t.Fatalf("test case #%d %q: invalid operation %q", i, tc, op)

		case "create":
			parts := strings.Split(arg, " ")
			if len(parts) != 4 || parts[2] != "want" ***REMOVED***
				t.Fatalf("test case #%d %q: invalid write", i, tc)
			***REMOVED***
			f, opErr := fs.OpenFile(ctx, parts[0], os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
			if got := errStr(opErr); got != parts[3] ***REMOVED***
				t.Fatalf("test case #%d %q: OpenFile: got %q (%v), want %q", i, tc, got, opErr, parts[3])
			***REMOVED***
			if f != nil ***REMOVED***
				if _, err := f.Write([]byte(parts[1])); err != nil ***REMOVED***
					t.Fatalf("test case #%d %q: Write: %v", i, tc, err)
				***REMOVED***
				if err := f.Close(); err != nil ***REMOVED***
					t.Fatalf("test case #%d %q: Close: %v", i, tc, err)
				***REMOVED***
			***REMOVED***

		case "find":
			got, err := find(ctx, nil, fs, "/")
			if err != nil ***REMOVED***
				t.Fatalf("test case #%d %q: find: %v", i, tc, err)
			***REMOVED***
			sort.Strings(got)
			want := strings.Split(arg, " ")
			if !reflect.DeepEqual(got, want) ***REMOVED***
				t.Fatalf("test case #%d %q:\ngot  %s\nwant %s", i, tc, got, want)
			***REMOVED***

		case "copy__", "mk-dir", "move__", "rm-all", "stat":
			nParts := 3
			switch op ***REMOVED***
			case "copy__":
				nParts = 6
			case "move__":
				nParts = 5
			***REMOVED***
			parts := strings.Split(arg, " ")
			if len(parts) != nParts ***REMOVED***
				t.Fatalf("test case #%d %q: invalid %s", i, tc, op)
			***REMOVED***

			got, opErr := "", error(nil)
			switch op ***REMOVED***
			case "copy__":
				depth := 0
				if parts[1] == "d=∞" ***REMOVED***
					depth = infiniteDepth
				***REMOVED***
				_, opErr = copyFiles(ctx, fs, parts[2], parts[3], parts[0] == "o=T", depth, 0)
			case "mk-dir":
				opErr = fs.Mkdir(ctx, parts[0], 0777)
			case "move__":
				_, opErr = moveFiles(ctx, fs, parts[1], parts[2], parts[0] == "o=T")
			case "rm-all":
				opErr = fs.RemoveAll(ctx, parts[0])
			case "stat":
				var stat os.FileInfo
				fileName := parts[0]
				if stat, opErr = fs.Stat(ctx, fileName); opErr == nil ***REMOVED***
					if stat.IsDir() ***REMOVED***
						got = "dir"
					***REMOVED*** else ***REMOVED***
						got = strconv.Itoa(int(stat.Size()))
					***REMOVED***

					if fileName == "/" ***REMOVED***
						// For a Dir FileSystem, the virtual file system root maps to a
						// real file system name like "/tmp/webdav-test012345", which does
						// not end with "/". We skip such cases.
					***REMOVED*** else if statName := stat.Name(); path.Base(fileName) != statName ***REMOVED***
						t.Fatalf("test case #%d %q: file name %q inconsistent with stat name %q",
							i, tc, fileName, statName)
					***REMOVED***
				***REMOVED***
			***REMOVED***
			if got == "" ***REMOVED***
				got = errStr(opErr)
			***REMOVED***

			if parts[len(parts)-2] != "want" ***REMOVED***
				t.Fatalf("test case #%d %q: invalid %s", i, tc, op)
			***REMOVED***
			if want := parts[len(parts)-1]; got != want ***REMOVED***
				t.Fatalf("test case #%d %q: got %q (%v), want %q", i, tc, got, opErr, want)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestDir(t *testing.T) ***REMOVED***
	switch runtime.GOOS ***REMOVED***
	case "nacl":
		t.Skip("see golang.org/issue/12004")
	case "plan9":
		t.Skip("see golang.org/issue/11453")
	***REMOVED***

	td, err := ioutil.TempDir("", "webdav-test")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer os.RemoveAll(td)
	testFS(t, Dir(td))
***REMOVED***

func TestMemFS(t *testing.T) ***REMOVED***
	testFS(t, NewMemFS())
***REMOVED***

func TestMemFSRoot(t *testing.T) ***REMOVED***
	ctx := context.Background()
	fs := NewMemFS()
	for i := 0; i < 5; i++ ***REMOVED***
		stat, err := fs.Stat(ctx, "/")
		if err != nil ***REMOVED***
			t.Fatalf("i=%d: Stat: %v", i, err)
		***REMOVED***
		if !stat.IsDir() ***REMOVED***
			t.Fatalf("i=%d: Stat.IsDir is false, want true", i)
		***REMOVED***

		f, err := fs.OpenFile(ctx, "/", os.O_RDONLY, 0)
		if err != nil ***REMOVED***
			t.Fatalf("i=%d: OpenFile: %v", i, err)
		***REMOVED***
		defer f.Close()
		children, err := f.Readdir(-1)
		if err != nil ***REMOVED***
			t.Fatalf("i=%d: Readdir: %v", i, err)
		***REMOVED***
		if len(children) != i ***REMOVED***
			t.Fatalf("i=%d: got %d children, want %d", i, len(children), i)
		***REMOVED***

		if _, err := f.Write(make([]byte, 1)); err == nil ***REMOVED***
			t.Fatalf("i=%d: Write: got nil error, want non-nil", i)
		***REMOVED***

		if err := fs.Mkdir(ctx, fmt.Sprintf("/dir%d", i), 0777); err != nil ***REMOVED***
			t.Fatalf("i=%d: Mkdir: %v", i, err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestMemFileReaddir(t *testing.T) ***REMOVED***
	ctx := context.Background()
	fs := NewMemFS()
	if err := fs.Mkdir(ctx, "/foo", 0777); err != nil ***REMOVED***
		t.Fatalf("Mkdir: %v", err)
	***REMOVED***
	readdir := func(count int) ([]os.FileInfo, error) ***REMOVED***
		f, err := fs.OpenFile(ctx, "/foo", os.O_RDONLY, 0)
		if err != nil ***REMOVED***
			t.Fatalf("OpenFile: %v", err)
		***REMOVED***
		defer f.Close()
		return f.Readdir(count)
	***REMOVED***
	if got, err := readdir(-1); len(got) != 0 || err != nil ***REMOVED***
		t.Fatalf("readdir(-1): got %d fileInfos with err=%v, want 0, <nil>", len(got), err)
	***REMOVED***
	if got, err := readdir(+1); len(got) != 0 || err != io.EOF ***REMOVED***
		t.Fatalf("readdir(+1): got %d fileInfos with err=%v, want 0, EOF", len(got), err)
	***REMOVED***
***REMOVED***

func TestMemFile(t *testing.T) ***REMOVED***
	testCases := []string***REMOVED***
		"wantData ",
		"wantSize 0",
		"write abc",
		"wantData abc",
		"write de",
		"wantData abcde",
		"wantSize 5",
		"write 5*x",
		"write 4*y+2*z",
		"write 3*st",
		"wantData abcdexxxxxyyyyzzststst",
		"wantSize 22",
		"seek set 4 want 4",
		"write EFG",
		"wantData abcdEFGxxxyyyyzzststst",
		"wantSize 22",
		"seek set 2 want 2",
		"read cdEF",
		"read Gx",
		"seek cur 0 want 8",
		"seek cur 2 want 10",
		"seek cur -1 want 9",
		"write J",
		"wantData abcdEFGxxJyyyyzzststst",
		"wantSize 22",
		"seek cur -4 want 6",
		"write ghijk",
		"wantData abcdEFghijkyyyzzststst",
		"wantSize 22",
		"read yyyz",
		"seek cur 0 want 15",
		"write ",
		"seek cur 0 want 15",
		"read ",
		"seek cur 0 want 15",
		"seek end -3 want 19",
		"write ZZ",
		"wantData abcdEFghijkyyyzzstsZZt",
		"wantSize 22",
		"write 4*A",
		"wantData abcdEFghijkyyyzzstsZZAAAA",
		"wantSize 25",
		"seek end 0 want 25",
		"seek end -5 want 20",
		"read Z+4*A",
		"write 5*B",
		"wantData abcdEFghijkyyyzzstsZZAAAABBBBB",
		"wantSize 30",
		"seek end 10 want 40",
		"write C",
		"wantData abcdEFghijkyyyzzstsZZAAAABBBBB..........C",
		"wantSize 41",
		"write D",
		"wantData abcdEFghijkyyyzzstsZZAAAABBBBB..........CD",
		"wantSize 42",
		"seek set 43 want 43",
		"write E",
		"wantData abcdEFghijkyyyzzstsZZAAAABBBBB..........CD.E",
		"wantSize 44",
		"seek set 0 want 0",
		"write 5*123456789_",
		"wantData 123456789_123456789_123456789_123456789_123456789_",
		"wantSize 50",
		"seek cur 0 want 50",
		"seek cur -99 want err",
	***REMOVED***

	ctx := context.Background()

	const filename = "/foo"
	fs := NewMemFS()
	f, err := fs.OpenFile(ctx, filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil ***REMOVED***
		t.Fatalf("OpenFile: %v", err)
	***REMOVED***
	defer f.Close()

	for i, tc := range testCases ***REMOVED***
		j := strings.IndexByte(tc, ' ')
		if j < 0 ***REMOVED***
			t.Fatalf("test case #%d %q: invalid command", i, tc)
		***REMOVED***
		op, arg := tc[:j], tc[j+1:]

		// Expand an arg like "3*a+2*b" to "aaabb".
		parts := strings.Split(arg, "+")
		for j, part := range parts ***REMOVED***
			if k := strings.IndexByte(part, '*'); k >= 0 ***REMOVED***
				repeatCount, repeatStr := part[:k], part[k+1:]
				n, err := strconv.Atoi(repeatCount)
				if err != nil ***REMOVED***
					t.Fatalf("test case #%d %q: invalid repeat count %q", i, tc, repeatCount)
				***REMOVED***
				parts[j] = strings.Repeat(repeatStr, n)
			***REMOVED***
		***REMOVED***
		arg = strings.Join(parts, "")

		switch op ***REMOVED***
		default:
			t.Fatalf("test case #%d %q: invalid operation %q", i, tc, op)

		case "read":
			buf := make([]byte, len(arg))
			if _, err := io.ReadFull(f, buf); err != nil ***REMOVED***
				t.Fatalf("test case #%d %q: ReadFull: %v", i, tc, err)
			***REMOVED***
			if got := string(buf); got != arg ***REMOVED***
				t.Fatalf("test case #%d %q:\ngot  %q\nwant %q", i, tc, got, arg)
			***REMOVED***

		case "seek":
			parts := strings.Split(arg, " ")
			if len(parts) != 4 ***REMOVED***
				t.Fatalf("test case #%d %q: invalid seek", i, tc)
			***REMOVED***

			whence := 0
			switch parts[0] ***REMOVED***
			default:
				t.Fatalf("test case #%d %q: invalid seek whence", i, tc)
			case "set":
				whence = os.SEEK_SET
			case "cur":
				whence = os.SEEK_CUR
			case "end":
				whence = os.SEEK_END
			***REMOVED***
			offset, err := strconv.Atoi(parts[1])
			if err != nil ***REMOVED***
				t.Fatalf("test case #%d %q: invalid offset %q", i, tc, parts[1])
			***REMOVED***

			if parts[2] != "want" ***REMOVED***
				t.Fatalf("test case #%d %q: invalid seek", i, tc)
			***REMOVED***
			if parts[3] == "err" ***REMOVED***
				_, err := f.Seek(int64(offset), whence)
				if err == nil ***REMOVED***
					t.Fatalf("test case #%d %q: Seek returned nil error, want non-nil", i, tc)
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				got, err := f.Seek(int64(offset), whence)
				if err != nil ***REMOVED***
					t.Fatalf("test case #%d %q: Seek: %v", i, tc, err)
				***REMOVED***
				want, err := strconv.Atoi(parts[3])
				if err != nil ***REMOVED***
					t.Fatalf("test case #%d %q: invalid want %q", i, tc, parts[3])
				***REMOVED***
				if got != int64(want) ***REMOVED***
					t.Fatalf("test case #%d %q: got %d, want %d", i, tc, got, want)
				***REMOVED***
			***REMOVED***

		case "write":
			n, err := f.Write([]byte(arg))
			if err != nil ***REMOVED***
				t.Fatalf("test case #%d %q: write: %v", i, tc, err)
			***REMOVED***
			if n != len(arg) ***REMOVED***
				t.Fatalf("test case #%d %q: write returned %d bytes, want %d", i, tc, n, len(arg))
			***REMOVED***

		case "wantData":
			g, err := fs.OpenFile(ctx, filename, os.O_RDONLY, 0666)
			if err != nil ***REMOVED***
				t.Fatalf("test case #%d %q: OpenFile: %v", i, tc, err)
			***REMOVED***
			gotBytes, err := ioutil.ReadAll(g)
			if err != nil ***REMOVED***
				t.Fatalf("test case #%d %q: ReadAll: %v", i, tc, err)
			***REMOVED***
			for i, c := range gotBytes ***REMOVED***
				if c == '\x00' ***REMOVED***
					gotBytes[i] = '.'
				***REMOVED***
			***REMOVED***
			got := string(gotBytes)
			if got != arg ***REMOVED***
				t.Fatalf("test case #%d %q:\ngot  %q\nwant %q", i, tc, got, arg)
			***REMOVED***
			if err := g.Close(); err != nil ***REMOVED***
				t.Fatalf("test case #%d %q: Close: %v", i, tc, err)
			***REMOVED***

		case "wantSize":
			n, err := strconv.Atoi(arg)
			if err != nil ***REMOVED***
				t.Fatalf("test case #%d %q: invalid size %q", i, tc, arg)
			***REMOVED***
			fi, err := fs.Stat(ctx, filename)
			if err != nil ***REMOVED***
				t.Fatalf("test case #%d %q: Stat: %v", i, tc, err)
			***REMOVED***
			if got, want := fi.Size(), int64(n); got != want ***REMOVED***
				t.Fatalf("test case #%d %q: got %d, want %d", i, tc, got, want)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// TestMemFileWriteAllocs tests that writing N consecutive 1KiB chunks to a
// memFile doesn't allocate a new buffer for each of those N times. Otherwise,
// calling io.Copy(aMemFile, src) is likely to have quadratic complexity.
func TestMemFileWriteAllocs(t *testing.T) ***REMOVED***
	if runtime.Compiler == "gccgo" ***REMOVED***
		t.Skip("gccgo allocates here")
	***REMOVED***
	ctx := context.Background()
	fs := NewMemFS()
	f, err := fs.OpenFile(ctx, "/xxx", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil ***REMOVED***
		t.Fatalf("OpenFile: %v", err)
	***REMOVED***
	defer f.Close()

	xxx := make([]byte, 1024)
	for i := range xxx ***REMOVED***
		xxx[i] = 'x'
	***REMOVED***

	a := testing.AllocsPerRun(100, func() ***REMOVED***
		f.Write(xxx)
	***REMOVED***)
	// AllocsPerRun returns an integral value, so we compare the rounded-down
	// number to zero.
	if a > 0 ***REMOVED***
		t.Fatalf("%v allocs per run, want 0", a)
	***REMOVED***
***REMOVED***

func BenchmarkMemFileWrite(b *testing.B) ***REMOVED***
	ctx := context.Background()
	fs := NewMemFS()
	xxx := make([]byte, 1024)
	for i := range xxx ***REMOVED***
		xxx[i] = 'x'
	***REMOVED***

	b.ResetTimer()
	for i := 0; i < b.N; i++ ***REMOVED***
		f, err := fs.OpenFile(ctx, "/xxx", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil ***REMOVED***
			b.Fatalf("OpenFile: %v", err)
		***REMOVED***
		for j := 0; j < 100; j++ ***REMOVED***
			f.Write(xxx)
		***REMOVED***
		if err := f.Close(); err != nil ***REMOVED***
			b.Fatalf("Close: %v", err)
		***REMOVED***
		if err := fs.RemoveAll(ctx, "/xxx"); err != nil ***REMOVED***
			b.Fatalf("RemoveAll: %v", err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestCopyMoveProps(t *testing.T) ***REMOVED***
	ctx := context.Background()
	fs := NewMemFS()
	create := func(name string) error ***REMOVED***
		f, err := fs.OpenFile(ctx, name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		_, wErr := f.Write([]byte("contents"))
		cErr := f.Close()
		if wErr != nil ***REMOVED***
			return wErr
		***REMOVED***
		return cErr
	***REMOVED***
	patch := func(name string, patches ...Proppatch) error ***REMOVED***
		f, err := fs.OpenFile(ctx, name, os.O_RDWR, 0666)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		_, pErr := f.(DeadPropsHolder).Patch(patches)
		cErr := f.Close()
		if pErr != nil ***REMOVED***
			return pErr
		***REMOVED***
		return cErr
	***REMOVED***
	props := func(name string) (map[xml.Name]Property, error) ***REMOVED***
		f, err := fs.OpenFile(ctx, name, os.O_RDWR, 0666)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		m, pErr := f.(DeadPropsHolder).DeadProps()
		cErr := f.Close()
		if pErr != nil ***REMOVED***
			return nil, pErr
		***REMOVED***
		if cErr != nil ***REMOVED***
			return nil, cErr
		***REMOVED***
		return m, nil
	***REMOVED***

	p0 := Property***REMOVED***
		XMLName:  xml.Name***REMOVED***Space: "x:", Local: "boat"***REMOVED***,
		InnerXML: []byte("pea-green"),
	***REMOVED***
	p1 := Property***REMOVED***
		XMLName:  xml.Name***REMOVED***Space: "x:", Local: "ring"***REMOVED***,
		InnerXML: []byte("1 shilling"),
	***REMOVED***
	p2 := Property***REMOVED***
		XMLName:  xml.Name***REMOVED***Space: "x:", Local: "spoon"***REMOVED***,
		InnerXML: []byte("runcible"),
	***REMOVED***
	p3 := Property***REMOVED***
		XMLName:  xml.Name***REMOVED***Space: "x:", Local: "moon"***REMOVED***,
		InnerXML: []byte("light"),
	***REMOVED***

	if err := create("/src"); err != nil ***REMOVED***
		t.Fatalf("create /src: %v", err)
	***REMOVED***
	if err := patch("/src", Proppatch***REMOVED***Props: []Property***REMOVED***p0, p1***REMOVED******REMOVED***); err != nil ***REMOVED***
		t.Fatalf("patch /src +p0 +p1: %v", err)
	***REMOVED***
	if _, err := copyFiles(ctx, fs, "/src", "/tmp", true, infiniteDepth, 0); err != nil ***REMOVED***
		t.Fatalf("copyFiles /src /tmp: %v", err)
	***REMOVED***
	if _, err := moveFiles(ctx, fs, "/tmp", "/dst", true); err != nil ***REMOVED***
		t.Fatalf("moveFiles /tmp /dst: %v", err)
	***REMOVED***
	if err := patch("/src", Proppatch***REMOVED***Props: []Property***REMOVED***p0***REMOVED***, Remove: true***REMOVED***); err != nil ***REMOVED***
		t.Fatalf("patch /src -p0: %v", err)
	***REMOVED***
	if err := patch("/src", Proppatch***REMOVED***Props: []Property***REMOVED***p2***REMOVED******REMOVED***); err != nil ***REMOVED***
		t.Fatalf("patch /src +p2: %v", err)
	***REMOVED***
	if err := patch("/dst", Proppatch***REMOVED***Props: []Property***REMOVED***p1***REMOVED***, Remove: true***REMOVED***); err != nil ***REMOVED***
		t.Fatalf("patch /dst -p1: %v", err)
	***REMOVED***
	if err := patch("/dst", Proppatch***REMOVED***Props: []Property***REMOVED***p3***REMOVED******REMOVED***); err != nil ***REMOVED***
		t.Fatalf("patch /dst +p3: %v", err)
	***REMOVED***

	gotSrc, err := props("/src")
	if err != nil ***REMOVED***
		t.Fatalf("props /src: %v", err)
	***REMOVED***
	wantSrc := map[xml.Name]Property***REMOVED***
		p1.XMLName: p1,
		p2.XMLName: p2,
	***REMOVED***
	if !reflect.DeepEqual(gotSrc, wantSrc) ***REMOVED***
		t.Fatalf("props /src:\ngot  %v\nwant %v", gotSrc, wantSrc)
	***REMOVED***

	gotDst, err := props("/dst")
	if err != nil ***REMOVED***
		t.Fatalf("props /dst: %v", err)
	***REMOVED***
	wantDst := map[xml.Name]Property***REMOVED***
		p0.XMLName: p0,
		p3.XMLName: p3,
	***REMOVED***
	if !reflect.DeepEqual(gotDst, wantDst) ***REMOVED***
		t.Fatalf("props /dst:\ngot  %v\nwant %v", gotDst, wantDst)
	***REMOVED***
***REMOVED***

func TestWalkFS(t *testing.T) ***REMOVED***
	testCases := []struct ***REMOVED***
		desc    string
		buildfs []string
		startAt string
		depth   int
		walkFn  filepath.WalkFunc
		want    []string
	***REMOVED******REMOVED******REMOVED***
		"just root",
		[]string***REMOVED******REMOVED***,
		"/",
		infiniteDepth,
		nil,
		[]string***REMOVED***
			"/",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		"infinite walk from root",
		[]string***REMOVED***
			"mkdir /a",
			"mkdir /a/b",
			"touch /a/b/c",
			"mkdir /a/d",
			"mkdir /e",
			"touch /f",
		***REMOVED***,
		"/",
		infiniteDepth,
		nil,
		[]string***REMOVED***
			"/",
			"/a",
			"/a/b",
			"/a/b/c",
			"/a/d",
			"/e",
			"/f",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		"infinite walk from subdir",
		[]string***REMOVED***
			"mkdir /a",
			"mkdir /a/b",
			"touch /a/b/c",
			"mkdir /a/d",
			"mkdir /e",
			"touch /f",
		***REMOVED***,
		"/a",
		infiniteDepth,
		nil,
		[]string***REMOVED***
			"/a",
			"/a/b",
			"/a/b/c",
			"/a/d",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		"depth 1 walk from root",
		[]string***REMOVED***
			"mkdir /a",
			"mkdir /a/b",
			"touch /a/b/c",
			"mkdir /a/d",
			"mkdir /e",
			"touch /f",
		***REMOVED***,
		"/",
		1,
		nil,
		[]string***REMOVED***
			"/",
			"/a",
			"/e",
			"/f",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		"depth 1 walk from subdir",
		[]string***REMOVED***
			"mkdir /a",
			"mkdir /a/b",
			"touch /a/b/c",
			"mkdir /a/b/g",
			"mkdir /a/b/g/h",
			"touch /a/b/g/i",
			"touch /a/b/g/h/j",
		***REMOVED***,
		"/a/b",
		1,
		nil,
		[]string***REMOVED***
			"/a/b",
			"/a/b/c",
			"/a/b/g",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		"depth 0 walk from subdir",
		[]string***REMOVED***
			"mkdir /a",
			"mkdir /a/b",
			"touch /a/b/c",
			"mkdir /a/b/g",
			"mkdir /a/b/g/h",
			"touch /a/b/g/i",
			"touch /a/b/g/h/j",
		***REMOVED***,
		"/a/b",
		0,
		nil,
		[]string***REMOVED***
			"/a/b",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		"infinite walk from file",
		[]string***REMOVED***
			"mkdir /a",
			"touch /a/b",
			"touch /a/c",
		***REMOVED***,
		"/a/b",
		0,
		nil,
		[]string***REMOVED***
			"/a/b",
		***REMOVED***,
	***REMOVED***, ***REMOVED***
		"infinite walk with skipped subdir",
		[]string***REMOVED***
			"mkdir /a",
			"mkdir /a/b",
			"touch /a/b/c",
			"mkdir /a/b/g",
			"mkdir /a/b/g/h",
			"touch /a/b/g/i",
			"touch /a/b/g/h/j",
			"touch /a/b/z",
		***REMOVED***,
		"/",
		infiniteDepth,
		func(path string, info os.FileInfo, err error) error ***REMOVED***
			if path == "/a/b/g" ***REMOVED***
				return filepath.SkipDir
			***REMOVED***
			return nil
		***REMOVED***,
		[]string***REMOVED***
			"/",
			"/a",
			"/a/b",
			"/a/b/c",
			"/a/b/z",
		***REMOVED***,
	***REMOVED******REMOVED***
	ctx := context.Background()
	for _, tc := range testCases ***REMOVED***
		fs, err := buildTestFS(tc.buildfs)
		if err != nil ***REMOVED***
			t.Fatalf("%s: cannot create test filesystem: %v", tc.desc, err)
		***REMOVED***
		var got []string
		traceFn := func(path string, info os.FileInfo, err error) error ***REMOVED***
			if tc.walkFn != nil ***REMOVED***
				err = tc.walkFn(path, info, err)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
			got = append(got, path)
			return nil
		***REMOVED***
		fi, err := fs.Stat(ctx, tc.startAt)
		if err != nil ***REMOVED***
			t.Fatalf("%s: cannot stat: %v", tc.desc, err)
		***REMOVED***
		err = walkFS(ctx, fs, tc.depth, tc.startAt, fi, traceFn)
		if err != nil ***REMOVED***
			t.Errorf("%s:\ngot error %v, want nil", tc.desc, err)
			continue
		***REMOVED***
		sort.Strings(got)
		sort.Strings(tc.want)
		if !reflect.DeepEqual(got, tc.want) ***REMOVED***
			t.Errorf("%s:\ngot  %q\nwant %q", tc.desc, got, tc.want)
			continue
		***REMOVED***
	***REMOVED***
***REMOVED***

func buildTestFS(buildfs []string) (FileSystem, error) ***REMOVED***
	// TODO: Could this be merged with the build logic in TestFS?

	ctx := context.Background()
	fs := NewMemFS()
	for _, b := range buildfs ***REMOVED***
		op := strings.Split(b, " ")
		switch op[0] ***REMOVED***
		case "mkdir":
			err := fs.Mkdir(ctx, op[1], os.ModeDir|0777)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		case "touch":
			f, err := fs.OpenFile(ctx, op[1], os.O_RDWR|os.O_CREATE, 0666)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			f.Close()
		case "write":
			f, err := fs.OpenFile(ctx, op[1], os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			_, err = f.Write([]byte(op[2]))
			f.Close()
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		default:
			return nil, fmt.Errorf("unknown file operation %q", op[0])
		***REMOVED***
	***REMOVED***
	return fs, nil
***REMOVED***
