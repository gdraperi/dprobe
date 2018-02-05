// Copyright © 2014 Steve Francia <spf@spf13.com>.
// Copyright 2009 The Go Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package afero

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"testing"
)

var testName = "test.txt"
var Fss = []Fs***REMOVED***&MemMapFs***REMOVED******REMOVED***, &OsFs***REMOVED******REMOVED******REMOVED***

var testRegistry map[Fs][]string = make(map[Fs][]string)

func testDir(fs Fs) string ***REMOVED***
	name, err := TempDir(fs, "", "afero")
	if err != nil ***REMOVED***
		panic(fmt.Sprint("unable to work with test dir", err))
	***REMOVED***
	testRegistry[fs] = append(testRegistry[fs], name)

	return name
***REMOVED***

func tmpFile(fs Fs) File ***REMOVED***
	x, err := TempFile(fs, "", "afero")

	if err != nil ***REMOVED***
		panic(fmt.Sprint("unable to work with temp file", err))
	***REMOVED***

	testRegistry[fs] = append(testRegistry[fs], x.Name())

	return x
***REMOVED***

//Read with length 0 should not return EOF.
func TestRead0(t *testing.T) ***REMOVED***
	for _, fs := range Fss ***REMOVED***
		f := tmpFile(fs)
		defer f.Close()
		f.WriteString("Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.")

		var b []byte
		// b := make([]byte, 0)
		n, err := f.Read(b)
		if n != 0 || err != nil ***REMOVED***
			t.Errorf("%v: Read(0) = %d, %v, want 0, nil", fs.Name(), n, err)
		***REMOVED***
		f.Seek(0, 0)
		b = make([]byte, 100)
		n, err = f.Read(b)
		if n <= 0 || err != nil ***REMOVED***
			t.Errorf("%v: Read(100) = %d, %v, want >0, nil", fs.Name(), n, err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestOpenFile(t *testing.T) ***REMOVED***
	defer removeAllTestFiles(t)
	for _, fs := range Fss ***REMOVED***
		tmp := testDir(fs)
		path := filepath.Join(tmp, testName)

		f, err := fs.OpenFile(path, os.O_RDWR|os.O_CREATE, 0600)
		if err != nil ***REMOVED***
			t.Error(fs.Name(), "OpenFile (O_CREATE) failed:", err)
			continue
		***REMOVED***
		io.WriteString(f, "initial")
		f.Close()

		f, err = fs.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil ***REMOVED***
			t.Error(fs.Name(), "OpenFile (O_APPEND) failed:", err)
			continue
		***REMOVED***
		io.WriteString(f, "|append")
		f.Close()

		f, err = fs.OpenFile(path, os.O_RDONLY, 0600)
		contents, _ := ioutil.ReadAll(f)
		expectedContents := "initial|append"
		if string(contents) != expectedContents ***REMOVED***
			t.Errorf("%v: appending, expected '%v', got: '%v'", fs.Name(), expectedContents, string(contents))
		***REMOVED***
		f.Close()

		f, err = fs.OpenFile(path, os.O_RDWR|os.O_TRUNC, 0600)
		if err != nil ***REMOVED***
			t.Error(fs.Name(), "OpenFile (O_TRUNC) failed:", err)
			continue
		***REMOVED***
		contents, _ = ioutil.ReadAll(f)
		if string(contents) != "" ***REMOVED***
			t.Errorf("%v: expected truncated file, got: '%v'", fs.Name(), string(contents))
		***REMOVED***
		f.Close()
	***REMOVED***
***REMOVED***

func TestCreate(t *testing.T) ***REMOVED***
	defer removeAllTestFiles(t)
	for _, fs := range Fss ***REMOVED***
		tmp := testDir(fs)
		path := filepath.Join(tmp, testName)

		f, err := fs.Create(path)
		if err != nil ***REMOVED***
			t.Error(fs.Name(), "Create failed:", err)
			f.Close()
			continue
		***REMOVED***
		io.WriteString(f, "initial")
		f.Close()

		f, err = fs.Create(path)
		if err != nil ***REMOVED***
			t.Error(fs.Name(), "Create failed:", err)
			f.Close()
			continue
		***REMOVED***
		secondContent := "second create"
		io.WriteString(f, secondContent)
		f.Close()

		f, err = fs.Open(path)
		if err != nil ***REMOVED***
			t.Error(fs.Name(), "Open failed:", err)
			f.Close()
			continue
		***REMOVED***
		buf, err := ReadAll(f)
		if err != nil ***REMOVED***
			t.Error(fs.Name(), "ReadAll failed:", err)
			f.Close()
			continue
		***REMOVED***
		if string(buf) != secondContent ***REMOVED***
			t.Error(fs.Name(), "Content should be", "\""+secondContent+"\" but is \""+string(buf)+"\"")
			f.Close()
			continue
		***REMOVED***
		f.Close()
	***REMOVED***
***REMOVED***

func TestMemFileRead(t *testing.T) ***REMOVED***
	f := tmpFile(new(MemMapFs))
	// f := MemFileCreate("testfile")
	f.WriteString("abcd")
	f.Seek(0, 0)
	b := make([]byte, 8)
	n, err := f.Read(b)
	if n != 4 ***REMOVED***
		t.Errorf("didn't read all bytes: %v %v %v", n, err, b)
	***REMOVED***
	if err != nil ***REMOVED***
		t.Errorf("err is not nil: %v %v %v", n, err, b)
	***REMOVED***
	n, err = f.Read(b)
	if n != 0 ***REMOVED***
		t.Errorf("read more bytes: %v %v %v", n, err, b)
	***REMOVED***
	if err != io.EOF ***REMOVED***
		t.Errorf("error is not EOF: %v %v %v", n, err, b)
	***REMOVED***
***REMOVED***

func TestRename(t *testing.T) ***REMOVED***
	defer removeAllTestFiles(t)
	for _, fs := range Fss ***REMOVED***
		tDir := testDir(fs)
		from := filepath.Join(tDir, "/renamefrom")
		to := filepath.Join(tDir, "/renameto")
		exists := filepath.Join(tDir, "/renameexists")
		file, err := fs.Create(from)
		if err != nil ***REMOVED***
			t.Fatalf("%s: open %q failed: %v", fs.Name(), to, err)
		***REMOVED***
		if err = file.Close(); err != nil ***REMOVED***
			t.Errorf("%s: close %q failed: %v", fs.Name(), to, err)
		***REMOVED***
		file, err = fs.Create(exists)
		if err != nil ***REMOVED***
			t.Fatalf("%s: open %q failed: %v", fs.Name(), to, err)
		***REMOVED***
		if err = file.Close(); err != nil ***REMOVED***
			t.Errorf("%s: close %q failed: %v", fs.Name(), to, err)
		***REMOVED***
		err = fs.Rename(from, to)
		if err != nil ***REMOVED***
			t.Fatalf("%s: rename %q, %q failed: %v", fs.Name(), to, from, err)
		***REMOVED***
		file, err = fs.Create(from)
		if err != nil ***REMOVED***
			t.Fatalf("%s: open %q failed: %v", fs.Name(), to, err)
		***REMOVED***
		if err = file.Close(); err != nil ***REMOVED***
			t.Errorf("%s: close %q failed: %v", fs.Name(), to, err)
		***REMOVED***
		err = fs.Rename(from, exists)
		if err != nil ***REMOVED***
			t.Errorf("%s: rename %q, %q failed: %v", fs.Name(), exists, from, err)
		***REMOVED***
		names, err := readDirNames(fs, tDir)
		if err != nil ***REMOVED***
			t.Errorf("%s: readDirNames error: %v", fs.Name(), err)
		***REMOVED***
		found := false
		for _, e := range names ***REMOVED***
			if e == "renamefrom" ***REMOVED***
				t.Error("File is still called renamefrom")
			***REMOVED***
			if e == "renameto" ***REMOVED***
				found = true
			***REMOVED***
		***REMOVED***
		if !found ***REMOVED***
			t.Error("File was not renamed to renameto")
		***REMOVED***

		_, err = fs.Stat(to)
		if err != nil ***REMOVED***
			t.Errorf("%s: stat %q failed: %v", fs.Name(), to, err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestRemove(t *testing.T) ***REMOVED***
	for _, fs := range Fss ***REMOVED***

		x, err := TempFile(fs, "", "afero")
		if err != nil ***REMOVED***
			t.Error(fmt.Sprint("unable to work with temp file", err))
		***REMOVED***

		path := x.Name()
		x.Close()

		tDir := filepath.Dir(path)

		err = fs.Remove(path)
		if err != nil ***REMOVED***
			t.Errorf("%v: Remove() failed: %v", fs.Name(), err)
			continue
		***REMOVED***

		_, err = fs.Stat(path)
		if !os.IsNotExist(err) ***REMOVED***
			t.Errorf("%v: Remove() didn't remove file", fs.Name())
			continue
		***REMOVED***

		// Deleting non-existent file should raise error
		err = fs.Remove(path)
		if !os.IsNotExist(err) ***REMOVED***
			t.Errorf("%v: Remove() didn't raise error for non-existent file", fs.Name())
		***REMOVED***

		f, err := fs.Open(tDir)
		if err != nil ***REMOVED***
			t.Error("TestDir should still exist:", err)
		***REMOVED***

		names, err := f.Readdirnames(-1)
		if err != nil ***REMOVED***
			t.Error("Readdirnames failed:", err)
		***REMOVED***

		for _, e := range names ***REMOVED***
			if e == testName ***REMOVED***
				t.Error("File was not removed from parent directory")
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestTruncate(t *testing.T) ***REMOVED***
	defer removeAllTestFiles(t)
	for _, fs := range Fss ***REMOVED***
		f := tmpFile(fs)
		defer f.Close()

		checkSize(t, f, 0)
		f.Write([]byte("hello, world\n"))
		checkSize(t, f, 13)
		f.Truncate(10)
		checkSize(t, f, 10)
		f.Truncate(1024)
		checkSize(t, f, 1024)
		f.Truncate(0)
		checkSize(t, f, 0)
		_, err := f.Write([]byte("surprise!"))
		if err == nil ***REMOVED***
			checkSize(t, f, 13+9) // wrote at offset past where hello, world was.
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestSeek(t *testing.T) ***REMOVED***
	defer removeAllTestFiles(t)
	for _, fs := range Fss ***REMOVED***
		f := tmpFile(fs)
		defer f.Close()

		const data = "hello, world\n"
		io.WriteString(f, data)

		type test struct ***REMOVED***
			in     int64
			whence int
			out    int64
		***REMOVED***
		var tests = []test***REMOVED***
			***REMOVED***0, 1, int64(len(data))***REMOVED***,
			***REMOVED***0, 0, 0***REMOVED***,
			***REMOVED***5, 0, 5***REMOVED***,
			***REMOVED***0, 2, int64(len(data))***REMOVED***,
			***REMOVED***0, 0, 0***REMOVED***,
			***REMOVED***-1, 2, int64(len(data)) - 1***REMOVED***,
			***REMOVED***1 << 33, 0, 1 << 33***REMOVED***,
			***REMOVED***1 << 33, 2, 1<<33 + int64(len(data))***REMOVED***,
		***REMOVED***
		for i, tt := range tests ***REMOVED***
			off, err := f.Seek(tt.in, tt.whence)
			if off != tt.out || err != nil ***REMOVED***
				if e, ok := err.(*os.PathError); ok && e.Err == syscall.EINVAL && tt.out > 1<<32 ***REMOVED***
					// Reiserfs rejects the big seeks.
					// http://code.google.com/p/go/issues/detail?id=91
					break
				***REMOVED***
				t.Errorf("#%d: Seek(%v, %v) = %v, %v want %v, nil", i, tt.in, tt.whence, off, err, tt.out)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestReadAt(t *testing.T) ***REMOVED***
	defer removeAllTestFiles(t)
	for _, fs := range Fss ***REMOVED***
		f := tmpFile(fs)
		defer f.Close()

		const data = "hello, world\n"
		io.WriteString(f, data)

		b := make([]byte, 5)
		n, err := f.ReadAt(b, 7)
		if err != nil || n != len(b) ***REMOVED***
			t.Fatalf("ReadAt 7: %d, %v", n, err)
		***REMOVED***
		if string(b) != "world" ***REMOVED***
			t.Fatalf("ReadAt 7: have %q want %q", string(b), "world")
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestWriteAt(t *testing.T) ***REMOVED***
	defer removeAllTestFiles(t)
	for _, fs := range Fss ***REMOVED***
		f := tmpFile(fs)
		defer f.Close()

		const data = "hello, world\n"
		io.WriteString(f, data)

		n, err := f.WriteAt([]byte("WORLD"), 7)
		if err != nil || n != 5 ***REMOVED***
			t.Fatalf("WriteAt 7: %d, %v", n, err)
		***REMOVED***

		f2, err := fs.Open(f.Name())
		if err != nil ***REMOVED***
			t.Fatalf("%v: ReadFile %s: %v", fs.Name(), f.Name(), err)
		***REMOVED***
		defer f2.Close()
		buf := new(bytes.Buffer)
		buf.ReadFrom(f2)
		b := buf.Bytes()
		if string(b) != "hello, WORLD\n" ***REMOVED***
			t.Fatalf("after write: have %q want %q", string(b), "hello, WORLD\n")
		***REMOVED***

	***REMOVED***
***REMOVED***

func setupTestDir(t *testing.T, fs Fs) string ***REMOVED***
	path := testDir(fs)
	return setupTestFiles(t, fs, path)
***REMOVED***

func setupTestDirRoot(t *testing.T, fs Fs) string ***REMOVED***
	path := testDir(fs)
	setupTestFiles(t, fs, path)
	return path
***REMOVED***

func setupTestDirReusePath(t *testing.T, fs Fs, path string) string ***REMOVED***
	testRegistry[fs] = append(testRegistry[fs], path)
	return setupTestFiles(t, fs, path)
***REMOVED***

func setupTestFiles(t *testing.T, fs Fs, path string) string ***REMOVED***
	testSubDir := filepath.Join(path, "more", "subdirectories", "for", "testing", "we")
	err := fs.MkdirAll(testSubDir, 0700)
	if err != nil && !os.IsExist(err) ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	f, err := fs.Create(filepath.Join(testSubDir, "testfile1"))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	f.WriteString("Testfile 1 content")
	f.Close()

	f, err = fs.Create(filepath.Join(testSubDir, "testfile2"))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	f.WriteString("Testfile 2 content")
	f.Close()

	f, err = fs.Create(filepath.Join(testSubDir, "testfile3"))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	f.WriteString("Testfile 3 content")
	f.Close()

	f, err = fs.Create(filepath.Join(testSubDir, "testfile4"))
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	f.WriteString("Testfile 4 content")
	f.Close()
	return testSubDir
***REMOVED***

func TestReaddirnames(t *testing.T) ***REMOVED***
	defer removeAllTestFiles(t)
	for _, fs := range Fss ***REMOVED***
		testSubDir := setupTestDir(t, fs)
		tDir := filepath.Dir(testSubDir)

		root, err := fs.Open(tDir)
		if err != nil ***REMOVED***
			t.Fatal(fs.Name(), tDir, err)
		***REMOVED***
		defer root.Close()

		namesRoot, err := root.Readdirnames(-1)
		if err != nil ***REMOVED***
			t.Fatal(fs.Name(), namesRoot, err)
		***REMOVED***

		sub, err := fs.Open(testSubDir)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		defer sub.Close()

		namesSub, err := sub.Readdirnames(-1)
		if err != nil ***REMOVED***
			t.Fatal(fs.Name(), namesSub, err)
		***REMOVED***

		findNames(fs, t, tDir, testSubDir, namesRoot, namesSub)
	***REMOVED***
***REMOVED***

func TestReaddirSimple(t *testing.T) ***REMOVED***
	defer removeAllTestFiles(t)
	for _, fs := range Fss ***REMOVED***
		testSubDir := setupTestDir(t, fs)
		tDir := filepath.Dir(testSubDir)

		root, err := fs.Open(tDir)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		defer root.Close()

		rootInfo, err := root.Readdir(1)
		if err != nil ***REMOVED***
			t.Log(myFileInfo(rootInfo))
			t.Error(err)
		***REMOVED***

		rootInfo, err = root.Readdir(5)
		if err != io.EOF ***REMOVED***
			t.Log(myFileInfo(rootInfo))
			t.Error(err)
		***REMOVED***

		sub, err := fs.Open(testSubDir)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		defer sub.Close()

		subInfo, err := sub.Readdir(5)
		if err != nil ***REMOVED***
			t.Log(myFileInfo(subInfo))
			t.Error(err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestReaddir(t *testing.T) ***REMOVED***
	defer removeAllTestFiles(t)
	for num := 0; num < 6; num++ ***REMOVED***
		outputs := make([]string, len(Fss))
		infos := make([]string, len(Fss))
		for i, fs := range Fss ***REMOVED***
			testSubDir := setupTestDir(t, fs)
			//tDir := filepath.Dir(testSubDir)
			root, err := fs.Open(testSubDir)
			if err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED***
			defer root.Close()

			for j := 0; j < 6; j++ ***REMOVED***
				info, err := root.Readdir(num)
				outputs[i] += fmt.Sprintf("%v  Error: %v\n", myFileInfo(info), err)
				infos[i] += fmt.Sprintln(len(info), err)
			***REMOVED***
		***REMOVED***

		fail := false
		for i, o := range infos ***REMOVED***
			if i == 0 ***REMOVED***
				continue
			***REMOVED***
			if o != infos[i-1] ***REMOVED***
				fail = true
				break
			***REMOVED***
		***REMOVED***
		if fail ***REMOVED***
			t.Log("Readdir outputs not equal for Readdir(", num, ")")
			for i, o := range outputs ***REMOVED***
				t.Log(Fss[i].Name())
				t.Log(o)
			***REMOVED***
			t.Fail()
		***REMOVED***
	***REMOVED***
***REMOVED***

type myFileInfo []os.FileInfo

func (m myFileInfo) String() string ***REMOVED***
	out := "Fileinfos:\n"
	for _, e := range m ***REMOVED***
		out += "  " + e.Name() + "\n"
	***REMOVED***
	return out
***REMOVED***

func TestReaddirAll(t *testing.T) ***REMOVED***
	defer removeAllTestFiles(t)
	for _, fs := range Fss ***REMOVED***
		testSubDir := setupTestDir(t, fs)
		tDir := filepath.Dir(testSubDir)

		root, err := fs.Open(tDir)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		defer root.Close()

		rootInfo, err := root.Readdir(-1)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		var namesRoot = []string***REMOVED******REMOVED***
		for _, e := range rootInfo ***REMOVED***
			namesRoot = append(namesRoot, e.Name())
		***REMOVED***

		sub, err := fs.Open(testSubDir)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		defer sub.Close()

		subInfo, err := sub.Readdir(-1)
		if err != nil ***REMOVED***
			t.Fatal(err)
		***REMOVED***
		var namesSub = []string***REMOVED******REMOVED***
		for _, e := range subInfo ***REMOVED***
			namesSub = append(namesSub, e.Name())
		***REMOVED***

		findNames(fs, t, tDir, testSubDir, namesRoot, namesSub)
	***REMOVED***
***REMOVED***

func findNames(fs Fs, t *testing.T, tDir, testSubDir string, root, sub []string) ***REMOVED***
	var foundRoot bool
	for _, e := range root ***REMOVED***
		f, err := fs.Open(filepath.Join(tDir, e))
		if err != nil ***REMOVED***
			t.Error("Open", filepath.Join(tDir, e), ":", err)
		***REMOVED***
		defer f.Close()

		if equal(e, "we") ***REMOVED***
			foundRoot = true
		***REMOVED***
	***REMOVED***
	if !foundRoot ***REMOVED***
		t.Logf("Names root: %v", root)
		t.Logf("Names sub: %v", sub)
		t.Error("Didn't find subdirectory we")
	***REMOVED***

	var found1, found2 bool
	for _, e := range sub ***REMOVED***
		f, err := fs.Open(filepath.Join(testSubDir, e))
		if err != nil ***REMOVED***
			t.Error("Open", filepath.Join(testSubDir, e), ":", err)
		***REMOVED***
		defer f.Close()

		if equal(e, "testfile1") ***REMOVED***
			found1 = true
		***REMOVED***
		if equal(e, "testfile2") ***REMOVED***
			found2 = true
		***REMOVED***
	***REMOVED***

	if !found1 ***REMOVED***
		t.Logf("Names root: %v", root)
		t.Logf("Names sub: %v", sub)
		t.Error("Didn't find testfile1")
	***REMOVED***
	if !found2 ***REMOVED***
		t.Logf("Names root: %v", root)
		t.Logf("Names sub: %v", sub)
		t.Error("Didn't find testfile2")
	***REMOVED***
***REMOVED***

func removeAllTestFiles(t *testing.T) ***REMOVED***
	for fs, list := range testRegistry ***REMOVED***
		for _, path := range list ***REMOVED***
			if err := fs.RemoveAll(path); err != nil ***REMOVED***
				t.Error(fs.Name(), err)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	testRegistry = make(map[Fs][]string)
***REMOVED***

func equal(name1, name2 string) (r bool) ***REMOVED***
	switch runtime.GOOS ***REMOVED***
	case "windows":
		r = strings.ToLower(name1) == strings.ToLower(name2)
	default:
		r = name1 == name2
	***REMOVED***
	return
***REMOVED***

func checkSize(t *testing.T, f File, size int64) ***REMOVED***
	dir, err := f.Stat()
	if err != nil ***REMOVED***
		t.Fatalf("Stat %q (looking for size %d): %s", f.Name(), size, err)
	***REMOVED***
	if dir.Size() != size ***REMOVED***
		t.Errorf("Stat %q: size %d want %d", f.Name(), dir.Size(), size)
	***REMOVED***
***REMOVED***
