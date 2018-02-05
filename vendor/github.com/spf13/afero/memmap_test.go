package afero

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestNormalizePath(t *testing.T) ***REMOVED***
	type test struct ***REMOVED***
		input    string
		expected string
	***REMOVED***

	data := []test***REMOVED***
		***REMOVED***".", FilePathSeparator***REMOVED***,
		***REMOVED***"./", FilePathSeparator***REMOVED***,
		***REMOVED***"..", FilePathSeparator***REMOVED***,
		***REMOVED***"../", FilePathSeparator***REMOVED***,
		***REMOVED***"./..", FilePathSeparator***REMOVED***,
		***REMOVED***"./../", FilePathSeparator***REMOVED***,
	***REMOVED***

	for i, d := range data ***REMOVED***
		cpath := normalizePath(d.input)
		if d.expected != cpath ***REMOVED***
			t.Errorf("Test %d failed. Expected %q got %q", i, d.expected, cpath)
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestPathErrors(t *testing.T) ***REMOVED***
	path := filepath.Join(".", "some", "path")
	path2 := filepath.Join(".", "different", "path")
	fs := NewMemMapFs()
	perm := os.FileMode(0755)

	// relevant functions:
	// func (m *MemMapFs) Chmod(name string, mode os.FileMode) error
	// func (m *MemMapFs) Chtimes(name string, atime time.Time, mtime time.Time) error
	// func (m *MemMapFs) Create(name string) (File, error)
	// func (m *MemMapFs) Mkdir(name string, perm os.FileMode) error
	// func (m *MemMapFs) MkdirAll(path string, perm os.FileMode) error
	// func (m *MemMapFs) Open(name string) (File, error)
	// func (m *MemMapFs) OpenFile(name string, flag int, perm os.FileMode) (File, error)
	// func (m *MemMapFs) Remove(name string) error
	// func (m *MemMapFs) Rename(oldname, newname string) error
	// func (m *MemMapFs) Stat(name string) (os.FileInfo, error)

	err := fs.Chmod(path, perm)
	checkPathError(t, err, "Chmod")

	err = fs.Chtimes(path, time.Now(), time.Now())
	checkPathError(t, err, "Chtimes")

	// fs.Create doesn't return an error

	err = fs.Mkdir(path2, perm)
	if err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***
	err = fs.Mkdir(path2, perm)
	checkPathError(t, err, "Mkdir")

	err = fs.MkdirAll(path2, perm)
	if err != nil ***REMOVED***
		t.Error("MkdirAll:", err)
	***REMOVED***

	_, err = fs.Open(path)
	checkPathError(t, err, "Open")

	_, err = fs.OpenFile(path, os.O_RDWR, perm)
	checkPathError(t, err, "OpenFile")

	err = fs.Remove(path)
	checkPathError(t, err, "Remove")

	err = fs.RemoveAll(path)
	if err != nil ***REMOVED***
		t.Error("RemoveAll:", err)
	***REMOVED***

	err = fs.Rename(path, path2)
	checkPathError(t, err, "Rename")

	_, err = fs.Stat(path)
	checkPathError(t, err, "Stat")
***REMOVED***

func checkPathError(t *testing.T, err error, op string) ***REMOVED***
	pathErr, ok := err.(*os.PathError)
	if !ok ***REMOVED***
		t.Error(op+":", err, "is not a os.PathError")
		return
	***REMOVED***
	_, ok = pathErr.Err.(*os.PathError)
	if ok ***REMOVED***
		t.Error(op+":", err, "contains another os.PathError")
	***REMOVED***
***REMOVED***

// Ensure Permissions are set on OpenFile/Mkdir/MkdirAll
func TestPermSet(t *testing.T) ***REMOVED***
	const fileName = "/myFileTest"
	const dirPath = "/myDirTest"
	const dirPathAll = "/my/path/to/dir"

	const fileMode = os.FileMode(0765)
	// directories will also have the directory bit set
	const dirMode = fileMode | os.ModeDir

	fs := NewMemMapFs()

	// Test Openfile
	f, err := fs.OpenFile(fileName, os.O_CREATE, fileMode)
	if err != nil ***REMOVED***
		t.Errorf("OpenFile Create failed: %s", err)
		return
	***REMOVED***
	f.Close()

	s, err := fs.Stat(fileName)
	if err != nil ***REMOVED***
		t.Errorf("Stat failed: %s", err)
		return
	***REMOVED***
	if s.Mode().String() != fileMode.String() ***REMOVED***
		t.Errorf("Permissions Incorrect: %s != %s", s.Mode().String(), fileMode.String())
		return
	***REMOVED***

	// Test Mkdir
	err = fs.Mkdir(dirPath, dirMode)
	if err != nil ***REMOVED***
		t.Errorf("MkDir Create failed: %s", err)
		return
	***REMOVED***
	s, err = fs.Stat(dirPath)
	if err != nil ***REMOVED***
		t.Errorf("Stat failed: %s", err)
		return
	***REMOVED***
	// sets File
	if s.Mode().String() != dirMode.String() ***REMOVED***
		t.Errorf("Permissions Incorrect: %s != %s", s.Mode().String(), dirMode.String())
		return
	***REMOVED***

	// Test MkdirAll
	err = fs.MkdirAll(dirPathAll, dirMode)
	if err != nil ***REMOVED***
		t.Errorf("MkDir Create failed: %s", err)
		return
	***REMOVED***
	s, err = fs.Stat(dirPathAll)
	if err != nil ***REMOVED***
		t.Errorf("Stat failed: %s", err)
		return
	***REMOVED***
	if s.Mode().String() != dirMode.String() ***REMOVED***
		t.Errorf("Permissions Incorrect: %s != %s", s.Mode().String(), dirMode.String())
		return
	***REMOVED***
***REMOVED***

// Fails if multiple file objects use the same file.at counter in MemMapFs
func TestMultipleOpenFiles(t *testing.T) ***REMOVED***
	defer removeAllTestFiles(t)
	const fileName = "afero-demo2.txt"

	var data = make([][]byte, len(Fss))

	for i, fs := range Fss ***REMOVED***
		dir := testDir(fs)
		path := filepath.Join(dir, fileName)
		fh1, err := fs.Create(path)
		if err != nil ***REMOVED***
			t.Error("fs.Create failed: " + err.Error())
		***REMOVED***
		_, err = fh1.Write([]byte("test"))
		if err != nil ***REMOVED***
			t.Error("fh.Write failed: " + err.Error())
		***REMOVED***
		_, err = fh1.Seek(0, os.SEEK_SET)
		if err != nil ***REMOVED***
			t.Error(err)
		***REMOVED***

		fh2, err := fs.OpenFile(path, os.O_RDWR, 0777)
		if err != nil ***REMOVED***
			t.Error("fs.OpenFile failed: " + err.Error())
		***REMOVED***
		_, err = fh2.Seek(0, os.SEEK_END)
		if err != nil ***REMOVED***
			t.Error(err)
		***REMOVED***
		_, err = fh2.Write([]byte("data"))
		if err != nil ***REMOVED***
			t.Error(err)
		***REMOVED***
		err = fh2.Close()
		if err != nil ***REMOVED***
			t.Error(err)
		***REMOVED***

		_, err = fh1.Write([]byte("data"))
		if err != nil ***REMOVED***
			t.Error(err)
		***REMOVED***
		err = fh1.Close()
		if err != nil ***REMOVED***
			t.Error(err)
		***REMOVED***
		// the file now should contain "datadata"
		data[i], err = ReadFile(fs, path)
		if err != nil ***REMOVED***
			t.Error(err)
		***REMOVED***
	***REMOVED***

	for i, fs := range Fss ***REMOVED***
		if i == 0 ***REMOVED***
			continue
		***REMOVED***
		if string(data[0]) != string(data[i]) ***REMOVED***
			t.Errorf("%s and %s don't behave the same\n"+
				"%s: \"%s\"\n%s: \"%s\"\n",
				Fss[0].Name(), fs.Name(), Fss[0].Name(), data[0], fs.Name(), data[i])
		***REMOVED***
	***REMOVED***
***REMOVED***

// Test if file.Write() fails when opened as read only
func TestReadOnly(t *testing.T) ***REMOVED***
	defer removeAllTestFiles(t)
	const fileName = "afero-demo.txt"

	for _, fs := range Fss ***REMOVED***
		dir := testDir(fs)
		path := filepath.Join(dir, fileName)

		f, err := fs.Create(path)
		if err != nil ***REMOVED***
			t.Error(fs.Name()+":", "fs.Create failed: "+err.Error())
		***REMOVED***
		_, err = f.Write([]byte("test"))
		if err != nil ***REMOVED***
			t.Error(fs.Name()+":", "Write failed: "+err.Error())
		***REMOVED***
		f.Close()

		f, err = fs.Open(path)
		if err != nil ***REMOVED***
			t.Error("fs.Open failed: " + err.Error())
		***REMOVED***
		_, err = f.Write([]byte("data"))
		if err == nil ***REMOVED***
			t.Error(fs.Name()+":", "No write error")
		***REMOVED***
		f.Close()

		f, err = fs.OpenFile(path, os.O_RDONLY, 0644)
		if err != nil ***REMOVED***
			t.Error("fs.Open failed: " + err.Error())
		***REMOVED***
		_, err = f.Write([]byte("data"))
		if err == nil ***REMOVED***
			t.Error(fs.Name()+":", "No write error")
		***REMOVED***
		f.Close()
	***REMOVED***
***REMOVED***

func TestWriteCloseTime(t *testing.T) ***REMOVED***
	defer removeAllTestFiles(t)
	const fileName = "afero-demo.txt"

	for _, fs := range Fss ***REMOVED***
		dir := testDir(fs)
		path := filepath.Join(dir, fileName)

		f, err := fs.Create(path)
		if err != nil ***REMOVED***
			t.Error(fs.Name()+":", "fs.Create failed: "+err.Error())
		***REMOVED***
		f.Close()

		f, err = fs.Create(path)
		if err != nil ***REMOVED***
			t.Error(fs.Name()+":", "fs.Create failed: "+err.Error())
		***REMOVED***
		fi, err := f.Stat()
		if err != nil ***REMOVED***
			t.Error(fs.Name()+":", "Stat failed: "+err.Error())
		***REMOVED***
		timeBefore := fi.ModTime()

		// sorry for the delay, but we have to make sure time advances,
		// also on non Un*x systems...
		switch runtime.GOOS ***REMOVED***
		case "windows":
			time.Sleep(2 * time.Second)
		case "darwin":
			time.Sleep(1 * time.Second)
		default: // depending on the FS, this may work with < 1 second, on my old ext3 it does not
			time.Sleep(1 * time.Second)
		***REMOVED***

		_, err = f.Write([]byte("test"))
		if err != nil ***REMOVED***
			t.Error(fs.Name()+":", "Write failed: "+err.Error())
		***REMOVED***
		f.Close()
		fi, err = fs.Stat(path)
		if err != nil ***REMOVED***
			t.Error(fs.Name()+":", "fs.Stat failed: "+err.Error())
		***REMOVED***
		if fi.ModTime().Equal(timeBefore) ***REMOVED***
			t.Error(fs.Name()+":", "ModTime was not set on Close()")
		***REMOVED***
	***REMOVED***
***REMOVED***

// This test should be run with the race detector on:
// go test -race -v -timeout 10s -run TestRacingDeleteAndClose
func TestRacingDeleteAndClose(t *testing.T) ***REMOVED***
	fs := NewMemMapFs()
	pathname := "testfile"
	f, err := fs.Create(pathname)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	in := make(chan bool)

	go func() ***REMOVED***
		<-in
		f.Close()
	***REMOVED***()
	go func() ***REMOVED***
		<-in
		fs.Remove(pathname)
	***REMOVED***()
	close(in)
***REMOVED***

// This test should be run with the race detector on:
// go test -run TestMemFsDataRace -race
func TestMemFsDataRace(t *testing.T) ***REMOVED***
	const dir = "test_dir"
	fs := NewMemMapFs()

	if err := fs.MkdirAll(dir, 0777); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	const n = 1000
	done := make(chan struct***REMOVED******REMOVED***)

	go func() ***REMOVED***
		defer close(done)
		for i := 0; i < n; i++ ***REMOVED***
			fname := filepath.Join(dir, fmt.Sprintf("%d.txt", i))
			if err := WriteFile(fs, fname, []byte(""), 0777); err != nil ***REMOVED***
				panic(err)
			***REMOVED***
			if err := fs.Remove(fname); err != nil ***REMOVED***
				panic(err)
			***REMOVED***
		***REMOVED***
	***REMOVED***()

loop:
	for ***REMOVED***
		select ***REMOVED***
		case <-done:
			break loop
		default:
			_, err := ReadDir(fs, dir)
			if err != nil ***REMOVED***
				t.Fatal(err)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func TestMemFsDirMode(t *testing.T) ***REMOVED***
	fs := NewMemMapFs()
	err := fs.Mkdir("/testDir1", 0644)
	if err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***
	err = fs.MkdirAll("/sub/testDir2", 0644)
	if err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***
	info, err := fs.Stat("/testDir1")
	if err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***
	if !info.IsDir() ***REMOVED***
		t.Error("should be a directory")
	***REMOVED***
	if !info.Mode().IsDir() ***REMOVED***
		t.Error("FileMode is not directory")
	***REMOVED***
	info, err = fs.Stat("/sub/testDir2")
	if err != nil ***REMOVED***
		t.Error(err)
	***REMOVED***
	if !info.IsDir() ***REMOVED***
		t.Error("should be a directory")
	***REMOVED***
	if !info.Mode().IsDir() ***REMOVED***
		t.Error("FileMode is not directory")
	***REMOVED***
***REMOVED***

func TestMemFsUnexpectedEOF(t *testing.T) ***REMOVED***
	t.Parallel()

	fs := NewMemMapFs()

	if err := WriteFile(fs, "file.txt", []byte("abc"), 0777); err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	f, err := fs.Open("file.txt")
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	defer f.Close()

	// Seek beyond the end.
	_, err = f.Seek(512, 0)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***

	buff := make([]byte, 256)
	_, err = io.ReadAtLeast(f, buff, 256)

	if err != io.ErrUnexpectedEOF ***REMOVED***
		t.Fatal("Expected ErrUnexpectedEOF")
	***REMOVED***
***REMOVED***
