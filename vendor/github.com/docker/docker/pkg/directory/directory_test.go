package directory

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

// Size of an empty directory should be 0
func TestSizeEmpty(t *testing.T) ***REMOVED***
	var dir string
	var err error
	if dir, err = ioutil.TempDir(os.TempDir(), "testSizeEmptyDirectory"); err != nil ***REMOVED***
		t.Fatalf("failed to create directory: %s", err)
	***REMOVED***

	var size int64
	if size, _ = Size(dir); size != 0 ***REMOVED***
		t.Fatalf("empty directory has size: %d", size)
	***REMOVED***
***REMOVED***

// Size of a directory with one empty file should be 0
func TestSizeEmptyFile(t *testing.T) ***REMOVED***
	var dir string
	var err error
	if dir, err = ioutil.TempDir(os.TempDir(), "testSizeEmptyFile"); err != nil ***REMOVED***
		t.Fatalf("failed to create directory: %s", err)
	***REMOVED***

	var file *os.File
	if file, err = ioutil.TempFile(dir, "file"); err != nil ***REMOVED***
		t.Fatalf("failed to create file: %s", err)
	***REMOVED***

	var size int64
	if size, _ = Size(file.Name()); size != 0 ***REMOVED***
		t.Fatalf("directory with one file has size: %d", size)
	***REMOVED***
***REMOVED***

// Size of a directory with one 5-byte file should be 5
func TestSizeNonemptyFile(t *testing.T) ***REMOVED***
	var dir string
	var err error
	if dir, err = ioutil.TempDir(os.TempDir(), "testSizeNonemptyFile"); err != nil ***REMOVED***
		t.Fatalf("failed to create directory: %s", err)
	***REMOVED***

	var file *os.File
	if file, err = ioutil.TempFile(dir, "file"); err != nil ***REMOVED***
		t.Fatalf("failed to create file: %s", err)
	***REMOVED***

	d := []byte***REMOVED***97, 98, 99, 100, 101***REMOVED***
	file.Write(d)

	var size int64
	if size, _ = Size(file.Name()); size != 5 ***REMOVED***
		t.Fatalf("directory with one 5-byte file has size: %d", size)
	***REMOVED***
***REMOVED***

// Size of a directory with one empty directory should be 0
func TestSizeNestedDirectoryEmpty(t *testing.T) ***REMOVED***
	var dir string
	var err error
	if dir, err = ioutil.TempDir(os.TempDir(), "testSizeNestedDirectoryEmpty"); err != nil ***REMOVED***
		t.Fatalf("failed to create directory: %s", err)
	***REMOVED***
	if dir, err = ioutil.TempDir(dir, "nested"); err != nil ***REMOVED***
		t.Fatalf("failed to create nested directory: %s", err)
	***REMOVED***

	var size int64
	if size, _ = Size(dir); size != 0 ***REMOVED***
		t.Fatalf("directory with one empty directory has size: %d", size)
	***REMOVED***
***REMOVED***

// Test directory with 1 file and 1 empty directory
func TestSizeFileAndNestedDirectoryEmpty(t *testing.T) ***REMOVED***
	var dir string
	var err error
	if dir, err = ioutil.TempDir(os.TempDir(), "testSizeFileAndNestedDirectoryEmpty"); err != nil ***REMOVED***
		t.Fatalf("failed to create directory: %s", err)
	***REMOVED***
	if dir, err = ioutil.TempDir(dir, "nested"); err != nil ***REMOVED***
		t.Fatalf("failed to create nested directory: %s", err)
	***REMOVED***

	var file *os.File
	if file, err = ioutil.TempFile(dir, "file"); err != nil ***REMOVED***
		t.Fatalf("failed to create file: %s", err)
	***REMOVED***

	d := []byte***REMOVED***100, 111, 99, 107, 101, 114***REMOVED***
	file.Write(d)

	var size int64
	if size, _ = Size(dir); size != 6 ***REMOVED***
		t.Fatalf("directory with 6-byte file and empty directory has size: %d", size)
	***REMOVED***
***REMOVED***

// Test directory with 1 file and 1 non-empty directory
func TestSizeFileAndNestedDirectoryNonempty(t *testing.T) ***REMOVED***
	var dir, dirNested string
	var err error
	if dir, err = ioutil.TempDir(os.TempDir(), "TestSizeFileAndNestedDirectoryNonempty"); err != nil ***REMOVED***
		t.Fatalf("failed to create directory: %s", err)
	***REMOVED***
	if dirNested, err = ioutil.TempDir(dir, "nested"); err != nil ***REMOVED***
		t.Fatalf("failed to create nested directory: %s", err)
	***REMOVED***

	var file *os.File
	if file, err = ioutil.TempFile(dir, "file"); err != nil ***REMOVED***
		t.Fatalf("failed to create file: %s", err)
	***REMOVED***

	data := []byte***REMOVED***100, 111, 99, 107, 101, 114***REMOVED***
	file.Write(data)

	var nestedFile *os.File
	if nestedFile, err = ioutil.TempFile(dirNested, "file"); err != nil ***REMOVED***
		t.Fatalf("failed to create file in nested directory: %s", err)
	***REMOVED***

	nestedData := []byte***REMOVED***100, 111, 99, 107, 101, 114***REMOVED***
	nestedFile.Write(nestedData)

	var size int64
	if size, _ = Size(dir); size != 12 ***REMOVED***
		t.Fatalf("directory with 6-byte file and nested directory with 6-byte file has size: %d", size)
	***REMOVED***
***REMOVED***

// Test migration of directory to a subdir underneath itself
func TestMoveToSubdir(t *testing.T) ***REMOVED***
	var outerDir, subDir string
	var err error

	if outerDir, err = ioutil.TempDir(os.TempDir(), "TestMoveToSubdir"); err != nil ***REMOVED***
		t.Fatalf("failed to create directory: %v", err)
	***REMOVED***

	if subDir, err = ioutil.TempDir(outerDir, "testSub"); err != nil ***REMOVED***
		t.Fatalf("failed to create subdirectory: %v", err)
	***REMOVED***

	// write 4 temp files in the outer dir to get moved
	filesList := []string***REMOVED***"a", "b", "c", "d"***REMOVED***
	for _, fName := range filesList ***REMOVED***
		if file, err := os.Create(filepath.Join(outerDir, fName)); err != nil ***REMOVED***
			t.Fatalf("couldn't create temp file %q: %v", fName, err)
		***REMOVED*** else ***REMOVED***
			file.WriteString(fName)
			file.Close()
		***REMOVED***
	***REMOVED***

	if err = MoveToSubdir(outerDir, filepath.Base(subDir)); err != nil ***REMOVED***
		t.Fatalf("Error during migration of content to subdirectory: %v", err)
	***REMOVED***
	// validate that the files were moved to the subdirectory
	infos, err := ioutil.ReadDir(subDir)
	if err != nil ***REMOVED***
		t.Fatal(err)
	***REMOVED***
	if len(infos) != 4 ***REMOVED***
		t.Fatalf("Should be four files in the subdir after the migration: actual length: %d", len(infos))
	***REMOVED***
	var results []string
	for _, info := range infos ***REMOVED***
		results = append(results, info.Name())
	***REMOVED***
	sort.Sort(sort.StringSlice(results))
	if !reflect.DeepEqual(filesList, results) ***REMOVED***
		t.Fatalf("Results after migration do not equal list of files: expected: %v, got: %v", filesList, results)
	***REMOVED***
***REMOVED***

// Test a non-existing directory
func TestSizeNonExistingDirectory(t *testing.T) ***REMOVED***
	if _, err := Size("/thisdirectoryshouldnotexist/TestSizeNonExistingDirectory"); err == nil ***REMOVED***
		t.Fatalf("error is expected")
	***REMOVED***
***REMOVED***
